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

package address

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/bech32"
	"perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/blake2b224"
)

// PubKeyLength is the length of the public verification key part of a Cardano `ed25519` keypair in bytes.
const PubKeyLength = 32

// PubKeyHashLength is the length of Cardano public key hashes using blake2b-224.
const PubKeyHashLength = blake2b224.Size224

const AddressLength = PubKeyLength + PubKeyHashLength

const MainnetIdentifier = "addr"
const TestnetIdentifier = "addr_test"
const mainnetIdentifierByte byte = 0x61
const testnetIdentifierByte byte = 0x60

// Address carries a public key that represents the public verification key part of a Cardano `ed25519` keypair.
type Address struct {
	// We currently do not implement BIP-44 address derivation. As a compromise, we allow to keep both a public key
	// (the public key under which signatures e.g. on channel state are verified) and a public key hash (under which
	// this address is supposed to receive payments).
	pubKey            [PubKeyLength]byte
	paymentPubKeyHash [PubKeyHashLength]byte
}

// MakeAddressFromSinglePubKey creates an Address from a single public key. This means the payment and signing public
// keys are the same. The (payment) public key hash is the blake2b-224 hash of the given public key.
func MakeAddressFromSinglePubKey(pubKey []byte) (Address, error) {
	a, err := MakeAddressFromPubKeyByteSlice(pubKey)
	if err != nil {
		return a, err
	}
	pubKeyHash, err := CalculatePubKeyHash(a.pubKey)
	if err != nil {
		return Address{}, err
	}
	a.paymentPubKeyHash = pubKeyHash
	return a, nil
}

func (a *Address) SetPaymentPubKeyHash(paymentPubKeyHash [PubKeyHashLength]byte) {
	copy(a.paymentPubKeyHash[:], paymentPubKeyHash[:])
}

func (a *Address) SetPaymentPubKeyHashFromSlice(paymentPubKeyHash []byte) error {
	if len(paymentPubKeyHash) != PubKeyHashLength {
		return fmt.Errorf(
			"payment public key hash has incorrect length. expected: %d bytes, got: %d bytes",
			PubKeyHashLength,
			len(paymentPubKeyHash),
		)
	}
	copy(a.paymentPubKeyHash[:], paymentPubKeyHash)
	return nil
}

func (a *Address) SetPaymentPubKeyHashFromHexString(paymentPubKeyHash string) error {
	paymentPubKeyHashBytes, err := hex.DecodeString(paymentPubKeyHash)
	if err != nil {
		return err
	}
	return a.SetPaymentPubKeyHashFromSlice(paymentPubKeyHashBytes)
}

// MakeAddressFromPubKeyByteArray returns a new Address for the given public key bytes.
// Note: This does not set the public key hash!
func MakeAddressFromPubKeyByteArray(pubKey [PubKeyLength]byte) Address {
	return Address{
		pubKey: pubKey,
	}
}

// MakeAddressFromPubKeyByteSlice returns a new Address for the given public key bytes.
// Note: This does not set the public key hash!
func MakeAddressFromPubKeyByteSlice(pubKey []byte) (Address, error) {
	if len(pubKey) != PubKeyLength {
		return Address{}, fmt.Errorf("public key has incorrect length. expected: %d bytes actual: %d bytes", PubKeyLength, len(pubKey))
	}
	addr := Address{}
	copy(addr.pubKey[:], pubKey)
	return addr, nil
}

// GetPubKey returns the public key of this address.
func (a Address) GetPubKey() [PubKeyLength]byte {
	return a.pubKey
}

// GetPubKeySlice returns the public key of this address. The returned slice is of length PubKeyLength.
func (a Address) GetPubKeySlice() []byte {
	return a.pubKey[:]
}

// MarshalBinary decodes this address into its byte representation.
// The returned byte slice has length PubKeyLength + PubKeyHashLength.
func (a Address) MarshalBinary() ([]byte, error) {
	return append(a.pubKey[:], a.paymentPubKeyHash[:]...), nil
}

