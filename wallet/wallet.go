package wallet

import (
	"perun.network/go-perun/wallet"
)

type RemoteWallet struct {
}

func (w *RemoteWallet) NewAccount() Account {
	//contact server to generate new account
	return *new(Account)
}

func (w *RemoteWallet) Unlock(a PaymentPubKeyHash) (wallet.Account, error) {
	// TODO: Verify with server that Address exists
	return Account{a}, nil
}
