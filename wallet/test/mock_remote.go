package test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"perun.network/perun-cardano-backend/channel/test"
	"perun.network/perun-cardano-backend/channel/types"
	"perun.network/perun-cardano-backend/wallet"
	"perun.network/perun-cardano-backend/wallet/address"
	"perun.network/perun-cardano-backend/wire"
)

// MockRemote should only be instantiated using NewMockRemote.
// The default implementation has one valid signature tuple:
// (MockMessage, MockSignature, MockAddress).
type MockRemote struct {
	MockAddress     address.Address
	MockPubKeyBytes [address.PubKeyLength]byte
	// UnavailableAddress is a valid wallet.PubKey that has associated account (private key) in this remote wallet.
	UnavailableAddress address.Address
	// InvalidPubKeyBytes is invalid because it is not exactly wallet.PubKeyLength bytes long.
	InvalidPubKeyBytes []byte

	MockSignature       []byte
	MockSignatureString string
	// OtherSignature is a correctly encoded signature that is not valid for any (message, public key) pair.
	OtherSignature       []byte
	OtherSignatureString string

	// InvalidSignatureLonger is a signature that has a length longer than wallet.SignatureLength.
	InvalidSignatureLonger []byte
	// InvalidSignatureShorter is a signature that has a length shorter than wallet.SignatureLength.
	InvalidSignatureShorter []byte

	MockMessage       []byte
	MockMessageString string

	MockChannelState types.ChannelState

	callEndpoint func(string, interface{}, interface{}) error
}

func NewMockRemote(rng *rand.Rand) *MockRemote {
	r := &MockRemote{}
	initializeRandomValues(r, rng)

	r.callEndpoint = makeCallEndpointDefault(r)
	return r
}

func initializeRandomValues(r *MockRemote, rng *rand.Rand) {
	const maxMessageLength = 0x100 // in bytes

	r.MockAddress = MakeRandomAddress(rng)
	r.MockPubKeyBytes = r.MockAddress.GetPubKey()
	r.UnavailableAddress = MakeRandomAddress(rng)
	for bytes.Equal(r.UnavailableAddress.GetPubKeySlice(), r.MockPubKeyBytes[:]) {
		r.UnavailableAddress = MakeRandomAddress(rng)
	}

	if rng.Int()%2 == 0 {
		r.InvalidPubKeyBytes = MakeTooFewPublicKeyBytes(rng)
	} else {
		r.InvalidPubKeyBytes = MakeTooManyPublicKeyBytes(rng)
	}
	rng.Read(r.InvalidPubKeyBytes)

	r.MockSignature = MakeRandomSignature(rng)
	r.MockSignatureString = hex.EncodeToString(r.MockSignature)

	r.OtherSignature = MakeRandomSignature(rng)
	for bytes.Equal(r.MockSignature, r.OtherSignature) {
		rng.Read(r.OtherSignature)
	}
	r.OtherSignatureString = hex.EncodeToString(r.OtherSignature)

	r.InvalidSignatureShorter = MakeTooShortSignature(rng)

	r.InvalidSignatureLonger = MakeTooLongSignature(rng)

	r.MockMessage = GetRandomByteSlice(0, maxMessageLength, rng)
	r.MockMessageString = hex.EncodeToString(r.MockMessage)

	r.MockChannelState = test.MakeRandomChannelState(rng)
}

func (m *MockRemote) SetCallEndpoint(f func(string, interface{}, interface{}) error) {
	m.callEndpoint = f
}

func makeCallEndpointDefault(r *MockRemote) func(string, interface{}, interface{}) error {
	return func(endpoint string, req interface{}, resp interface{}) error {
		switch endpoint {
		case wallet.EndpointSignData:
			request, ok := req.(wire.SigningRequest)
			if !ok {
				return fmt.Errorf("unable to cast request to SingingRequest")
			}
			response, ok := resp.(*wire.SigningResponse)
			if !ok {
				return fmt.Errorf("unable to cast resp to SingingResponse")
			}
			return callSign(r, request, response)
		case wallet.EndpointVerifyDataSignature:
			request, ok := req.(wire.VerificationRequest)
			if !ok {
				return fmt.Errorf("unable to cast request to VerificationRequest")
			}
			response, ok := resp.(*wire.VerificationResponse)
			if !ok {
				return fmt.Errorf("unable to cast resp to VerificationResponse")
			}
			return callVerify(r, request, response)
		case wallet.EndpointKeyAvailable:
			request, ok := req.(wire.KeyAvailabilityRequest)
			if !ok {
				return fmt.Errorf("unable to cast request to KeyAvailabilityRequst")
			}
			response, ok := resp.(*wire.KeyAvailabilityResponse)
			if !ok {
				return fmt.Errorf("unable to cast resp to KeyAvailabilityResponse")
			}
			return callKeyAvailable(r, request, response)
		case wallet.EndpointSignChannelState:
			request, ok := req.(wire.ChannelStateSigningRequest)
			if !ok {
				return fmt.Errorf("unable to cast request to ChannelStateSigningRequest")
			}
			response, ok := resp.(*wire.SigningResponse)
			if !ok {
				return fmt.Errorf("unable to cast resp to SingingResponse")
			}
			return callSignChannelState(r, request, response)
		case wallet.EndpointVerifyChannelStateSignature:
			request, ok := req.(wire.ChannelStateVerificationRequest)
			if !ok {
				return fmt.Errorf("unable to cast request to ChannelStateVerificationRequest")
			}
			response, ok := resp.(*wire.VerificationResponse)
			if !ok {
				return fmt.Errorf("unable to cast resp to VerificationResponse")
			}
			return callVerifyChannelState(r, request, response)
		default:
			return fmt.Errorf("unable to recognize endpoint: %s", endpoint)

		}
	}
}

