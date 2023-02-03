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

func MakeChannelParameters(params channel.Params) (ChannelParameters, error) {
	// TODO assert that the params are valid for the current state of the Cardano backend (e.g. no virtual channels)
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
