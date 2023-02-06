package wire

import "time"

type ChannelDatum struct {
	ChannelParameters ChannelParameters `json:"channelParameters"`
	ChannelState      ChannelState      `json:"state"`
	Time              time.Time         `json:"time"`
	Funding           []uint64          `json:"funding"`
	Funded            bool              `json:"funded"`
	Disputed          bool              `json:"disputed"`
}