// UnmarshalBinary expects a byte slice of length PubKeyLength + PubKeyHashLength and decodes it into the
// receiver Address.
func (a *Address) UnmarshalBinary(data []byte) error {
	if len(data) != AddressLength {
		return fmt.Errorf("public key has incorrect length. expected: %d bytes actual: %d bytes",
			AddressLength,
			len(data))
	}
	copy(a.pubKey[:], data[:PubKeyLength])
	copy(a.paymentPubKeyHash[:], data[PubKeyLength:])
	return nil
}

// String returns the public key as hex string.
func (a Address) String() string {
	return hex.EncodeToString(a.pubKey[:])
}

// GetTestnetAddressOfPubKey returns the testnet address string representation of this addresses' public key
// (i.e. `addr_test1...`).
func (a Address) GetTestnetAddressOfPubKey() (string, error) {
	hash, err := CalculatePubKeyHash(a.pubKey)
	if err != nil {
		return "", fmt.Errorf("unable to compute blake2b hash of public key: %w", err)
	}
	return convertToAddress(hash, testnetIdentifierByte, TestnetIdentifier)
}

// GetMainnetAddressOfPubKey returns the mainnet address string representation of this addresses' public key
// (i.e. `addr1...`).
func (a Address) GetMainnetAddressOfPubKey() (string, error) {
	hash, err := CalculatePubKeyHash(a.pubKey)
	if err != nil {
		return "", fmt.Errorf("unable to compute blake2b hash of public key: %w", err)
	}
	return convertToAddress(hash, mainnetIdentifierByte, MainnetIdentifier)
}

// GetTestnetAddressOfPubKeyHash returns the testnet address string representation of this addresses' payment public key
// hash (i.e. `addr1...`).
func (a Address) GetTestnetAddressOfPubKeyHash() (string, error) {
	return convertToAddress(a.paymentPubKeyHash, testnetIdentifierByte, TestnetIdentifier)
}

// GetMainnetAddressOfPubKeyHash returns the mainnet address string representation of this addresses' payment public key
// hash (i.e. `addr1...`).
func (a Address) GetMainnetAddressOfPubKeyHash() (string, error) {
	return convertToAddress(a.paymentPubKeyHash, mainnetIdentifierByte, MainnetIdentifier)
}

// convertToAddress returns the address string for given network parameters.
func convertToAddress(hash [PubKeyHashLength]byte, networkIdentifierByte byte, networkIdentifierString string) (string, error) {

	// Bech32-encode the PubKeyHash.
	conv, err := bech32.ConvertBits(append([]byte{networkIdentifierByte}, hash[:]...), 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("unable to convert bits for bech32 encoding: %w", err)
	}
	encodedHash, err := bech32.Encode(networkIdentifierString, conv)
	if err != nil {
		return "", fmt.Errorf("unable bech32-encode: %w", err)
	}
	return encodedHash, nil
}

// GetPubKeyHash returns the public key hash associated with payments to this address.
func (a Address) GetPubKeyHash() [PubKeyHashLength]byte {
	return a.paymentPubKeyHash
}

func (a Address) GetPubKeyHashSlice() []byte {
	return a.paymentPubKeyHash[:]
}

func CalculatePubKeyHash(pubKey [PubKeyLength]byte) ([PubKeyHashLength]byte, error) {
	return blake2b224.Sum224(pubKey[:])
}

// Equal returns true, iff the given address is of type Address and their public keys are equal.
func (a Address) Equal(other wallet.Address) bool {
	otherAddress, ok := other.(*Address)
	if !ok {
		return false
	}
	return a.pubKey == otherAddress.pubKey && a.paymentPubKeyHash == otherAddress.paymentPubKeyHash
}

var _ wallet.Address = (*Address)(nil)
