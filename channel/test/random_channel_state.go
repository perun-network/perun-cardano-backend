package test

import (
	"math/rand"
	gpchannel "perun.network/go-perun/channel"
	"perun.network/perun-cardano-backend/channel/types"
)

func MakeRandomChannelState(rng *rand.Rand) types.ChannelState {
	var channelID = gpchannel.ID{}
	rng.Read(channelID[:])
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
