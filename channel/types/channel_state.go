package types

import (
	"perun.network/go-perun/channel"
	"reflect"
)

// ChannelState is the cardano backend equivalent to go-perun's channel.State.
type ChannelState struct {
	ID       channel.ID
	Balances []uint64
	Version  uint64
	Final    bool
}

func MakeChannelState(id channel.ID, balances []uint64, version uint64, final bool) ChannelState {
	return ChannelState{
		ID:       id,
		Balances: balances,
		Version:  version,
		Final:    final,
	}
}

func (cs ChannelState) Equal(other ChannelState) bool {
	return cs.ID == other.ID &&
		reflect.DeepEqual(cs.Balances, other.Balances) &&
		cs.Version == other.Version &&
		cs.Final == other.Final
}
