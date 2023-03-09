package wire

import (
	"perun.network/perun-cardano-backend/channel/types"
	"time"
)

type ChannelDatum struct {
	ChannelParameters ChannelParameters `json:"channelParameters"`
	ChannelToken      ChannelToken      `json:"channelToken"`
	Disputed          bool              `json:"disputed"`
	Funded            bool              `json:"funded"`
	Funding           []uint64          `json:"funding"`
	ChannelState      ChannelState      `json:"state"`
	Time              int64             `json:"time"`
}

func (c ChannelDatum) Decode() (types.ChannelDatum, error) {
	p, err := c.ChannelParameters.Decode()
	if err != nil {
		return types.ChannelDatum{}, err
	}
	return types.ChannelDatum{
		ChannelParameters: p,
		ChannelToken:      c.ChannelToken.Decode(),
		ChannelState:      c.ChannelState.Decode(),
		Time:              time.UnixMilli(c.Time),
		FundingBalances:   c.Funding,
		Funded:            c.Funded,
		Disputed:          c.Disputed,
	}, nil
}

func MakeChannelDatum(datum types.ChannelDatum) ChannelDatum {
	return ChannelDatum{
		ChannelParameters: MakeChannelParameters(datum.ChannelParameters),
		ChannelToken:      MakeChannelToken(datum.ChannelToken),
		ChannelState:      MakeChannelState(datum.ChannelState),
		Time:              datum.Time.UnixMilli(),
		Funding:           datum.FundingBalances,
		Funded:            datum.Funded,
		Disputed:          datum.Disputed,
	}
}
