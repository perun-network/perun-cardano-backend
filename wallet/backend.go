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
	"errors"
	"fmt"
	"io"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/channel/types"
	"perun.network/perun-cardano-backend/wallet/address"
	"perun.network/perun-cardano-backend/wire"
)

// RemoteBackend is a wallet.Backend implementation with a remote server for signing data and verifying signatures.
type RemoteBackend struct {
	walletServer Remote
}

// MakeRemoteBackend returns a new RemoteBackend struct.
func MakeRemoteBackend(remote Remote) RemoteBackend {
	return RemoteBackend{remote}
}

// NewAddress returns a pointer to a new, empty address.
func (b RemoteBackend) NewAddress() wallet.Address {
	return new(address.Address)
}

// DecodeSig reads SignatureLength bytes from the given reader and returns the read signature.
func (b RemoteBackend) DecodeSig(reader io.Reader) (wallet.Sig, error) {
	sig := make([]byte, wire.SignatureLength)
	if _, err := io.ReadFull(reader, sig); err != nil {
		return nil, fmt.Errorf("unable to read signature from reader: %w", err)
	}
	return sig, nil
}

// VerifySignature returns true, iff the given signature is valid for the given message under the public key associated
// with the given address.
func (b RemoteBackend) VerifySignature(msg []byte, sig wallet.Sig, a wallet.Address) (bool, error) {
	addr, ok := a.(*address.Address)
	if !ok {
		return false, fmt.Errorf("invalid PubKey for signature verification")
	}
	if len(sig) != wire.SignatureLength {
		return false, fmt.Errorf(
			"signature has incorrect length. expected: %d bytes actual: %d bytes",
			wire.SignatureLength,
			len(sig),
		)
	}
	request := wire.MakeVerificationRequest(sig, *addr, msg)
	var response wire.VerificationResponse
	err := b.walletServer.CallEndpoint(EndpointVerifyDataSignature, request, &response)
	if err != nil {
		return false, fmt.Errorf("wallet server could not verify message: %w", err)
	}
	return response, nil
}

// VerifyChannelStateSignature returns true, iff the given signature is valid for the given ChannelState under the
// public key associated with the given address.
func (b RemoteBackend) VerifyChannelStateSignature(state types.ChannelState, sig wallet.Sig, a wallet.Address) (bool, error) {
	addr, ok := a.(*address.Address)
	if !ok {
		return false, fmt.Errorf("invalid PubKey for signature verification")
	}
	if len(sig) != wire.SignatureLength {
		return false, fmt.Errorf(
			"signature has incorrect length. expected: %d bytes actual: %d bytes",
			wire.SignatureLength,
			len(sig),
		)
	}
	request := wire.MakeChannelStateVerificationRequest(sig, *addr, state)
	var response wire.VerificationResponse
	err := b.walletServer.CallEndpoint(EndpointVerifyChannelStateSignature, request, &response)
	if err != nil {
		return false, fmt.Errorf("wallet server could not verify message: %w", err)
	}
	return response, nil
}

// CalculateChannelID returns the channelId for the given parameters as calculated by the remote instance.
func (b RemoteBackend) CalculateChannelID(parameters types.ChannelParameters) (channel.ID, error) {
	request := wire.MakeChannelParameters(parameters)
	var response wire.ChannelID
	err := b.walletServer.CallEndpoint(EndpointCalculateChannelID, request, &response)
	if err != nil {
		return response, fmt.Errorf("wallet server was unable to compute ChannelID: %w", err)
	}
	return response, nil
}

func (b RemoteBackend) ToChannelStateSigningAccount(account wallet.Account) (types.ChannelStateSigningAccount, error) {
	acc, ok := account.(RemoteAccount)
	if !ok {
		return acc, errors.New("account is not a RemoteAccount")
	}
	return acc, nil
}

var _ wallet.Backend = RemoteBackend{}
