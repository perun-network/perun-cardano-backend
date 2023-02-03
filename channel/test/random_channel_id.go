package test

import (
	"math/rand"
	"perun.network/go-perun/channel"
	"perun.network/perun-cardano-backend/channel/types"
)

func MakeRandomChannelID(rng *rand.Rand) types.ID {
	id := channel.ID{}
	rng.Read(id[:])
	return id
}
