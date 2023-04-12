package types

import (
	pchannel "perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

type ExtendedWalletBackend interface {
	wallet.Backend
	VerifyChannelStateSignature(state ChannelState, sig wallet.Sig, a wallet.Address) (bool, error)
	CalculateChannelID(parameters ChannelParameters) (pchannel.ID, error)
	ToChannelStateSigningAccount(account wallet.Account) (ChannelStateSigningAccount, error)
}

type ChannelStateSigningAccount interface {
	wallet.Account
	SignChannelState(state ChannelState) (wallet.Sig, error)
}
