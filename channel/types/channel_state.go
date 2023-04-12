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

package types

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"perun.network/go-perun/channel"
)

var MaxBalance = new(big.Int).SetUint64(math.MaxUint64)

type (
	ID      = channel.ID
	Balance = uint64
	Version = uint64

	// ChannelState is the cardano backend equivalent to go-perun's channel.State.
	ChannelState struct {
		ID       ID
		Balances []Balance
		Version  Version
		Final    bool
	}
)

func MakeChannelState(id channel.ID, balances []uint64, version uint64, final bool) ChannelState {
	return ChannelState{
		ID:       id,
		Balances: balances,
		Version:  version,
		Final:    final,
	}
}

func MakeAlloc(a channel.Allocation) ([]Balance, error) {
	var err error
	if len(a.Balances) < 1 {
		return nil, fmt.Errorf("state has invalid balance")
	}
	ret := make([]uint64, len(a.Balances[0]))

	if len(a.Assets) != len(a.Balances) {
		return ret, errors.New("invalid allocation")
	}

	// Necessary because this backend currently only supports a single (native) asset and no sub-channels.
	if len(a.Assets) != 1 || len(a.Balances) != 1 || len(a.Locked) != 0 {
		return ret, fmt.Errorf("allocation incompatible with this backend")
	}
	// Necessary because this backend currently only supports a single (native) asset.
	if !a.Assets[0].Equal(Asset) {
		return ret, errors.New("allocation has asset other than native asset")
	}

	for i, balance := range a.Balances[0] {
		if ret[i], err = MakeBalance(*balance); err != nil {
			break
		}
	}
	return ret, err
}

func MakeBalance(balance big.Int) (Balance, error) {
	if balance.Sign() < 0 || balance.Cmp(MaxBalance) > 0 {
		return 0, fmt.Errorf("invalid balance")
	}
	return balance.Uint64(), nil
}

// ConvertChannelState converts a go-perun channel.State to a ChannelState.
func ConvertChannelState(state channel.State) (ChannelState, error) {
	if err := state.Valid(); err != nil {
		return ChannelState{}, fmt.Errorf("state is invalid")
	}

	balances, err := MakeAlloc(state.Allocation)
	if err != nil {
		return ChannelState{}, fmt.Errorf("unable to make allocation: %w", err)
	}
	return ChannelState{
		ID:       state.ID,
		Balances: balances,
		Version:  state.Version,
		Final:    state.IsFinal,
	}, nil
}

func (cs ChannelState) Equal(other ChannelState) bool {
	equal := cs.ID == other.ID &&
		cs.Version == other.Version &&
		cs.Final == other.Final
	if !equal {
		return false
	}
	if len(cs.Balances) != len(other.Balances) {
		return false
	}
	for i, bal := range cs.Balances {
		if bal != other.Balances[i] {
			return false
		}
	}
	return true
}
