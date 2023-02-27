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

package wire_test

import (
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/channel/test"
	"perun.network/perun-cardano-backend/wire"
	pkgtest "polycry.pt/poly-go/test"
	"testing"
)

func TestMakeChannelState(t *testing.T) {
	rng := pkgtest.Prng(t)
	testChannelState := test.MakeRandomChannelState(rng)

	expected := wire.ChannelState{
		Balances:  testChannelState.Balances,
		ChannelID: testChannelState.ID,
		Final:     testChannelState.Final,
		Version:   testChannelState.Version,
	}
	actual := wire.MakeChannelState(testChannelState)
	require.Equal(t, expected, actual, "wire.ChannelState not as expected")
}

func TestChannelState_Decode(t *testing.T) {
	rng := pkgtest.Prng(t)
	expected := test.MakeRandomChannelState(rng)
	uut := wire.ChannelState{
		Balances:  expected.Balances,
		ChannelID: expected.ID,
		Final:     expected.Final,
		Version:   expected.Version,
	}
	actual := uut.Decode()
	require.Equal(t, expected, actual, "decoded channel state is wrong")
}
func TestChannelState(t *testing.T) {
	rng := pkgtest.Prng(t)
	expected := test.MakeRandomChannelState(rng)
	actual := wire.MakeChannelState(expected).Decode()
	require.Equal(t, expected, actual, "channel state not as expected")
}
