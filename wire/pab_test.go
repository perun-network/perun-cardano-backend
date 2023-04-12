package wire_test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wire"
	"testing"
)

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

func TestFundParams(t *testing.T) {
	var ct wire.ChannelToken
	err := json.Unmarshal([]byte(jsonChannelToken), &ct)
	require.NoError(t, err)
	cid := [32]byte{}
	fp := wire.MakeFundParams(cid, ct.Decode(), 1)
	res, err := json.Marshal(fp)
	fmt.Println(string(res))
}
