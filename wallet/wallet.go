package wallet

import (
	"fmt"
	"perun.network/go-perun/wallet"
)

// RemoteWallet is a cardano signing wallet using a remote wallet server.
// Note: If we decide to stick with a remote signing wallet, we will have to harden the wallet server!
type RemoteWallet struct {
	walletServer Remote
}

func NewRemoteWallet(remote Remote) *RemoteWallet {
	return &RemoteWallet{walletServer: remote}
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
	pubKey, ok := address.(*PubKey)
	if !ok {
		return nil, fmt.Errorf("invalid address for signature verification (expected type PubKey)")
	}

	available, err := w.walletServer.CallKeyAvailable(*pubKey)
	if err != nil {
		return nil, fmt.Errorf("wallet server could not assert key availability: %w", err)
	}
	if !available {
		return nil, fmt.Errorf("wallet server has no private key for public key %s", pubKey.KeyString)
	}
	return RemoteAccount{*pubKey, w.walletServer}, nil
}
