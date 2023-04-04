package test

import (
	"math/rand"
	gpwallet "perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/wallet"
	"perun.network/perun-cardano-backend/wallet/address"
	"perun.network/perun-cardano-backend/wire"
)

func MakeRemoteAccount(address address.Address, remote wallet.Remote) wallet.RemoteAccount {
	return wallet.MakeRemoteAccount(address, remote, "")
}

func NewRemoteWallet(remote wallet.Remote) *wallet.RemoteWallet {
	return wallet.NewRemoteWallet(remote, "")
}

// GetRandomByteSlice returns a byte slice of a random length in [minLength, maxLength] filled with random values.
func GetRandomByteSlice(minLength int, maxLength int, rng *rand.Rand) []byte {
	randomBytes := make([]byte, rng.Intn(maxLength-minLength+1)+minLength)
	rng.Read(randomBytes)
	return randomBytes
}

func MakeRandomAddress(rng *rand.Rand) address.Address {
	addrBytes := [address.PubKeyLength]byte{}
	rng.Read(addrBytes[:])
	addr := address.MakeAddressFromPubKeyByteArray(addrBytes)
	pkhBytes := [address.PubKeyHashLength]byte{}
	rng.Read(pkhBytes[:])
	_ = addr.SetPaymentPubKeyHashFromSlice(pkhBytes[:])
	return addr
}

func MakeRandomSignature(rng *rand.Rand) gpwallet.Sig {
	sig := make([]byte, wire.SignatureLength)
	rng.Read(sig)
	return sig
}

func MakeTooFewPubKeyBytes(rng *rand.Rand) []byte {
	return GetRandomByteSlice(0, address.PubKeyLength-1, rng)
}

func MakeTooManyPubKeyBytes(rng *rand.Rand) []byte {
	const maxInvalidPubKeyLength = address.PubKeyLength * 2
	return GetRandomByteSlice(address.PubKeyLength+1, maxInvalidPubKeyLength, rng)
}

func MakeTooFewAddressBytes(rng *rand.Rand) []byte {
	return GetRandomByteSlice(0, address.AddressLength-1, rng)
}

func MakeTooManyAddressBytes(rng *rand.Rand) []byte {
	const maxInvalidAddressLength = address.AddressLength * 2
	return GetRandomByteSlice(address.AddressLength+1, maxInvalidAddressLength, rng)
}

func MakeTooLongSignature(rng *rand.Rand) gpwallet.Sig {
	const maxInvalidSigLength = 0x100
	return GetRandomByteSlice(wire.SignatureLength+1, maxInvalidSigLength, rng)
}

func MakeTooShortSignature(rng *rand.Rand) gpwallet.Sig {
	return GetRandomByteSlice(0, wire.SignatureLength-1, rng)
}
