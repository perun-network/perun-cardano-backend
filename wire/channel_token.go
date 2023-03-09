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
