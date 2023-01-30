package types

import (
	"perun.network/go-perun/channel"
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
