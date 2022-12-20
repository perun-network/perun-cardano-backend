package test

import (
	"math/rand"
	"perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/wallet/address"
	"perun.network/perun-cardano-backend/wire"
)

// GetRandomByteSlice returns a byte slice of a random length in [minLength, maxLength] filled with random values.
func GetRandomByteSlice(minLength int, maxLength int, rng *rand.Rand) []byte {
	randomBytes := make([]byte, rng.Intn(maxLength-minLength+1)+minLength)
	rng.Read(randomBytes)
	return randomBytes
}

func MakeRandomAddress(rng *rand.Rand) address.Address {
	addrBytes := [address.PubKeyLength]byte{}
	rng.Read(addrBytes[:])
	return address.MakeAddressFromByteArray(addrBytes)
}

func MakeTooFewPublicKeyBytes(rng *rand.Rand) []byte {
	return GetRandomByteSlice(0, address.PubKeyLength-1, rng)
}

func MakeTooManyPublicKeyBytes(rng *rand.Rand) []byte {
	const maxInvalidPubKeyLength = address.PubKeyLength * 2
	return GetRandomByteSlice(address.PubKeyLength+1, maxInvalidPubKeyLength, rng)
}

func MakeTooLongSignature(rng *rand.Rand) wallet.Sig {
	const maxInvalidSigLength = 0x100
	return GetRandomByteSlice(wire.SignatureLength+1, maxInvalidSigLength, rng)
}

func MakeTooShortSignature(rng *rand.Rand) wallet.Sig {
	return GetRandomByteSlice(0, wire.SignatureLength-1, rng)
}
