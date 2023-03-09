package wire_test

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wire"
	"testing"
)

func TestChannelDatum_Decode(t *testing.T) {
	const jsonString = `
{
  "channelParameters": {
    "pNonce": 69,
    "pPaymentPKs": [
      {
        "unPaymentPubKeyHash": {
          "getPubKeyHash": "a2c20c77887ace1cd986193e4e75babd8993cfd56995cd5cfce609c2"
        }
      },
      {
        "unPaymentPubKeyHash": {
          "getPubKeyHash": "80a4f45b56b88d1139da23bc4c3c75ec6d32943c087f250b86193ca7"
        }
      }
    ],
    "pSigningPKs": [
      {
        "unPaymentPubKey": {
          "getPubKey": "8d9de88fbf445b7f6c3875a14daba94caee2ffcbc9ac211c95aba0a2f5711853"
        }
      },
      {
        "unPaymentPubKey": {
          "getPubKey": "98c77c40ccc536e0d433874dae97d4a0787b10b3bca0dc2e1bdc7be0a544f0ac"
        }
      }
    ],
    "pTimeLock": 15
  },
  "channelToken": {
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
  },
  "disputed": false,
  "funded": true,
  "funding": [
    7,
    8
  ],
  "state": {
    "balances": [
      1,
      2
    ],
    "channelId": "ea0d44056537e06dd7f38c94b099f7556072a163ad16f40204bc23a4c2e20c53",
    "final": false,
    "version": 1337
  },
  "time": 1000
}
`
	var channelDatum wire.ChannelDatum
	err := json.Unmarshal([]byte(jsonString), &channelDatum)
	require.NoError(t, err)
	tChannelDatum, err := channelDatum.Decode()
	require.NoError(t, err)
	wChannelDatum := wire.MakeChannelDatum(tChannelDatum)
	require.Equal(t, channelDatum, wChannelDatum)
	res, err := json.Marshal(wChannelDatum)
	require.NoError(t, err)
	var cmp wire.ChannelDatum
	err = json.Unmarshal(res, &cmp)
	require.NoError(t, err)
	require.Equal(t, channelDatum, cmp)
}
