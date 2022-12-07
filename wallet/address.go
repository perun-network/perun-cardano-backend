package wallet

import (
	"encoding/hex"
	"fmt"
	"perun.network/go-perun/wallet"
)

// PubKeyLength is the length of the public verification key part of a Cardano `ed25519` keypair in bytes.
const PubKeyLength = 32

// TODO: the PubKey implementation of wallet.Address only represents the public verification key.
// 	Cardano addresses also carry staking information though and look something like this:
//  `addr1vpu5vlrf4xkxv2qpwngf6cjhtw542ayty80v8dyr49rf5eg0yu80w`.
//  We should move to the latter address representation at some point.

// PubKey represents the public verification key part of a Cardano `ed25519` keypair
// (Ledger.Crypto.PubKey in the plutus-ledger library). The Key is a hex string (without "0x") representing the bytes of
// the public key.
type PubKey struct {
	Key string `json:"getPubKey"`
}

// MarshalBinary decodes the hexadecimal Key into byte representation. The returned byte slice has length PubKeyLength.
func (a PubKey) MarshalBinary() ([]byte, error) {
	key, err := hex.DecodeString(a.Key)
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

// UnmarshalBinary expects a byte slice of length PubKeyLength and decodes it into the PubKey receiver.
func (a *PubKey) UnmarshalBinary(data []byte) error {
	if len(data) != PubKeyLength {
		return fmt.Errorf("public key has incorrect length. expected: %d bytes actual: %d bytes",
			PubKeyLength,
			len(data))
	}
	a.Key = hex.EncodeToString(data)
	return nil
}

// String returns the key string
func (a PubKey) String() string {
	return a.Key
}

// Equal returns true, iff the given address is of type PubKey and their Key values are equal.
func (a PubKey) Equal(other wallet.Address) bool {
	otherPubKey, ok := other.(*PubKey)
	if !ok {
		return false
	}
	return a.Key == otherPubKey.Key
}

var _ wallet.Address = (*PubKey)(nil)
