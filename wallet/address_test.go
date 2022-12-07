package wallet_test

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wallet"
	"perun.network/perun-cardano-backend/wallet/test"
	"testing"
)

func TestPubKey_MarshalBinary_ValidPubKey(t *testing.T) {
	seed := test.SetSeed()
	r := test.NewMockRemote()
	uut := r.MockPubKey
	actualBytes, err := uut.MarshalBinary()
	require.NoErrorf(t, err, "unable to marshal valid public key, test-seed: %d", seed)
	require.Equalf(
		t,
		r.MockPubKeyBytes,
		actualBytes,
		"wrong bytes representation of marshalled public key, test-seed: %d",
		seed,
	)
}

func TestPubKey_MarshalBinary_InvalidPubKey(t *testing.T) {
	seed := test.SetSeed()
	r := test.NewMockRemote()
	uut := r.InvalidPubKey
	actualBytes, err := uut.MarshalBinary()
	require.Errorf(
		t,
		err,
		"failed to error when marshalling invalid public key: %s with length: %d, test-seed: %d",
		uut.String(),
		len(actualBytes),
		seed,
	)
}

func TestPubKey_UnmarshalBinary_ValidPubKeyBytes(t *testing.T) {
	seed := test.SetSeed()
	r := test.NewMockRemote()
	uut := wallet.Address{}
	err := uut.UnmarshalBinary(r.MockPubKeyBytes)
	require.NoErrorf(t, err, "unable to unmarshal valid public key bytes, test-seed: %d", seed)
	require.Equalf(
		t,
		r.MockPubKey,
		uut,
		"marshalled public key is not as expected, test-seed: %d",
		seed,
	)
}

func TestPubKey_UnmarshalBinary_InvalidPubKeyBytes(t *testing.T) {
	seed := test.SetSeed()
	r := test.NewMockRemote()
	uut := wallet.Address{}
	err := uut.UnmarshalBinary(r.InvalidPubKeyBytes)
	require.Errorf(
		t,
		err,
		"failed to error when unmarshalling invalid public key bytes with length: %d, test-seed: %d",
		len(r.InvalidPubKeyBytes),
		seed,
	)
}

func TestPubKey_String(t *testing.T) {
	seed := test.SetSeed()
	r := test.NewMockRemote()
	require.Equalf(
		t,
		r.MockPubKey.PubKey,
		r.MockPubKey.String(),
		"wrong string representation for public key, test-seed: %d",
		seed,
	)
}

func TestPubKey_Equal(t *testing.T) {
	seed := test.SetSeed()
	r := test.NewMockRemote()
	a := &r.MockPubKey
	b := &wallet.Address{PubKey: hex.EncodeToString(r.MockPubKeyBytes)}
	require.Truef(t, a.Equal(b), "public keys that have the same key string should be equal, test-seed: %d", seed)
	require.Truef(t, b.Equal(a), "public key equality should be commutative, test-seed: %d", seed)
	require.Falsef(
		t,
		r.MockPubKey.Equal(&r.UnavailablePubKey),
		"public keys with different key string should not be equal, test-seed: %d",
		seed,
	)
	require.Falsef(
		t,
		r.UnavailablePubKey.Equal(&r.MockPubKey),
		"public key equality should be commutative, test-seed: %d",
		seed,
	)
	c := test.NewOtherAddressImpl(t, seed)
	require.Falsef(
		t,
		a.Equal(c),
		"public key should not be equal to address of different type, test-seed: %d",
		seed,
	)
}
