package channel

import (
	"fmt"
	"golang.org/x/crypto/sha3"
	pchannel "perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
)

// backend implements the backend interface
// The type is private since it only needs to be exposed as singleton by the
// `Backend` variable.
type backend struct{}

// CalcID calculates the channel-id from the parameters
func (b backend) CalcID(params *pchannel.Params) pchannel.ID {
	encodedParams, err := EncodeParams(params)
	if err != nil {
		panic(fmt.Sprintf("cannot calculate channel id: %v", err))
	}

	return sha3.Sum256(encodedParams)
}

// Sign signs the given state with the given account.
func (b backend) Sign(account wallet.Account, state *pchannel.State) (wallet.Sig, error) {
	encodedState, err := EncodeState(state)
	if err != nil {
		return nil, fmt.Errorf("unable to encode state for signing: %w", err)
	}
	return account.SignData(encodedState)
}

// Verify returns true, iff the signature is correct for the given state and address.
func (b backend) Verify(addr wallet.Address, state *pchannel.State, sig wallet.Sig) (bool, error) {
	encodedState, err := EncodeState(state)
	if err != nil {
		return false, fmt.Errorf("unable to encode state for verifying: %w", err)
	}
	return wallet.VerifySignature(encodedState, sig, addr)
}

// NewAsset returns a variable of type Asset, which can be used for unmarshalling an asset from its binary
// representation.
func (b backend) NewAsset() pchannel.Asset {
	return Asset
}

// EncodeState is a placeholder for state encoding.
func EncodeState(state *pchannel.State) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

// EncodeParams placeholder for parameter encoding.
func EncodeParams(params *pchannel.Params) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

var Backend backend

var _ pchannel.Backend = Backend
