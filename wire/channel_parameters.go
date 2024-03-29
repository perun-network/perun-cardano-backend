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
	"encoding/hex"
	"fmt"
	"math/big"
	"perun.network/perun-cardano-backend/channel/types"
	"perun.network/perun-cardano-backend/wallet/address"
	"time"
)

// ChannelParameters reflects the Haskell type `Channel` of the Channel Smart Contract in respect to its json encoding.
type ChannelParameters struct {
	Nonce               string              `json:"pNonce"`
	PaymentPubKeyHashes []PaymentPubKeyHash `json:"pPaymentPKs"`
	SigningPubKeys      []PaymentPubKey     `json:"pSigningPKs"`
	TimeLock            int64               `json:"pTimeLock"`
}

type PaymentPubKey struct {
	PubKey PubKey `json:"unPaymentPubKey"`
}

type PaymentPubKeyHash struct {
	PubKeyHash PubKeyHash `json:"unPaymentPubKeyHash"`
}

type PubKeyHash struct {
	Hex string `json:"getPubKeyHash"`
}

func MakePaymentPubKey(address address.Address) PaymentPubKey {
	return PaymentPubKey{
		PubKey: MakePubKey(address),
	}
}

func MakeChannelParameters(parameters types.ChannelParameters) ChannelParameters {
	pubKeyHashes := make([]PaymentPubKeyHash, len(parameters.Parties))
	signingPubKeys := make([]PaymentPubKey, len(parameters.Parties))
	for i, addr := range parameters.Parties {
		pubKeyHashes[i] = MakePaymentPubKeyHash(addr)
		signingPubKeys[i] = MakePaymentPubKey(addr)
	}
	return ChannelParameters{
		Nonce:               fmt.Sprintf("%x", parameters.Nonce),
		PaymentPubKeyHashes: pubKeyHashes,
		SigningPubKeys:      signingPubKeys,
		TimeLock:            parameters.Timeout.Milliseconds(),
	}

}

func MakePaymentPubKeyHash(address address.Address) PaymentPubKeyHash {
	hash := address.GetPubKeyHash()
	return PaymentPubKeyHash{
		PubKeyHash: PubKeyHash{
			Hex: hex.EncodeToString(hash[:]),
		},
	}
}

func (pkh PubKeyHash) Decode() ([]byte, error) {
	bytes, err := hex.DecodeString(pkh.Hex)
	if err != nil {
		return nil, err
	}
	if len(bytes) != address.PubKeyHashLength {
		return nil, fmt.Errorf("pubKeyHash has wrong length: %d", len(bytes))
	}
	return bytes, nil
}

func (cp ChannelParameters) Decode() (types.ChannelParameters, error) {
	parties := make([]address.Address, len(cp.SigningPubKeys))
	n, ok := new(big.Int).SetString(cp.Nonce, 16)
	if !ok {
		return types.ChannelParameters{}, fmt.Errorf("unable to decode nonce: %s", cp.Nonce)
	}
	for i, ppk := range cp.SigningPubKeys {
		addr, err := ppk.PubKey.Decode()
		if err != nil {
			return types.ChannelParameters{}, err
		}
		pkh, err := cp.PaymentPubKeyHashes[i].PubKeyHash.Decode()
		if err != nil {
			return types.ChannelParameters{}, err
		}
		err = addr.SetPaymentPubKeyHashFromSlice(pkh)
		if err != nil {
			return types.ChannelParameters{}, err
		}
		parties[i] = addr
	}
	return types.ChannelParameters{
		Parties: parties,
		Nonce:   n,
		Timeout: time.Duration(cp.TimeLock) * time.Millisecond,
	}, nil
}
