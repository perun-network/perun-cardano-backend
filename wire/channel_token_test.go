package wire_test

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wire"
	"testing"
)

func TestMakeChannelToken(t *testing.T) {
	const jsonString = `{
  "ctName": {
    "unTokenName": "\u00000x9ed6434876ffcb22d6fddd51c7f3c56675470d4b32f8ff0e9051ac1c"
  },
  "ctSymbol": {
    "unCurrencySymbol": "2bea49efaf89f14462c697b29471434c095316af444acf1988caeb14"
  },
  "ctTxOutRef": {
    "txOutRefId": {
      "getTxId": "deadbeef"
    },
    "txOutRefIdx": 1
  }
}`
	var ct wire.ChannelToken
	err := json.Unmarshal([]byte(jsonString), &ct)
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
