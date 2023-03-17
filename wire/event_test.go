package wire_test

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wire"
	"testing"
)

const jsonCreatedSubMessage = `{
  "contents": [
    {
      "eventDatums": [
        {
          "channelParameters": {
            "pNonce": "0101000000000000000894190425c5d2ce24",
            "pPaymentPKs": [
              {
                "unPaymentPubKeyHash": {
                  "getPubKeyHash": "2d21719061b5b09a640130cc96c564ac7768de447995f0d301355c41"
                }
              },
              {
                "unPaymentPubKeyHash": {
                  "getPubKeyHash": "955a3785c9fa16e94018d80a41831f8733030a693910e7829bcb1046"
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
              "unTokenName": "\u00000xfe832508cdec04ffa3ecb5117b8b5368ca78415b5239cdf78f0036df"
            },
            "ctSymbol": {
              "unCurrencySymbol": "08fcb65e79d5517a210a2e0e50efbefe08832ff98fcf901348ad8b30"
            },
            "ctTxOutRef": {
              "txOutRefId": {
                "getTxId": "123d0277fc69704b261f073447b4a6e535449df36eb8c6bfd5c3ed8431680537"
              },
              "txOutRefIdx": 1
            }
          },
          "disputed": false,
          "funded": true,
          "funding": [
            20000000,
            20000000
          ],
          "state": {
            "balances": [
              20000000,
              20000000
            ],
            "channelId": "8eaad94121089e008b04bad9c76bed769ab0a282d19513316529d94e8c5faaee",
            "final": false,
            "version": 0
          },
          "time": 1679041735000
        }
      ],
      "eventSigs": [],
      "tag": "Created"
    }
  ],
  "tag": "NewObservableState"
}`

const jsonDisputedSubMessage = `{
  "contents": [
    {
      "eventDatums": [
        {
          "channelParameters": {
            "pNonce": "0101000000000000000894190425c5d2ce24",
            "pPaymentPKs": [
              {
                "unPaymentPubKeyHash": {
                  "getPubKeyHash": "2d21719061b5b09a640130cc96c564ac7768de447995f0d301355c41"
                }
              },
              {
                "unPaymentPubKeyHash": {
                  "getPubKeyHash": "955a3785c9fa16e94018d80a41831f8733030a693910e7829bcb1046"
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
              "unTokenName": "\u00000xfe832508cdec04ffa3ecb5117b8b5368ca78415b5239cdf78f0036df"
            },
            "ctSymbol": {
              "unCurrencySymbol": "08fcb65e79d5517a210a2e0e50efbefe08832ff98fcf901348ad8b30"
            },
            "ctTxOutRef": {
              "txOutRefId": {
                "getTxId": "123d0277fc69704b261f073447b4a6e535449df36eb8c6bfd5c3ed8431680537"
              },
              "txOutRefIdx": 1
            }
          },
          "disputed": false,
          "funded": true,
          "funding": [
            20000000,
            20000000
          ],
          "state": {
            "balances": [
              20000000,
              20000000
            ],
            "channelId": "8eaad94121089e008b04bad9c76bed769ab0a282d19513316529d94e8c5faaee",
            "final": false,
            "version": 0
          },
          "time": 1679041735000
        },
        {
          "channelParameters": {
            "pNonce": "0101000000000000000894190425c5d2ce24",
            "pPaymentPKs": [
              {
                "unPaymentPubKeyHash": {
                  "getPubKeyHash": "2d21719061b5b09a640130cc96c564ac7768de447995f0d301355c41"
                }
              },
              {
                "unPaymentPubKeyHash": {
                  "getPubKeyHash": "955a3785c9fa16e94018d80a41831f8733030a693910e7829bcb1046"
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
              "unTokenName": "\u00000xfe832508cdec04ffa3ecb5117b8b5368ca78415b5239cdf78f0036df"
            },
            "ctSymbol": {
              "unCurrencySymbol": "08fcb65e79d5517a210a2e0e50efbefe08832ff98fcf901348ad8b30"
            },
            "ctTxOutRef": {
              "txOutRefId": {
                "getTxId": "123d0277fc69704b261f073447b4a6e535449df36eb8c6bfd5c3ed8431680537"
              },
              "txOutRefIdx": 1
            }
          },
          "disputed": true,
          "funded": true,
          "funding": [
            20000000,
            20000000
          ],
          "state": {
            "balances": [
              10000000,
              30000000
            ],
            "channelId": "8eaad94121089e008b04bad9c76bed769ab0a282d19513316529d94e8c5faaee",
            "final": false,
            "version": 1
          },
          "time": 1679041783000
        }
      ],
      "eventSigs": [
        {
          "getSignature": "b3979f0983cfeb02828abe30c35b0ccc23ec6512355a3999acd2ae22c884eb9a711807e52b9678d71211abf6bc101f8e85354331c37672c896fd75ab55fdc000"
        },
        {
          "getSignature": "6e37b8900e9d3e9d1dff5c25d858d230aa70fdb2fdb05e905eba75a16ccd85cac1d5984cb58ceb489b2e16a956cc3be6554419b5e6a7aaed75b37eff4cecca09"
        }
      ],
      "tag": "Disputed"
    }
  ],
  "tag": "NewObservableState"
}`

func TestCreatedSubscriptionMessage(t *testing.T) {
	var m wire.SubscriptionMessage
	err := json.Unmarshal([]byte(jsonCreatedSubMessage), &m)
	require.NoError(t, err)
	require.Equal(t, wire.EventMessageTag, m.Tag)
	var events []wire.Event
	err = json.Unmarshal(m.Contents, &events)
	require.NoError(t, err)
	require.Equal(t, 1, len(events))
	require.Equal(t, 1, len(events[0].DatumList))
	require.Equal(t, "Created", events[0].Tag)
}

func TestDisputedSubscriptionMessage(t *testing.T) {
	var m wire.SubscriptionMessage
	err := json.Unmarshal([]byte(jsonDisputedSubMessage), &m)
	require.NoError(t, err)
	require.Equal(t, wire.EventMessageTag, m.Tag)
	var events []wire.Event
	err = json.Unmarshal(m.Contents, &events)
	require.NoError(t, err)
	require.Equal(t, 1, len(events))
	require.Equal(t, 2, len(events[0].DatumList))
	require.Equal(t, 2, len(events[0].Signatures))
	require.Equal(t, "Disputed", events[0].Tag)
}
