package wire

import (
	"perun.network/perun-cardano-backend/channel/types"
	"time"
)

type ChannelDatum struct {
	ChannelParameters ChannelParameters `json:"channelParameters"`
	ChannelState      ChannelState      `json:"state"`
	Time              time.Time         `json:"time"`
	Funding           []uint64          `json:"funding"`
	Funded            bool              `json:"funded"`
	Disputed          bool              `json:"disputed"`
}

func (c ChannelDatum) Decode() (types.ChannelDatum, error) {
	p, err := c.ChannelParameters.Decode()
	if err != nil {
		return types.ChannelDatum{}, err
	}
	return types.ChannelDatum{
		ChannelParameters: p,
		ChannelState:      c.ChannelState.Decode(),
		Time:              c.Time,
		FundingBalances:   c.Funding,
		Funded:            c.Funded,
		Disputed:          c.Disputed,
	}, nil
}
