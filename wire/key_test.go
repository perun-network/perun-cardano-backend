package wire_test

import (
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wallet/test"
	"perun.network/perun-cardano-backend/wire"
	pkgtest "polycry.pt/poly-go/test"
	"testing"
)

func TestPubKey(t *testing.T) {
	rng := pkgtest.Prng(t)
	expected := test.MakeRandomAddress(rng)
	actual, err := wire.MakePubKey(expected).Decode()
	require.NoError(t, err, "unexpected error when decoding public key")
	require.Equal(t, expected, actual, "PubKey.Decode returned a wrong address")
}
