// Copyright 2023 - See NOTICE file for copyright holders.
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

package channel

import (
	"fmt"
	pchannel "perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/channel/types"
)

// backend implements the wallet.Backend interface
// The type is private since it only needs to be exposed as a singleton by the
// `Backend` variable.
// The current version of backend needs to use our wallet.RemoteBackend implementation.
// This is a workaround that makes encoding state for signing and verifying possible.
type backend struct {
	walletBackend types.ExtendedWalletBackend
}

// SetWalletBackend needs to be called initially.
func SetWalletBackend(remoteBackend types.ExtendedWalletBackend) {
	Backend = backend{walletBackend: remoteBackend}
}

// CalcID calculates the channel-id from the parameters.
// Note that the remote wallet is used for this.
func (b backend) CalcID(params *pchannel.Params) pchannel.ID {
	if params == nil {
		panic("params must not be nil for channel id calculation")
	}
	p, err := types.MakeChannelParameters(*params)
	if err != nil {
		panic(err)
	}
	id, err := b.walletBackend.CalculateChannelID(p)
	if err != nil {
		panic(err)
	}
	return id
}

// Sign signs the given state with the given account.
func (b backend) Sign(account wallet.Account, state *pchannel.State) (wallet.Sig, error) {
	if account == nil {
		return nil, fmt.Errorf("account must not be nil for signing")
	}
	if state == nil {
		return nil, fmt.Errorf("state must not be nil for signing")
	}
	acc, err := b.walletBackend.ToChannelStateSigningAccount(account)
	if err != nil {
		return nil, err
	}

	channelState, err := types.ConvertChannelState(*state)
	if err != nil {
		return nil, fmt.Errorf("unable to convert state for signing: %w", err)
	}

	return acc.SignChannelState(channelState)
}

// Verify returns true, iff the signature is correct for the given state and address.
func (b backend) Verify(addr wallet.Address, state *pchannel.State, sig wallet.Sig) (bool, error) {
	if addr == nil {
		return false, fmt.Errorf("address must not be nil for verification")
	}
	if state == nil {
		return false, fmt.Errorf("state must not be nil for verification")
	}
	if sig == nil {
		return false, fmt.Errorf("signature must not be nil for verification")
	}
	channelState, err := types.ConvertChannelState(*state)
	if err != nil {
		return false, fmt.Errorf("unable to encode state for verifying: %w", err)
	}
	return Backend.walletBackend.VerifyChannelStateSignature(channelState, sig, addr)
}

// NewAsset returns a variable of type Asset, which can be used for unmarshalling an asset from its binary
// representation.
func (b backend) NewAsset() pchannel.Asset {
	return types.Asset
}

var Backend backend

var _ pchannel.Backend = Backend
