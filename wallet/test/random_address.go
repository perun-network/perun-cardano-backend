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
