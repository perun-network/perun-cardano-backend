// Copyright 2022, 2023 - See NOTICE file for copyright holders.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	var response wire.KeyAvailabilityResponse
	err := w.walletServer.CallEndpoint(EndpointKeyAvailable, wire.MakeKeyAvailabilityRequest(*rwAddress), &response)
	if err != nil {
		return nil, fmt.Errorf("wallet server could not assert key availability: %w", err)
	}
	if !response {
		return nil, fmt.Errorf("wallet server has no private key for address %s", rwAddress)
	}
	return MakeRemoteAccount(*rwAddress, w.walletServer), nil
}
