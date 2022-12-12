package wire_test

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/wallet/test"
	"perun.network/perun-cardano-backend/wire"
	pkgtest "polycry.pt/poly-go/test"
	"testing"
)

func TestMakeSigningRequest(t *testing.T) {
	rng := pkgtest.Prng(t)

	const maxMessageLength = 0x100
	message := make([]byte, maxMessageLength)
	rng.Read(message)
	addr := test.MakeRandomAddress(rng)

	expected := wire.SigningRequest{
		PubKey:  wire.MakePubKey(addr),
		Message: hex.EncodeToString(message),
	}
	actual := wire.MakeSigningRequest(addr, message)
	require.Equal(t, expected, actual, "signing request is not as expected")
}

func TestMakeVerificationRequest(t *testing.T) {
	rng := pkgtest.Prng(t)

	const maxMessageLength = 0x100
	message := make([]byte, maxMessageLength)
	rng.Read(message)
	signature := make(wallet.Sig, wire.SignatureLength)
	rng.Read(signature)
	addr := test.MakeRandomAddress(rng)

	expected := wire.VerificationRequest{
		Signature: wire.MakeSignature(signature),
		PubKey:    wire.MakePubKey(addr),
		Message:   hex.EncodeToString(message),
	}
	actual := wire.MakeVerificationRequest(signature, addr, message)
	require.Equal(t, expected, actual, "signature request is not as expected")
}

func TestMakeSignature(t *testing.T) {
	rng := pkgtest.Prng(t)
	testSig := make(wallet.Sig, wire.SignatureLength)
	rng.Read(testSig)
	expected := wire.Signature{Signature: hex.EncodeToString(testSig)}
	actual := wire.MakeSignature(testSig)
	require.Equal(t, expected, actual, "signature is not as expected")
}

func TestMakeKeyAvailabilityRequest(t *testing.T) {
	rng := pkgtest.Prng(t)
	addr := test.MakeRandomAddress(rng)
	expected := wire.MakePubKey(addr)
	actual := wire.MakeKeyAvailabilityRequest(addr)
	require.Equal(t, expected, actual, "KeyAvailabilityRequest is not as expected")
}

func TestSigningResponse_Decode(t *testing.T) {
	rng := pkgtest.Prng(t)
	expected := make(wallet.Sig, wire.SignatureLength)
	rng.Read(expected)
	uut := wire.MakeSignature(expected)
	actual, err := uut.Decode()
	require.NoError(t, err, "unexpected error when Decoding valid signature")
	require.Equal(t, expected, actual, "decoded signature is not as expected")

	const maxInvalidSigLength = 0x100
	var invalidSig wallet.Sig
	if rng.Int()%2 == 0 {
		invalidSig = make(wallet.Sig, rng.Intn(wire.SignatureLength))
	} else {
		invalidSig = make(wallet.Sig, rng.Intn(maxInvalidSigLength-wire.SignatureLength)+wire.SignatureLength+1)
	}
	rng.Read(invalidSig)
	uut = wire.MakeSignature(invalidSig)
	_, err = uut.Decode()
	require.Errorf(
		t,
		err,
		"failed to return error when decoding invalid signature of length %d from SigningResponse",
		len(invalidSig),
	)
}
