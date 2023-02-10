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
	"encoding/json"
	"fmt"
	"perun.network/perun-cardano-backend/channel/types"
)

const ChannelIDLength = 32

// ChannelState reflects the Haskell type `ChannelState` of the Channel Smart Contract in respect to its json encoding.
type ChannelState struct {
	Balances  []uint64  `json:"balances"`
	ChannelID ChannelID `json:"channelId"`
	Final     bool      `json:"final"`
	Version   uint64    `json:"version"`
}

func MakeChannelState(cs types.ChannelState) ChannelState {
	return ChannelState{
		Balances:  cs.Balances,
		ChannelID: cs.ID,
		Final:     cs.Final,
		Version:   cs.Version,
	}
}

func (cs ChannelState) Decode() types.ChannelState {
	return types.MakeChannelState(cs.ChannelID, cs.Balances, cs.Version, cs.Final)
}

type ChannelID [ChannelIDLength]byte

func (cid *ChannelID) UnmarshalJSON(bytes []byte) error {
	var hexString string
	err := json.Unmarshal(bytes, &hexString)
	if err != nil {
		return err
	}
	decode, err := hex.DecodeString(hexString)
	if err != nil {
		return err
	}
	if len(decode) != ChannelIDLength {
		return fmt.Errorf(
			"decoded channel id has wrong length. expected: %d, actual: %d",
			ChannelIDLength,
			len(decode),
		)
	}
	copy(cid[:], decode)
	return nil
}

func (cid ChannelID) MarshalJSON() ([]byte, error) {
	hexString := hex.EncodeToString(cid[:])
	return json.Marshal(hexString)
}
