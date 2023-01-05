package types

import (
	"fmt"
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
		Nonce:   params.Nonce,
		Timeout: time.Duration(params.ChallengeDuration) * time.Second,
	}, nil
}
