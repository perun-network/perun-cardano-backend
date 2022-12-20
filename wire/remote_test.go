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

func TestMake(t *testing.T) {
	rng := pkgtest.Prng(t)

	const maxMessageLength = 0x100

	randomSigningRequestTest := func() func(*testing.T) {
		referenceAddress := test.MakeRandomAddress(rng)
		referenceMessage := make([]byte, rng.Intn(maxMessageLength))
		rng.Read(referenceMessage)

		return func(t *testing.T) {
			t.Parallel()
			expected := wire.SigningRequest{
				PubKey:  wire.MakePubKey(referenceAddress),
				Message: hex.EncodeToString(referenceMessage),
			}
			actual := wire.MakeSigningRequest(referenceAddress, referenceMessage)
			require.Equal(t, expected, actual, "signing request is not as expected")
		}
	}
	randomVerificationRequestTest := func() func(*testing.T) {
		referenceAddress := test.MakeRandomAddress(rng)
		referenceSignature := make(wallet.Sig, wire.SignatureLength)
		rng.Read(referenceSignature)
		referenceMessage := make([]byte, rng.Intn(maxMessageLength))
		rng.Read(referenceMessage)

		return func(t *testing.T) {
			t.Parallel()
			expected := wire.VerificationRequest{
				Signature: wire.MakeSignature(referenceSignature),
				PubKey:    wire.MakePubKey(referenceAddress),
				Message:   hex.EncodeToString(referenceMessage),
			}
			actual := wire.MakeVerificationRequest(referenceSignature, referenceAddress, referenceMessage)
			require.Equal(t, expected, actual, "signature request is not as expected")
		}
	}

	randomKeyAvailabilityRequestTest := func() func(*testing.T) {
		referenceAddress := test.MakeRandomAddress(rng)

		return func(t *testing.T) {
			t.Parallel()
			expected := wire.MakePubKey(referenceAddress)
			actual := wire.MakeKeyAvailabilityRequest(referenceAddress)
			require.Equal(t, expected, actual, "KeyAvailabilityRequest is not as expected")
		}
	}

	randomSignatureTest := func() func(*testing.T) {
		referenceSignature := make(wallet.Sig, wire.SignatureLength)
		rng.Read(referenceSignature)

		return func(t *testing.T) {
			t.Parallel()
			expected := wire.Signature{Hex: hex.EncodeToString(referenceSignature)}
			actual := wire.MakeSignature(referenceSignature)
			require.Equal(t, expected, actual, "signature is not as expected")
		}
	}

	randomPubKeyTest := func() func(*testing.T) {
		referenceAddress := test.MakeRandomAddress(rng)

		return func(t *testing.T) {
			t.Parallel()
			expected := wire.PubKey{Hex: hex.EncodeToString(referenceAddress.GetPubKeySlice())}
			actual := wire.MakePubKey(referenceAddress)
			require.Equal(t, expected, actual, "PubKey not as expected")
		}
	}

	for i := 0; i < 100; i++ {
		t.Run("MakeSigningRequest", randomSigningRequestTest())
		t.Run("MakeVerificationRequest", randomVerificationRequestTest())
		t.Run("MakeKeyAvailabilityRequest", randomKeyAvailabilityRequestTest())
		t.Run("MakeSignature", randomSignatureTest())
		t.Run("MakePubKey", randomPubKeyTest())
	}
}

func TestDecode(t *testing.T) {
	rng := pkgtest.Prng(t)

	randomSigningResponseTest := func() func(*testing.T) {
		referenceSignature := make(wallet.Sig, wire.SignatureLength)
		rng.Read(referenceSignature)

		return func(t *testing.T) {
			t.Parallel()
			uut := wire.MakeSignature(referenceSignature)
			actual, err := uut.Decode()
			require.NoError(t, err, "unexpected error when Decoding valid signature")
			require.Equal(t, referenceSignature, actual, "decoded signature is not as expected")
		}
	}
	randomSigningResponseInvalidTest := func(invalidSig wallet.Sig) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			uut := wire.MakeSignature(invalidSig)
			_, err := uut.Decode()
			require.Errorf(
				t,
				err,
				"failed to return error when decoding invalid signature of length %d from SigningResponse",
				len(invalidSig),
			)
		}
	}

	randomPubKeyTest := func() func(*testing.T) {
		referenceAddress := test.MakeRandomAddress(rng)

		return func(t *testing.T) {
			t.Parallel()
			uut := wire.PubKey{Hex: hex.EncodeToString(referenceAddress.GetPubKeySlice())}
			actual, err := uut.Decode()
			require.NoError(t, err, "unexpected error when decoding public key")
			require.Equal(t, referenceAddress, actual, "decoded address is wrong")
		}
	}

	randomPubKeyInvalidTest := func(invalidPubKeyBytes []byte) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			uut := wire.PubKey{Hex: hex.EncodeToString(invalidPubKeyBytes)}
			_, err := uut.Decode()
			require.Errorf(
				t,
				err,
				"failed to return error when decoding PubKey with invalid pubKeyLength %d",
				len(invalidPubKeyBytes),
			)
		}
	}

	for i := 0; i < 100; i++ {
		t.Run("SigningResponse Decode - Valid", randomSigningResponseTest())

		t.Run(
			"SigningResponse Decode - Invalid - Signature too short",
			randomSigningResponseInvalidTest(test.MakeTooShortSignature(rng)),
		)
		t.Run(
			"SigningResponse Decode - Invalid - Signature too long",
			randomSigningResponseInvalidTest(test.MakeTooLongSignature(rng)))
		t.Run("PubKey Decode - Valid", randomPubKeyTest())
		t.Run(
			"PubKey Decode - Invalid - PubKey too short",
			randomPubKeyInvalidTest(test.MakeTooFewPublicKeyBytes(rng)),
		)
		t.Run(
			"PubKey Decode - Invalid - PubKey too long",
			randomPubKeyInvalidTest(test.MakeTooManyPublicKeyBytes(rng)),
		)
	}
}
