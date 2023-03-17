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

package wire

import (
	"perun.network/perun-cardano-backend/channel/types"
)

type ChannelToken struct {
	TokenName      TokenName      `json:"ctName"`
	CurrencySymbol CurrencySymbol `json:"ctSymbol"`
	CtTxOutRef     struct {
		TxOutRefId struct {
			GetTxId string `json:"getTxId"`
		} `json:"txOutRefId"`
		TxOutRefIdx int `json:"txOutRefIdx"`
	} `json:"ctTxOutRef"`
}

type TokenName struct {
	Name string `json:"unTokenName"`
}

type CurrencySymbol struct {
	Symbol string `json:"unCurrencySymbol"`
}

type AssetClass struct {
	A []interface{} `json:"unAssetClass"`
}

func MakeAssetClass(token types.ChannelToken) AssetClass {
	return AssetClass{
		A: []interface{}{CurrencySymbol{Symbol: token.TokenSymbol}, TokenName{Name: token.TokenName}},
	}
}

func (t ChannelToken) Decode() types.ChannelToken {
	return types.ChannelToken{
		TokenName:   t.TokenName.Name,
		TokenSymbol: t.CurrencySymbol.Symbol,
		TxOutRef: types.TxOutRef{
			TxID:  t.CtTxOutRef.TxOutRefId.GetTxId,
			Index: t.CtTxOutRef.TxOutRefIdx,
		},
	}
}

func MakeChannelToken(token types.ChannelToken) ChannelToken {
	return ChannelToken{
		TokenName: TokenName{
			Name: token.TokenName,
		},
		CurrencySymbol: CurrencySymbol{
			Symbol: token.TokenSymbol,
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
