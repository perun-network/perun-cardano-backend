package wallet

import (
	"encoding/hex"
	"fmt"
	"perun.network/go-perun/wallet"
)

// PubKeyLength is the length of the public verification key part of a Cardano `ed25519` keypair in bytes.
const PubKeyLength = 32

// TODO: the Address implementation of wallet.Address only represents the public verification key.
// 	Cardano addresses also carry staking information though and look something like this:
//  `addr1vpu5vlrf4xkxv2qpwngf6cjhtw542ayty80v8dyr49rf5eg0yu80w`.
//  We should move to the latter address representation at some point.

// Address carries a public key that represents the public verification key part of a Cardano `ed25519` keypair
// (Ledger.Crypto.Address in the plutus-ledger library). The PubKey is a hex string (without "0x") representing the
// bytes of the public key.
type Address struct {
	PubKey string `json:"getPubKey"`
}

// MarshalBinary decodes public key of this address into its byte representation.
// The returned byte slice has length PubKeyLength.
func (a Address) MarshalBinary() ([]byte, error) {
	key, err := hex.DecodeString(a.PubKey)
	if err != nil {
		return nil, fmt.Errorf("pubkey string is not a hex string: %w", err)
	}
	if len(key) != PubKeyLength {
		return nil, fmt.Errorf(
			"public key has incorrect length. expected: %d bytes actual: %d bytes",
			PubKeyLength,
			len(key),
		)
	}
	return key, nil
}

// UnmarshalBinary expects a byte slice of length PubKeyLength and decodes it into the received Address.
func (a *Address) UnmarshalBinary(data []byte) error {
	if len(data) != PubKeyLength {
		return fmt.Errorf("public key has incorrect length. expected: %d bytes actual: %d bytes",
			PubKeyLength,
			len(data))
	}
	a.PubKey = hex.EncodeToString(data)
	return nil
}

// TODO: This should probably return an address like `addr1vpu5vlrf4xkxv2qpwngf6cjhtw542ayty80v8dyr49rf5eg0yu80w`
// String returns the public key string.
func (a Address) String() string {
	return a.PubKey
}

// Equal returns true, iff the given address is of type Address and their public keys are equal.
func (a Address) Equal(other wallet.Address) bool {
	otherAddress, ok := other.(*Address)
	if !ok {
		return false
	}
	return a.PubKey == otherAddress.PubKey
}

var _ wallet.Address = (*Address)(nil)
