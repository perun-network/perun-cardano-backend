package address

import (
	"encoding/hex"
	"fmt"
	"perun.network/go-perun/wallet"
)

// PubKeyLength is the length of the public verification key part of a Cardano `ed25519` keypair in bytes.
const PubKeyLength = 32

// Address carries a public key that represents the public verification key part of a Cardano `ed25519` keypair
type Address struct {
	pubKey [PubKeyLength]byte
}

// MakeAddressFromByteArray returns a new Address for the given public key bytes.
func MakeAddressFromByteArray(pubKey [PubKeyLength]byte) Address {
	return Address{
		pubKey: pubKey,
	}
}

// MakeAddressFromByteSlice returns a new Address for the given public key bytes.
func MakeAddressFromByteSlice(pubKey []byte) (Address, error) {
	addr := Address{}
	err := addr.UnmarshalBinary(pubKey)
	return addr, err
}

// GetPubKey returns the public key of this address
func (a Address) GetPubKey() [PubKeyLength]byte {
	return a.pubKey
}

// GetPubKeySlice returns the public key of this address. The returned slice is of length PubKeyLength.
func (a Address) GetPubKeySlice() []byte {
	return a.pubKey[:]
}

// MarshalBinary decodes public key of this address into its byte representation.
// The returned byte slice has length PubKeyLength.
func (a Address) MarshalBinary() ([]byte, error) {
	return a.pubKey[:], nil
}

// UnmarshalBinary expects a byte slice of length PubKeyLength and decodes it into the received Address.
func (a *Address) UnmarshalBinary(data []byte) error {
	if len(data) != PubKeyLength {
		return fmt.Errorf("public key has incorrect length. expected: %d bytes actual: %d bytes",
			PubKeyLength,
			len(data))
	}
	copy(a.pubKey[:], data)
	return nil
}

// TODO: This should probably return an address like `addr1vpu5vlrf4xkxv2qpwngf6cjhtw542ayty80v8dyr49rf5eg0yu80w`
// String returns the public key as hex string.
func (a Address) String() string {
	return hex.EncodeToString(a.pubKey[:])
}

// Equal returns true, iff the given address is of type Address and their public keys are equal.
func (a Address) Equal(other wallet.Address) bool {
	otherAddress, ok := other.(*Address)
	if !ok {
		return false
	}
	return a.pubKey == otherAddress.pubKey
}

var _ wallet.Address = (*Address)(nil)
