package wire_test

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wallet/address"
	"perun.network/perun-cardano-backend/wallet/test"
	"perun.network/perun-cardano-backend/wire"
	pkgtest "polycry.pt/poly-go/test"
	"testing"
)

func TestMakePubKey(t *testing.T) {
	rng := pkgtest.Prng(t)
	testAddr := test.MakeRandomAddress(rng)

	expected := wire.PubKey{PubKey: hex.EncodeToString(testAddr.GetPubKeySlice())}
	actual := wire.MakePubKey(testAddr)
	require.Equal(t, expected, actual, "PubKey not as expected")
}

func TestPubKey_Decode(t *testing.T) {
	const maxInvalidPubKeyLength = 128

	rng := pkgtest.Prng(t)
	expected := test.MakeRandomAddress(rng)
	uut := wire.PubKey{PubKey: hex.EncodeToString(expected.GetPubKeySlice())}
	actual, err := uut.Decode()
	require.NoError(t, err, "unexpected error when decoding public key")
	require.Equal(t, expected, actual, "decoded address is wrong")

	var invalidPubKeyBytes []byte
	if rng.Int()%2 == 0 {
		invalidPubKeyBytes = make([]byte, rng.Intn(address.PubKeyLength))
	} else {
		invalidPubKeyBytes = make([]byte, rng.Intn(maxInvalidPubKeyLength-address.PubKeyLength)+address.PubKeyLength+1)
	}
	uut = wire.PubKey{PubKey: hex.EncodeToString(invalidPubKeyBytes)}
	_, err = uut.Decode()
	require.Errorf(t, err, "failed to return error when decoding PubKey with invalid pubKeyLength %d", len(invalidPubKeyBytes))
}
func TestPubKey(t *testing.T) {
	rng := pkgtest.Prng(t)
	expected := test.MakeRandomAddress(rng)
	actual, err := wire.MakePubKey(expected).Decode()
	require.NoError(t, err, "unexpected error when decoding public key")
	require.Equal(t, expected, actual, "PubKey.Decode returned a wrong address")
}
