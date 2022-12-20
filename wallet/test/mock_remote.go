package test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
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

	callSign         func(wire.SigningRequest) (wire.SigningResponse, error)
	callVerify       func(wire.VerificationRequest) (wire.VerificationResponse, error)
	callKeyAvailable func(wire.KeyAvailabilityRequest) (wire.KeyAvailabilityResponse, error)
}

func NewMockRemote(rng *rand.Rand) *MockRemote {
	r := &MockRemote{}
	initializeRandomValues(r, rng)

	r.callSign = makeCallSignDefault(r)
	r.callVerify = makeCallVerifyDefault(r)
	r.callKeyAvailable = makeCallKeyAvailableDefault(r)
	return r
}

func initializeRandomValues(r *MockRemote, rng *rand.Rand) {
	const maxMessageLength = 0x100 // in bytes

	r.MockAddress = MakeRandomAddress(rng)
	r.MockPubKeyBytes = r.MockAddress.GetPubKey()

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

func (m *MockRemote) SetCallSign(f func(request wire.SigningRequest) (wire.SigningResponse, error)) {
	m.callSign = f
}

func makeCallSignDefault(r *MockRemote) func(request wire.SigningRequest) (wire.SigningResponse, error) {
	return func(request wire.SigningRequest) (wire.SigningResponse, error) {
		reqAddr, err := request.PubKey.Decode()
		if err != nil {
			return wire.SigningResponse{}, fmt.Errorf("unable to decode PubKey from request")
		}
		if !reqAddr.Equal(&r.MockAddress) {
			return wire.SigningResponse{}, fmt.Errorf("invalid public key for mock remote")
		}

		if request.Message != r.MockMessageString {
			return wire.SigningResponse{}, fmt.Errorf("invalid data for mock remote")
		}
		return wire.Signature{Hex: r.MockSignatureString}, nil
	}
}

func makeCallVerifyDefault(r *MockRemote) func(wire.VerificationRequest) (wire.VerificationResponse, error) {
	return func(request wire.VerificationRequest) (wire.VerificationResponse, error) {
		reqAddr, err := request.PubKey.Decode()
		if err != nil {
			return false, fmt.Errorf("unable to decode PubKey from request")
		}
		if !reqAddr.Equal(&r.MockAddress) && !reqAddr.Equal(&r.UnavailableAddress) {
			return false, fmt.Errorf("invalid public key for mock remote")
		}
		if reqAddr.Equal(&r.UnavailableAddress) {
			return false, nil
		}

		if request.Message != r.MockMessageString {
			return false, fmt.Errorf("invalid data for mock remote")
		}
		if request.Signature.Hex == r.MockSignatureString {
			return true, nil
		}
		if request.Signature.Hex == r.OtherSignatureString {
			return false, nil
		}
		if request.Signature.Hex == hex.EncodeToString(r.InvalidSignatureShorter) ||
			request.Signature.Hex == hex.EncodeToString(r.InvalidSignatureLonger) {
			panic("mock remote received signature of invalid length to verify")
		}
		return false, fmt.Errorf("invalid signature for mock remote")
	}
}

func makeCallKeyAvailableDefault(r *MockRemote) func(wire.KeyAvailabilityRequest) (wire.KeyAvailabilityResponse, error) {
	return func(request wire.KeyAvailabilityRequest) (wire.KeyAvailabilityResponse, error) {
		reqAddr, err := request.Decode()
		if err != nil {
			return false, fmt.Errorf("unable to decode address from request: %w", err)
		}
		return reqAddr.Equal(&r.MockAddress), nil
	}
}

func (m *MockRemote) CallSign(request wire.SigningRequest) (wire.SigningResponse, error) {
	return m.callSign(request)
}

func (m *MockRemote) CallVerify(request wire.VerificationRequest) (wire.VerificationResponse, error) {
	return m.callVerify(request)
}

func (m *MockRemote) CallKeyAvailable(request wire.KeyAvailabilityRequest) (wire.KeyAvailabilityResponse, error) {
	return m.callKeyAvailable(request)
}
