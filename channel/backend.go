package channel

import (
	"fmt"
	pchannel "perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/blake2b224"
)

// backend implements the backend interface
// The type is private since it only needs to be exposed as singleton by the
// `Backend` variable.
type backend struct{}

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
