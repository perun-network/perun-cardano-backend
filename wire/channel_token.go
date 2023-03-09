package wire

import "perun.network/perun-cardano-backend/channel/types"

type ChannelToken struct {
	CtName struct {
		UnTokenName string `json:"unTokenName"`
	} `json:"ctName"`
	CtSymbol struct {
		UnCurrencySymbol string `json:"unCurrencySymbol"`
	} `json:"ctSymbol"`
	CtTxOutRef struct {
		TxOutRefId struct {
			GetTxId string `json:"getTxId"`
		} `json:"txOutRefId"`
		TxOutRefIdx int `json:"txOutRefIdx"`
	} `json:"ctTxOutRef"`
}

func (t ChannelToken) Decode() types.ChannelToken {
	return types.ChannelToken{
		TokenName:   t.CtName.UnTokenName,
		TokenSymbol: t.CtSymbol.UnCurrencySymbol,
		TxOutRef: types.TxOutRef{
			TxID:  t.CtTxOutRef.TxOutRefId.GetTxId,
			Index: t.CtTxOutRef.TxOutRefIdx,
		},
	}
}

func MakeChannelToken(token types.ChannelToken) ChannelToken {
	return ChannelToken{
		CtName: struct {
			UnTokenName string `json:"unTokenName"`
		}{
			UnTokenName: token.TokenName,
		},
		CtSymbol: struct {
			UnCurrencySymbol string `json:"unCurrencySymbol"`
		}{
			UnCurrencySymbol: token.TokenSymbol,
		},
		CtTxOutRef: struct {
			TxOutRefId struct {
				GetTxId string `json:"getTxId"`
			} `json:"txOutRefId"`
			TxOutRefIdx int `json:"txOutRefIdx"`
		}{
			TxOutRefId: struct {
				GetTxId string `json:"getTxId"`
			}{
				GetTxId: token.TxOutRef.TxID,
			},
			TxOutRefIdx: token.TxOutRef.Index,
		},
	}
}
