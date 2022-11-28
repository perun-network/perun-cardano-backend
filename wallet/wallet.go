package wallet

import (
	"perun.network/go-perun/wallet"
)

// RemoteWallet is a cardano signing wallet using a remote wallet server.
// Note: If we decide to stick with a remote signing wallet, we will have to harden the wallet server!
type RemoteWallet struct {
	walletServer Remote
}

// LockAll is unimplemented due to the remote nature of this wallet
func (w *RemoteWallet) LockAll() {
}

// IncrementUsage is unimplemented due to the remote nature of this wallet
func (w *RemoteWallet) IncrementUsage(address wallet.Address) {
}

// DecrementUsage is unimplemented due to the remote nature of this wallet
func (w *RemoteWallet) DecrementUsage(address wallet.Address) {
}

// Unlock verifies that the remote wallet server has the private key for this address before returning an Account
func (w *RemoteWallet) Unlock(address wallet.Address) (wallet.Account, error) {
	// Todo: verify that the wallet server has the private key associated with this address before returning an Account
	addr := *address.(*PubKey)
	return RemoteAccount{addr, w.walletServer}, nil
}