func (m *MockRemote) CallEndpoint(endpoint string, request interface{}, response interface{}) error {
	return m.callEndpoint(endpoint, request, response)
}

func callSign(r *MockRemote, request wire.SigningRequest, response *wire.SigningResponse) error {
	reqAddr, err := request.PubKey.Decode()
	if err != nil {
		return fmt.Errorf("unable to decode PubKey from request")
	}
	if !reqAddr.Equal(&r.MockAddress) {
		return fmt.Errorf("invalid public key for mock remote")
	}

	if request.Message != r.MockMessageString {
		return fmt.Errorf("invalid data for mock remote")
	}
	response.Hex = r.MockSignatureString
	return nil
}

func callVerify(r *MockRemote, request wire.VerificationRequest, response *wire.VerificationResponse) error {
	reqAddr, err := request.PubKey.Decode()
	if err != nil {
		return fmt.Errorf("unable to decode PubKey from request")
	}
	if !reqAddr.Equal(&r.MockAddress) && !reqAddr.Equal(&r.UnavailableAddress) {
		return fmt.Errorf("invalid public key for mock remote")
	}
	if reqAddr.Equal(&r.UnavailableAddress) {
		*response = false
		return nil
	}
	if request.Message != r.MockMessageString {
		return fmt.Errorf("invalid data for mock remote")
	}
	if request.Signature.Hex == r.MockSignatureString {
		*response = true
		return nil
	}
	if request.Signature.Hex == r.OtherSignatureString {
		*response = false
		return nil
	}
	if request.Signature.Hex == hex.EncodeToString(r.InvalidSignatureShorter) ||
		request.Signature.Hex == hex.EncodeToString(r.InvalidSignatureLonger) {
		panic("mock remote received signature of invalid length to verify")
	}
	return fmt.Errorf("invalid signature for mock remote")
}

func callSignChannelState(r *MockRemote, request wire.ChannelStateSigningRequest, response *wire.SigningResponse) error {
	reqAddr, err := request.PubKey.Decode()
	if err != nil {
		return fmt.Errorf("unable to decode PubKey from request")
	}
	if !reqAddr.Equal(&r.MockAddress) {
		return fmt.Errorf("invalid public key for mock remote")
	}
	if !request.ChannelState.Decode().Equal(r.MockChannelState) {
		return fmt.Errorf("invalid channel state for mock remote")
	}
	response.Hex = r.MockSignatureString
	return nil
}

func callVerifyChannelState(r *MockRemote, request wire.ChannelStateVerificationRequest, response *wire.VerificationResponse) error {
	reqAddr, err := request.PubKey.Decode()
	if err != nil {
		return fmt.Errorf("unable to decode PubKey from request")
	}
	if !reqAddr.Equal(&r.MockAddress) && !reqAddr.Equal(&r.UnavailableAddress) {
		return fmt.Errorf("invalid public key for mock remote")
	}
	if reqAddr.Equal(&r.UnavailableAddress) {
		*response = false
		return nil
	}
	if !request.ChannelState.Decode().Equal(r.MockChannelState) {
		return fmt.Errorf("invalid data for mock remote")
	}
	if request.Signature.Hex == r.MockSignatureString {
		*response = true
		return nil
	}
	if request.Signature.Hex == r.OtherSignatureString {
		*response = false
		return nil
	}
	if request.Signature.Hex == hex.EncodeToString(r.InvalidSignatureShorter) ||
		request.Signature.Hex == hex.EncodeToString(r.InvalidSignatureLonger) {
		panic("mock remote received signature of invalid length to verify")
	}
	return fmt.Errorf("invalid signature for mock remote")
}

func callKeyAvailable(r *MockRemote, request wire.KeyAvailabilityRequest, response *wire.KeyAvailabilityResponse) error {
	reqAddr, err := request.Decode()
	if err != nil {
		return fmt.Errorf("unable to decode address from request: %w", err)
	}
	*response = reqAddr.Equal(&r.MockAddress)
	return nil
}
