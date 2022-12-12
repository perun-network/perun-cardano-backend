package address_test

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wallet/address"
	"perun.network/perun-cardano-backend/wallet/test"
	pkgtest "polycry.pt/poly-go/test"
	"testing"
)

func TestAddress_MarshalBinary(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := r.MockAddress
	actualBytes, err := uut.MarshalBinary()
	require.NoError(t, err, "unable to marshal valid address")
	require.Equal(
		t,
		r.MockPubKeyBytes[:],
		actualBytes,
		"wrong bytes representation of marshalled address")
}

func TestAddress_UnmarshalBinary_ValidAddressBytes(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := address.Address{}
	err := uut.UnmarshalBinary(r.MockPubKeyBytes[:])
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
	uut := address.Address{}
	err := uut.UnmarshalBinary(r.InvalidPubKeyBytes)
	require.Errorf(
		t,
		err,
		"failed to error when unmarshalling invalid public key bytes with length: %d",
		len(r.InvalidPubKeyBytes),
	)
}

func TestAddress_String(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	require.Equal(
		t,
		hex.EncodeToString(r.MockPubKeyBytes[:]),
		r.MockAddress.String(),
		"wrong string representation for public key",
	)
}

func TestAddress_Equal(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	a := &r.MockAddress
	b := &address.Address{PubKey: r.MockPubKeyBytes}
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
