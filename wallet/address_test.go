package wallet_test

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wallet"
	"perun.network/perun-cardano-backend/wallet/test"
	pkgtest "polycry.pt/poly-go/test"
	"testing"
)

func TestAddress_MarshalBinary_ValidAddress(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := r.MockAddress
	actualBytes, err := uut.MarshalBinary()
	require.NoError(t, err, "unable to marshal valid address")
	require.Equal(
		t,
		r.MockAddressBytes,
		actualBytes,
		"wrong bytes representation of marshalled address")
}

func TestAddress_MarshalBinary_InvalidAddress(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := r.InvalidAddress
	actualBytes, err := uut.MarshalBinary()
	require.Errorf(
		t,
		err,
		"failed to error when marshalling invalid address: %s with length: %d",
		uut.String(),
		len(actualBytes),
	)
}

func TestAddress_UnmarshalBinary_ValidAddressBytes(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := wallet.Address{}
	err := uut.UnmarshalBinary(r.MockAddressBytes)
	require.NoError(t, err, "unable to unmarshal valid address bytes")
	require.Equal(
		t,
		r.MockAddress,
		uut,
		"marshalled address is not as expected",
	)
}

func TestAddress_UnmarshalBinary_InvalidAddressBytes(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := wallet.Address{}
	err := uut.UnmarshalBinary(r.InvalidAddressBytes)
	require.Errorf(
		t,
		err,
		"failed to error when unmarshalling invalid public key bytes with length: %d",
		len(r.InvalidAddressBytes),
	)
}

func TestAddress_String(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	require.Equal(
		t,
		r.MockAddress.PubKey,
		r.MockAddress.String(),
		"wrong string representation for public key",
	)
}

func TestAddress_Equal(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	a := &r.MockAddress
	b := &wallet.Address{PubKey: hex.EncodeToString(r.MockAddressBytes)}
	require.True(t, a.Equal(b), "addresses that have the same public key should be equal")
	require.True(t, b.Equal(a), "address equality should be commutative")
	require.False(
		t,
		r.MockAddress.Equal(&r.UnavailableAddress),
		"addresses with different public keys should not be equal",
	)
	require.False(
		t,
		r.UnavailableAddress.Equal(&r.MockAddress),
		"address equality should be commutative",
	)
	c := test.NewOtherAddressImpl(t)
	require.False(
		t,
		a.Equal(c),
		"addresses should not be equal to address of different type",
	)
}
