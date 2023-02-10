// Copyright 2022, 2023 - See NOTICE file for copyright holders.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
