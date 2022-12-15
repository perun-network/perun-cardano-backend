package test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
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
}

func (m *MockRemote) SetCallEndpoint(f func(string, interface{}, interface{}) error) {
	m.callEndpoint = f
}

func makeCallEndpointDefault(r *MockRemote) func(string, interface{}, interface{}) error {
	return func(endpoint string, req interface{}, response interface{}) error {
		switch endpoint {
		case wallet.EndpointSignData:
			request, ok := req.(wire.SigningRequest)
			if !ok {
				return fmt.Errorf("unable to cast request to SingingRequest")
			}
			resp, ok := response.(*wire.SigningResponse)
			if !ok {
				return fmt.Errorf("unable to cast response to SingingResponse")
			}
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
			resp.Hex = r.MockSignatureString
			return nil
		case wallet.EndpointVerifyDataSignature:
			request, ok := req.(wire.VerificationRequest)
			if !ok {
				return fmt.Errorf("unable to cast request to VerificationRequest")
			}
			resp, ok := response.(*wire.VerificationResponse)
			if !ok {
				return fmt.Errorf("unable to cast response to VerificationResponse")
			}
			reqAddr, err := request.PubKey.Decode()
			if err != nil {
				return fmt.Errorf("unable to decode PubKey from request")
			}
			if !reqAddr.Equal(&r.MockAddress) && !reqAddr.Equal(&r.UnavailableAddress) {
				return fmt.Errorf("invalid public key for mock remote")
			}
			if reqAddr.Equal(&r.UnavailableAddress) {
				*resp = false
				return nil
			}
			if request.Message != r.MockMessageString {
				return fmt.Errorf("invalid data for mock remote")
			}
			if request.Signature.Hex == r.MockSignatureString {
				*resp = true
				return nil
			}
			if request.Signature.Hex == r.OtherSignatureString {
				*resp = false
				return nil
			}
			if request.Signature.Hex == hex.EncodeToString(r.InvalidSignatureShorter) ||
				request.Signature.Hex == hex.EncodeToString(r.InvalidSignatureLonger) {
				panic("mock remote received signature of invalid length to verify")
			}
			return fmt.Errorf("invalid signature for mock remote")
		case wallet.EndpointKeyAvailable:
			request, ok := req.(wire.KeyAvailabilityRequest)
			if !ok {
				return fmt.Errorf("unable to cast request to KeyAvailabilityRequst")
			}
			resp, ok := response.(*wire.KeyAvailabilityResponse)
			if !ok {
				return fmt.Errorf("unable to cast response to KeyAvailabilityResponse")
			}
			reqAddr, err := request.Decode()
			if err != nil {
				return fmt.Errorf("unable to decode address from request: %w", err)
			}
			*resp = wire.KeyAvailabilityResponse(reqAddr.Equal(&r.MockAddress))
			return nil
		case wallet.EndpointSignChannelState:
			panic("implement me")
		case wallet.EndpointVerifyChannelStateSignature:
			panic("implement me")
		default:
			return fmt.Errorf("unable to recognize endpoint: %s", endpoint)

		}
	}
}
func (m *MockRemote) CallEndpoint(endpoint string, request interface{}, response interface{}) error {
	return m.callEndpoint(endpoint, request, response)
}
