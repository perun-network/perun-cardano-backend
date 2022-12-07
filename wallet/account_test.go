package wallet_test

import (
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wallet"
	"perun.network/perun-cardano-backend/wallet/test"
	pkgtest "polycry.pt/poly-go/test"
	"testing"
)

func TestRemoteAccount_Address(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := wallet.MakeRemoteAccount(r.MockPubKey, r)
	actualAddress := uut.Address()
	require.Equal(t, &r.MockPubKey, actualAddress, "Address returns the wrong account address")
}

func TestRemoteAccount_SignData(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := wallet.MakeRemoteAccount(r.MockPubKey, r)
	actualSignature, err := uut.SignData(r.MockMessage)
	require.NoError(t, err, "unable to sign valid data for valid address")
	require.Equal(t, r.MockSignature, actualSignature, "signature is wrong")
}
