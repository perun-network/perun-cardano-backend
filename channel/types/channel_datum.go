package types

import "time"

type ChannelDatum struct {
	ChannelParameters ChannelParameters
	ChannelState      ChannelState
	Time              time.Time
	FundingBalances   []Balance
	Funded            bool
	Disputed          bool
}
