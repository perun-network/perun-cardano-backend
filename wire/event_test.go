package wire_test

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wire"
	"testing"
)

const jsonSubscriptionMessage = `{
  "contents": [
    {
      "contents": [
        {
          "channelParameters": {
            "pNonce": 6794739380500570000,
            "pPaymentPKs": [
              {
                "unPaymentPubKeyHash": {
                  "getPubKeyHash": "9706069d2e482d1612cdf062d0d2f9bb3db01ab074f7c3eeb741bcd4"
                }
              },
              {
                "unPaymentPubKeyHash": {
                  "getPubKeyHash": "b50a436ae002343d30c9ddd48608a13e0e38b6785a47121c80cf45ff"
                }
              }
            ],
            "pSigningPKs": [
              {
                "unPaymentPubKey": {
                  "getPubKey": "5a3aeed83ffe0e41408a41de4cf9e1f1e39416643ea21231a2d00be46f5446a9"
                }
              },
              {
                "unPaymentPubKey": {
                  "getPubKey": "04960fbc5fe4f1ae939fdfed8a13569384474db2a38ce7b65b328d1cd578fded"
                }
              }
            ],
            "pTimeLock": 90000
          },
          "channelToken": {
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
          },
          "disputed": false,
          "funded": false,
          "funding": [
            20000000,
            0
          ],
          "state": {
            "balances": [
              20000000,
              20000000
            ],
            "channelId": "85a9eebf8644e02adef030c5ebbb5e6e9b7fbd5fa182caf57db2c8e12fc9d004",
            "final": false,
            "version": 0
          },
          "time": 1678692493000
        }
      ],
      "tag": "Created"
    }
  ],
  "tag": "NewObservableState"
}`

func TestSubscriptionMessage(t *testing.T) {
	var m wire.SubscriptionMessage
	err := json.Unmarshal([]byte(jsonSubscriptionMessage), &m)
	require.NoError(t, err)
	require.Equal(t, wire.EventMessageTag, m.Tag)
	var events []wire.Event
	err = json.Unmarshal(m.Contents, &events)
	require.NoError(t, err)
	require.Equal(t, 1, len(events))
	require.Equal(t, 1, len(events[0].DatumList))
	require.Equal(t, "Created", events[0].Tag)
}
