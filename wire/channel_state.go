package wire

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"perun.network/perun-cardano-backend/channel/types"
)

const ChannelIDLength = 32

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
