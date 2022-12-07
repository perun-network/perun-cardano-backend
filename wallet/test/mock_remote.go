package test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"perun.network/perun-cardano-backend/wallet"
	"time"
)

var testSeed *int64

func SetSeed() int64 {
	if testSeed == nil {
		testSeed = new(int64)
		*testSeed = time.Now().UnixNano()
		rand.Seed(*testSeed)
		return *testSeed
	}
	*testSeed += 1
	rand.Seed(*testSeed)
	return *testSeed
}

// MockRemote should only be instantiated using NewMockRemote.
// The default implementation has one valid signature tuple:
// (MockMessage, MockSignature, MockPubKey).
type MockRemote struct {
	MockPubKey      wallet.Address
	MockPubKeyBytes []byte
	// UnavailablePubKey is a valid wallet.Address that has associated account (private key) in this remote wallet.
	UnavailablePubKey wallet.Address
	// InvalidPubKey is invalid because it is not exactly wallet.PubKeyLength bytes long.
	InvalidPubKey      wallet.Address
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

	callSign         func(wallet.SigningRequest) (wallet.SigningResponse, error)
	callVerify       func(wallet.VerificationRequest) (wallet.VerificationResponse, error)
	callKeyAvailable func(wallet.KeyAvailabilityRequest) (wallet.KeyAvailabilityResponse, error)
}

func NewMockRemote() *MockRemote {
	r := &MockRemote{}
	initializeRandomValues(r)

	r.callSign = makeCallSignDefault(r)
	r.callVerify = makeCallVerifyDefault(r)
	r.callKeyAvailable = makeCallKeyAvailableDefault(r)
	return r
}

func initializeRandomValues(r *MockRemote) {
	const maxMessageLength = 2 ^ 8 // in bytes
	const maxInvalidPubKeyLength = wallet.PubKeyLength * 2
	const maxInvalidSignatureLength = wallet.SignatureLength * 2

	r.MockPubKeyBytes = make([]byte, wallet.PubKeyLength)
	rand.Read(r.MockPubKeyBytes)
	r.MockPubKey = wallet.Address{PubKey: hex.EncodeToString(r.MockPubKeyBytes)}

	unavailablePubKeyBytes := make([]byte, wallet.PubKeyLength)
	for bytes.Equal(r.MockPubKeyBytes, unavailablePubKeyBytes) {
		rand.Read(unavailablePubKeyBytes)
	}
	r.UnavailablePubKey = wallet.Address{PubKey: hex.EncodeToString(unavailablePubKeyBytes)}

	if rand.Int()%2 == 0 {
		r.InvalidPubKeyBytes = make([]byte, rand.Intn(wallet.PubKeyLength))
	} else {
		r.InvalidPubKeyBytes = make([]byte, rand.Intn(maxInvalidPubKeyLength-wallet.PubKeyLength)+wallet.PubKeyLength+1)
	}
	rand.Read(r.InvalidPubKeyBytes)
	r.InvalidPubKey = wallet.Address{PubKey: hex.EncodeToString(r.InvalidPubKeyBytes)}

	r.MockSignature = make([]byte, wallet.SignatureLength)
	rand.Read(r.MockSignature)
	r.MockSignatureString = hex.EncodeToString(r.MockSignature)

	r.OtherSignature = make([]byte, wallet.SignatureLength)
	for bytes.Equal(r.MockSignature, r.OtherSignature) {
		rand.Read(r.OtherSignature)
	}
	r.OtherSignatureString = hex.EncodeToString(r.OtherSignature)

	r.InvalidSignatureShorter = make([]byte, rand.Intn(wallet.SignatureLength))
	rand.Read(r.InvalidSignatureShorter)

	r.InvalidSignatureLonger = make([]byte, rand.Intn(maxInvalidSignatureLength-wallet.SignatureLength)+wallet.SignatureLength+1)
	rand.Read(r.InvalidSignatureLonger)

	r.MockMessage = make([]byte, rand.Intn(maxMessageLength+1))
	rand.Read(r.MockMessage)
	r.MockMessageString = hex.EncodeToString(r.MockMessage)
}

func (m *MockRemote) SetCallSign(f func(request wallet.SigningRequest) (wallet.SigningResponse, error)) {
	m.callSign = f
}

func makeCallSignDefault(r *MockRemote) func(request wallet.SigningRequest) (wallet.SigningResponse, error) {
	return func(request wallet.SigningRequest) (wallet.SigningResponse, error) {
		if !request.AccountPubKey.Equal(&r.MockPubKey) {
			return wallet.SigningResponse{}, fmt.Errorf("invalid public key for mock remote")
		}

		if request.Message != r.MockMessageString {
			return wallet.SigningResponse{}, fmt.Errorf("invalid data for mock remote")
		}
		return wallet.SignatureWrapper{Signature: r.MockSignatureString}, nil
	}
}

func makeCallVerifyDefault(r *MockRemote) func(wallet.VerificationRequest) (wallet.VerificationResponse, error) {
	return func(request wallet.VerificationRequest) (wallet.VerificationResponse, error) {
		if !request.Address.Equal(&r.MockPubKey) && !request.Address.Equal(&r.UnavailablePubKey) {
			return false, fmt.Errorf("invalid public key for mock remote")
		}
		if request.Address.Equal(&r.UnavailablePubKey) {
			return false, nil
		}

		if request.Message != r.MockMessageString {
			return false, fmt.Errorf("invalid data for mock remote")
		}
		if request.SigWrapper.Signature == r.MockSignatureString {
			return true, nil
		}
		if request.SigWrapper.Signature == r.OtherSignatureString {
			return false, nil
		}
		if request.SigWrapper.Signature == hex.EncodeToString(r.InvalidSignatureShorter) ||
			request.SigWrapper.Signature == hex.EncodeToString(r.InvalidSignatureLonger) {
			panic("mock remote received signature of invalid length to verify")
		}
		return false, fmt.Errorf("invalid signature for mock remote")
	}
}

func makeCallKeyAvailableDefault(r *MockRemote) func(wallet.KeyAvailabilityRequest) (wallet.KeyAvailabilityResponse, error) {
	return func(request wallet.KeyAvailabilityRequest) (wallet.KeyAvailabilityResponse, error) {
		_, err := request.MarshalBinary()
		if err != nil {
			return false, fmt.Errorf("invalid pubKey: %w", err)
		}
		return request.Equal(&r.MockPubKey), nil
	}
}

func (m *MockRemote) CallSign(request wallet.SigningRequest) (wallet.SigningResponse, error) {
	return m.callSign(request)
}

func (m *MockRemote) CallVerify(request wallet.VerificationRequest) (wallet.VerificationResponse, error) {
	return m.callVerify(request)
}

func (m *MockRemote) CallKeyAvailable(request wallet.KeyAvailabilityRequest) (wallet.KeyAvailabilityResponse, error) {
	return m.callKeyAvailable(request)
}
