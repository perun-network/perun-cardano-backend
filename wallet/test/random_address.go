package test

import (
	"math/rand"
	"perun.network/perun-cardano-backend/wallet/address"
)

func MakeRandomAddress(rng *rand.Rand) address.Address {
	addrBytes := [address.PubKeyLength]byte{}
	rng.Read(addrBytes[:])
	return address.MakeAddressFromByteArray(addrBytes)
}

func MakeTooFewPublicKeyBytes(rng *rand.Rand) []byte {
	tooFewBytes := make([]byte, rng.Intn(address.PubKeyLength))
	rng.Read(tooFewBytes)
	return tooFewBytes
}

func MakeTooManyPublicKeyBytes(rng *rand.Rand) []byte {
	const maxInvalidPubKeyLength = address.PubKeyLength * 2
	tooManyBytes := make([]byte, rng.Intn(maxInvalidPubKeyLength-address.PubKeyLength)+address.PubKeyLength+1)
	rng.Read(tooManyBytes)
	return tooManyBytes
}
