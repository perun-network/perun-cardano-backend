package types

type ChannelToken struct {
	TokenSymbol string
	TokenName   string
	TxOutRef    TxOutRef
}

type TxOutRef struct {
	TxID  string
	Index int
}
