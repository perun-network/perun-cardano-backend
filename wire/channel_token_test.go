// Copyright 2023 - See NOTICE file for copyright holders.
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
	"encoding/json"
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wire"
	"testing"
)

func TestMakeChannelToken(t *testing.T) {
	const jsonChannelToken = `{
  "ctName": {
    "unTokenName": "\u00000x1f726fdf149afa045f6d2dea978cf4d6d82304d2d68e4263fcd9a518"
  },
  "ctSymbol": {
    "unCurrencySymbol": "5783a64780a2aa5a14e1824999713727087a2c2eb423c7080475570d"
  },
  "ctTxOutRef": {
    "txOutRefId": {
      "getTxId": "a707536d289e9eb51f251ffb193704f8c1c99692148bcc992616b1b11d364c41"
    },
    "txOutRefIdx": 2
  }
}`
	var ct wire.ChannelToken
	err := json.Unmarshal([]byte(jsonChannelToken), &ct)
	require.NoError(t, err)
	tCt := ct.Decode()
	wCt := wire.MakeChannelToken(tCt)
	res, err := json.Marshal(wCt)
	require.NoError(t, err)
	var cmp wire.ChannelToken
	err = json.Unmarshal(res, &cmp)
	require.NoError(t, err)
	require.Equal(t, ct, cmp)
}
