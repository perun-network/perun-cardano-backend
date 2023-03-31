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

package types

import (
	"fmt"
	"math/big"
	"perun.network/go-perun/channel"
	"perun.network/perun-cardano-backend/wallet/address"
	"time"
)

// ChannelParameters is the cardano backend equivalent to go-perun's channel.Params.
type ChannelParameters struct {
	Parties []address.Address
	Nonce   channel.Nonce
	Timeout time.Duration
}

// MakeChannelParameters constructs a ChannelParameters from a go-perun channel.Params.
func MakeChannelParameters(params channel.Params) (ChannelParameters, error) {
	if params.App != channel.NoApp() {
		return ChannelParameters{}, fmt.Errorf("the backend does not support an app in parameters")
	}
	if params.VirtualChannel {
		return ChannelParameters{}, fmt.Errorf("the backend does not support Virtual Channels")
	}
	if !params.LedgerChannel {
		return ChannelParameters{}, fmt.Errorf("the backend only supports Ledger Channels")
	}
	parties := make([]address.Address, len(params.Parts))
	for i, party := range params.Parts {
		addr, ok := party.(*address.Address)
		if !ok {
			return ChannelParameters{}, fmt.Errorf("address %s is not of type address.Address", party.String())
		}
		parties[i] = *addr
	}
	return ChannelParameters{
		Parties: parties,
		Nonce:   new(big.Int).Set(params.Nonce),
		Timeout: time.Duration(params.ChallengeDuration) * time.Second,
	}, nil
}

func (cp ChannelParameters) Equal(other ChannelParameters) bool {
	if cp.Timeout != other.Timeout {
		return false
	}
	if cp.Nonce.Cmp(other.Nonce) != 0 {
		return false
	}
	if len(cp.Parties) != len(other.Parties) {
		return false
	}
	for i, party := range cp.Parties {
		if !party.Equal(&other.Parties[i]) {
			return false
		}
	}
	return true
}
