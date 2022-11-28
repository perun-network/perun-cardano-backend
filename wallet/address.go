package wallet

import (
	"encoding/hex"
	"fmt"
	"perun.network/go-perun/wallet"
)

// PubKeyLength is the length of a public key in bytes (not the length of hex string!)
const PubKeyLength = 32

// PubKey represents a Cardano public key struct (Ledger.Crypto.PubKey in the plutus-ledger library).
// The KeyString is a hex string (without "0x") representing the bytes of the public key
type PubKey struct {
	KeyString string `json:"getPubKey"`
}

// MarshalBinary decodes the hexadecimal key string and returns the represented bytes after verifying the PubKeyLength
func (a PubKey) MarshalBinary() ([]byte, error) {
	key, err := hex.DecodeString(a.KeyString)
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

// UnmarshalBinary verifies the PubKeyLength and then sets KeyString to the hexadecimal encoding of the bytes
func (a *PubKey) UnmarshalBinary(data []byte) error {
	if len(data) != PubKeyLength {
		return fmt.Errorf("public key has incorrect length. expected: %d bytes actual: %d bytes",
			PubKeyLength,
			len(data))
	}
	a.KeyString = hex.EncodeToString(data)
	return nil
}

// String returns the KeyString
func (a PubKey) String() string {
	return a.KeyString
}

// Equal tests two public keys for equality by comparing their KeyString values
func (a PubKey) Equal(b wallet.Address) bool {
	b_, ok := b.(*PubKey)
	if !ok {
		return false
	}
	return a.KeyString == b_.KeyString
}

var _ wallet.Address = (*PubKey)(nil)
