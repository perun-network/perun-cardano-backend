package test

import (
	"github.com/stretchr/testify/require"
	pw "perun.network/go-perun/wallet"
	"testing"
)

type OtherAddressImpl struct {
	t    *testing.T
	seed int64
}

func NewOtherAddressImpl(t *testing.T, seed int64) *OtherAddressImpl {
	return &OtherAddressImpl{
		t:    t,
		seed: seed,
	}
}

func (n *OtherAddressImpl) MarshalBinary() (data []byte, err error) {
	require.Failf(
		n.t,
		"failure",
		"MarshalBinary() should not be called on a OtherAddressImpl, test-seed: %d",
		n.seed,
	)
	return nil, nil
}

func (n *OtherAddressImpl) UnmarshalBinary(data []byte) error {
	require.Failf(
		n.t,
		"failure",
		"UnmarshalBinary() should not be called on a OtherAddressImpl, test-seed: %d",
		n.seed,
	)
	return nil
}

func (n *OtherAddressImpl) String() string {
	require.Failf(
		n.t,
		"failure",
		"String() should not be called on a OtherAddressImpl, test-seed: %d",
		n.seed,
	)
	return ""
}

func (n *OtherAddressImpl) Equal(address pw.Address) bool {
	require.Failf(
		n.t,
		"failure",
		"Equal() should not be called on a OtherAddressImpl, test-seed: %d",
		n.seed,
	)
	return false
}

var _ pw.Address = &OtherAddressImpl{}
