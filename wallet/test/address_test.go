package test

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	pw "perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/wallet"
	"testing"
)

type NotPubKeyAddress struct {
	T    *testing.T
	Seed int64
}

func (n *NotPubKeyAddress) MarshalBinary() (data []byte, err error) {
	require.Failf(
		n.T,
		"failure",
		"MarshalBinary() should not be called on a NotPubKeyAddress, test-seed: %d",
		n.Seed,
	)
	return nil, nil
}

func (n *NotPubKeyAddress) UnmarshalBinary(data []byte) error {
	require.Failf(
		n.T,
		"failure",
		"UnmarshalBinary() should not be called on a NotPubKeyAddress, test-seed: %d",
		n.Seed,
	)
	return nil
}

func (n *NotPubKeyAddress) String() string {
	require.Failf(
		n.T,
		"failure",
		"String() should not be called on a NotPubKeyAddress, test-seed: %d",
		n.Seed,
	)
	return ""
}

func (n *NotPubKeyAddress) Equal(address pw.Address) bool {
	require.Failf(
		n.T,
		"failure",
		"Equal() should not be called on a NotPubKeyAddress, test-seed: %d",
		n.Seed,
	)
	return false
}

var _ pw.Address = &NotPubKeyAddress{}

func TestPubKey_MarshalBinary_ValidPubKey(t *testing.T) {
	seed := SetSeed()
	r := NewMockRemote()
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
	seed := SetSeed()
	r := NewMockRemote()
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
	seed := SetSeed()
	r := NewMockRemote()
	uut := wallet.PubKey{}
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
	seed := SetSeed()
	r := NewMockRemote()
	uut := wallet.PubKey{}
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
	seed := SetSeed()
	r := NewMockRemote()
	require.Equalf(
		t,
		r.MockPubKey.Key,
		r.MockPubKey.String(),
		"wrong string representation for public key, test-seed: %d",
		seed,
	)
}

func TestPubKey_Equal(t *testing.T) {
	seed := SetSeed()
	r := NewMockRemote()
	a := &r.MockPubKey
	b := &wallet.PubKey{Key: hex.EncodeToString(r.MockPubKeyBytes)}
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
	c := &NotPubKeyAddress{
		T:    t,
		Seed: seed,
	}
	require.Falsef(
		t,
		a.Equal(c),
		"public key should not be equal to address of different type, test-seed: %d",
		seed,
	)
}
