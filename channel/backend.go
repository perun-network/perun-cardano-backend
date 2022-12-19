package channel

import (
	"fmt"
	pchannel "perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/blake2b224"
	"perun.network/perun-cardano-backend/channel/types"
	remotewallet "perun.network/perun-cardano-backend/wallet"
)

// backend implements the backend interface
// The type is private since it only needs to be exposed as singleton by the
// `Backend` variable.
// The current version of backend needs to use our wallet.RemoteBackend implementation.
// This is a workaround that makes encoding state for signing and verifying possible.
type backend struct {
	walletBackend *remotewallet.RemoteBackend
}

// SetWalletBackend needs to be called initially.
func SetWalletBackend(remoteBackend *remotewallet.RemoteBackend) {
	Backend = backend{walletBackend: remoteBackend}
}

// CalcID calculates the channel-id from the parameters.
func (b backend) CalcID(params *pchannel.Params) pchannel.ID {
	encodedParams, err := EncodeParams(params)
	if err != nil {
		panic(fmt.Sprintf("cannot calculate channel id: %v", err))
	}
	hash, err := blake2b224.Sum224(encodedParams)
	if err != nil {
		panic(fmt.Sprintf("unable to hash encoded parameters to compute channel-id: %v", err))
	}
	// We extend the hash with zero-padding to arrive at the 32 byte channel id used by go-perun.
	var id pchannel.ID
	copy(id[:blake2b224.Size224], hash[:])
	return id
}

// Sign signs the given state with the given account.
func (b backend) Sign(account wallet.Account, state *pchannel.State) (wallet.Sig, error) {
	remoteAccount, ok := account.(*remotewallet.RemoteAccount)
	if !ok {
		return nil, fmt.Errorf("unable to cast Account to RemoteAccount")
	}

	channelState, err := types.ConvertChannelState(*state)
	if err != nil {
		return nil, fmt.Errorf("unable to convert state for signing: %w", err)
	}

	return remoteAccount.SignChannelState(channelState)
}

// Verify returns true, iff the signature is correct for the given state and address.
func (b backend) Verify(addr wallet.Address, state *pchannel.State, sig wallet.Sig) (bool, error) {
	channelState, err := types.ConvertChannelState(*state)
	if err != nil {
		return false, fmt.Errorf("unable to encode state for verifying: %w", err)
	}
	return Backend.walletBackend.VerifyChannelStateSignature(channelState, sig, addr)
}

// NewAsset returns a variable of type Asset, which can be used for unmarshalling an asset from its binary
// representation.
func (b backend) NewAsset() pchannel.Asset {
	return Asset
}

// EncodeParams placeholder for parameter encoding.
func EncodeParams(params *pchannel.Params) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

var Backend backend

var _ pchannel.Backend = Backend
