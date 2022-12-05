package test

import (
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wallet"
	"testing"
)

func TestRemoteAccount_Address(t *testing.T) {
	seed := SetSeed()
	r := NewMockRemote()
	uut := wallet.MakeRemoteAccount(r.MockPubKey, r)
	actualAddress := uut.Address()
	require.Equalf(t, &r.MockPubKey, actualAddress, "Address returns the wrong account address, test-seed: %d", seed)
}

func TestRemoteAccount_SignData(t *testing.T) {
	seed := SetSeed()
	r := NewMockRemote()
	uut := wallet.MakeRemoteAccount(r.MockPubKey, r)
	actualSignature, err := uut.SignData(r.MockMessage)
	require.NoErrorf(t, err, "unable to sign valid data for valid address, test-seed: %d", seed)
	require.Equalf(t, r.MockSignature, actualSignature, "signature is wrong, test-seed: %d", seed)
}
