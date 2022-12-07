package wallet_test

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wallet"
	"perun.network/perun-cardano-backend/wallet/test"
	pkgtest "polycry.pt/poly-go/test"
	"testing"
)

func TestPubKey_MarshalBinary_ValidPubKey(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := r.MockPubKey
	actualBytes, err := uut.MarshalBinary()
	require.NoError(t, err, "unable to marshal valid public key")
	require.Equal(
		t,
		r.MockPubKeyBytes,
		actualBytes,
		"wrong bytes representation of marshalled public key")
}

func TestPubKey_MarshalBinary_InvalidPubKey(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := r.InvalidPubKey
	actualBytes, err := uut.MarshalBinary()
	require.Errorf(
		t,
		err,
		"failed to error when marshalling invalid public key: %s with length: %d",
		uut.String(),
		len(actualBytes),
	)
}

func TestPubKey_UnmarshalBinary_ValidPubKeyBytes(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := wallet.Address{}
	err := uut.UnmarshalBinary(r.MockPubKeyBytes)
	require.NoError(t, err, "unable to unmarshal valid public key bytes")
	require.Equal(
		t,
		r.MockPubKey,
		uut,
		"marshalled public key is not as expected",
	)
}

func TestPubKey_UnmarshalBinary_InvalidPubKeyBytes(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := wallet.Address{}
	err := uut.UnmarshalBinary(r.InvalidPubKeyBytes)
	require.Errorf(
		t,
		err,
		"failed to error when unmarshalling invalid public key bytes with length: %d",
		len(r.InvalidPubKeyBytes),
	)
}

func TestPubKey_String(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	require.Equal(
		t,
		r.MockPubKey.PubKey,
		r.MockPubKey.String(),
		"wrong string representation for public key",
	)
}

func TestPubKey_Equal(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	a := &r.MockPubKey
	b := &wallet.Address{PubKey: hex.EncodeToString(r.MockPubKeyBytes)}
	require.True(t, a.Equal(b), "public keys that have the same key string should be equal")
	require.True(t, b.Equal(a), "public key equality should be commutative")
	require.False(
		t,
		r.MockPubKey.Equal(&r.UnavailablePubKey),
		"public keys with different key string should not be equal",
	)
	require.False(
		t,
		r.UnavailablePubKey.Equal(&r.MockPubKey),
		"public key equality should be commutative",
	)
	c := test.NewOtherAddressImpl(t)
	require.False(
		t,
		a.Equal(c),
		"public key should not be equal to address of different type",
	)
}
