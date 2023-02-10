// Copyright 2022, 2023 - See NOTICE file for copyright holders.
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
	"perun.network/perun-cardano-backend/wallet/address"
)

// PubKey is a json serializable public key to communicate with cardano apis (see: Ledger.Crypto.PubKey).
type PubKey struct {
	Hex string `json:"getPubKey"`
}

// MakePubKey returns a PubKey
func MakePubKey(address address.Address) PubKey {
	return PubKey{
		Hex: hex.EncodeToString(address.GetPubKeySlice()),
	}
}

func (key PubKey) Decode() (address.Address, error) {
	pubKey, err := hex.DecodeString(key.Hex)
	if err != nil {
		return address.Address{}, fmt.Errorf("unable to decode PubKey hex string: %w", err)
	}
	return address.MakeAddressFromByteSlice(pubKey)
}
