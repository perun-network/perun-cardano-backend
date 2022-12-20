package wallet

import (
	"fmt"
	"perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/wallet/address"
	"perun.network/perun-cardano-backend/wire"
)

// RemoteWallet is a cardano signing wallet using a remote wallet server.
// Note: If we decide to stick with a remote signing wallet, we will have to harden the wallet server!
type RemoteWallet struct {
	walletServer Remote
}

// NewRemoteWallet returns a pointer to a new RemoteWallet struct associated with the given Remote wallet server.
func NewRemoteWallet(remote Remote) *RemoteWallet {
	return &RemoteWallet{walletServer: remote}
}

// LockAll is unimplemented due to the remote nature of this wallet.
func (w *RemoteWallet) LockAll() {
}

// IncrementUsage is unimplemented due to the remote nature of this wallet.
func (w *RemoteWallet) IncrementUsage(address wallet.Address) {
}

// DecrementUsage is unimplemented due to the remote nature of this wallet.
func (w *RemoteWallet) DecrementUsage(address wallet.Address) {
}

// Unlock returns the account of the given address, iff the wallet server associated with this RemoteWallet
// has that account.
func (w *RemoteWallet) Unlock(addr wallet.Address) (wallet.Account, error) {
	rwAddress, ok := addr.(*address.Address)
	if !ok {
		return nil, fmt.Errorf("invalid address for signature verification (expected type Address)")
	}

	available, err := w.walletServer.CallKeyAvailable(wire.MakeKeyAvailabilityRequest(*rwAddress))
	if err != nil {
		return nil, fmt.Errorf("wallet server could not assert key availability: %w", err)
	}
	if !available {
		return nil, fmt.Errorf("wallet server has no private key for address %s", rwAddress)
	}
	return MakeRemoteAccount(*rwAddress, w.walletServer), nil
}
