package wire

type Event struct {
	Tag          string         `json:"tag"`
	ChannelDatum []ChannelDatum `json:"channelDatum"`
}
