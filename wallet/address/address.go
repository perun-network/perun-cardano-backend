package address

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/bech32"
	"golang.org/x/crypto/blake2b"
	"perun.network/go-perun/wallet"
)

// PubKeyLength is the length of the public verification key part of a Cardano `ed25519` keypair in bytes.
const PubKeyLength = 32
const PubKeyHashLength = 28

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

// GetPubKey returns the public key of this address.
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

// String returns the public key as hex string.
func (a Address) String() string {
	return hex.EncodeToString(a.pubKey[:])
}

// GetTestnetAddress returns the testnet address string representation of this Address (i.e. `addr_test1...`).
func (a Address) GetTestnetAddress() (string, error) {
	const testnetIdentifierByte byte = 0x60
	const testnetIdentifierString = "addr_test"
	return a.convertToAddress(testnetIdentifierByte, testnetIdentifierString)
}

// GetMainnetAddress returns the mainnet address string representation of this Address (i.e. `addr1...`).
func (a Address) GetMainnetAddress() (string, error) {
	const mainnetIdentifierByte byte = 0x61
	const mainnetIdentifierString = "addr"
	return a.convertToAddress(mainnetIdentifierByte, mainnetIdentifierString)
}

// convertToAddress returns the address string for given network parameters.
func (a Address) convertToAddress(networkIdentifierByte byte, networkIdentifierString string) (string, error) {
	pubKeyHash, err := a.GetPubKeyHash()
	if err != nil {
		return "", fmt.Errorf("unable to compute PubKeyHash: %w", err)
	}

	// Bech32-encode the PubKeyHash.
	conv, err := bech32.ConvertBits(append([]byte{networkIdentifierByte}, pubKeyHash[:]...), 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("unable to convert bits for bech32 encoding: %w", err)
	}
	encodedHash, err := bech32.Encode(networkIdentifierString, conv)
	if err != nil {
		return "", fmt.Errorf("unable bech32-encode: %w", err)
	}
	return encodedHash, nil
}

// GetPubKeyHash returns the blake2b224-hash of the public key associated with this address.
func (a Address) GetPubKeyHash() ([PubKeyHashLength]byte, error) {
	// Get a blake2b-224 hash function instance.
	blake2b224, err := blake2b.New(28, nil)
	if err != nil {
		return [PubKeyHashLength]byte{}, fmt.Errorf("unable to create blake2b-224 instance: %w", err)
	}
	// Compute the hash of the public key
	n, err := blake2b224.Write(a.GetPubKeySlice())
	if err != nil {
		return [PubKeyHashLength]byte{}, fmt.Errorf("unable to compute blake2b224 hash of public key: %w", err)
	}
	if n != PubKeyLength {
		return [PubKeyHashLength]byte{}, fmt.Errorf("did not write enough bytes. expected: %d, wrote: %d", PubKeyLength, n)
	}
	res := blake2b224.Sum(nil)
	if len(res) != PubKeyHashLength {
		return [PubKeyHashLength]byte{}, fmt.Errorf(
			"hash does not have expected length. expected: %d, actual: %d",
			PubKeyHashLength,
			len(res),
		)
	}
	ret := [PubKeyHashLength]byte{}

	copy(ret[:], res)
	// Read and return the hash.
	return ret, nil
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
