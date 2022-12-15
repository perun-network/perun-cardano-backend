package wire

type ChannelState struct {
	Balances  []int `json:"balances"`
	ChannelID int   `json:"channelId"`
	Final     bool  `json:"final"`
	Version   int   `json:"version"`
}
