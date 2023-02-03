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

package test

import (
	"math/rand"
	"perun.network/perun-cardano-backend/channel/types"
)

func MakeRandomChannelState(rng *rand.Rand) types.ChannelState {
	channelID := MakeRandomChannelID(rng)
	balances := []uint64{rng.Uint64(), rng.Uint64()}
	version := rng.Uint64()
	final := rng.Intn(2) == 1
	return types.ChannelState{
		ID:       channelID,
		Balances: balances,
		Version:  version,
		Final:    final,
	}
}
