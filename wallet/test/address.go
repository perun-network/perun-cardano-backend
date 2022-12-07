package test

import (
	"github.com/stretchr/testify/require"
	pw "perun.network/go-perun/wallet"
	"testing"
)

type OtherAddressImpl struct {
	t *testing.T
}

func NewOtherAddressImpl(t *testing.T) *OtherAddressImpl {
	return &OtherAddressImpl{
		t: t,
	}
}

func (n *OtherAddressImpl) MarshalBinary() (data []byte, err error) {
	require.Fail(
		n.t,
		"failure",
		"MarshalBinary() should not be called on a OtherAddressImpl",
	)
	return nil, nil
}

func (n *OtherAddressImpl) UnmarshalBinary(data []byte) error {
	require.Fail(
		n.t,
		"failure",
		"UnmarshalBinary() should not be called on a OtherAddressImpl",
	)
	return nil
}

func (n *OtherAddressImpl) String() string {
	require.Fail(
		n.t,
		"failure",
		"String() should not be called on a OtherAddressImpl",
	)
	return ""
}

func (n *OtherAddressImpl) Equal(address pw.Address) bool {
	require.Fail(
		n.t,
		"failure",
		"Equal() should not be called on a OtherAddressImpl",
	)
	return false
}

var _ pw.Address = &OtherAddressImpl{}
