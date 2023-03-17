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
	"perun.network/perun-cardano-backend/wallet/test"
	"perun.network/perun-cardano-backend/wire"
	pkgtest "polycry.pt/poly-go/test"
	"testing"
)

func TestPubKey(t *testing.T) {
	rng := pkgtest.Prng(t)
	expected := test.MakeRandomAddress(rng)
	actual, err := wire.MakePubKey(expected).Decode()
	require.NoError(t, err, "unexpected error when decoding public key")
	require.Equal(t, expected.GetPubKey(), actual.GetPubKey(), "PubKey.Decode returned a wrong address")
}
