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
	uut := wallet.MakeRemoteAccount(r.MockAddress, r)
	actualAddress := uut.Address()
	require.Equal(t, &r.MockAddress, actualAddress, "Address returns the wrong account address")
}

func TestRemoteAccount_SignData(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := wallet.MakeRemoteAccount(r.MockAddress, r)
	actualSignature, err := uut.SignData(r.MockMessage)
	require.NoError(t, err, "unable to sign valid data for valid address")
	require.Equal(t, r.MockSignature, actualSignature, "signature is wrong")
}

func TestRemoteAccount_SignChannelState(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := wallet.MakeRemoteAccount(r.MockAddress, r)
	actualSignature, err := uut.SignChannelState(r.MockChannelState)
	require.NoError(t, err, "unable to sign valid channel state")
	require.Equal(t, r.MockSignature, actualSignature, "signature is wrong")
}
