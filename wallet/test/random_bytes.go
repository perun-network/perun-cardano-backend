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
	addr := address.MakeAddressFromPubKeyByteArray(addrBytes)
	pkhBytes := [address.PubKeyHashLength]byte{}
	rng.Read(pkhBytes[:])
	_ = addr.SetPaymentPubKeyHashFromSlice(pkhBytes[:])
	return addr
}

func MakeRandomSignature(rng *rand.Rand) wallet.Sig {
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
	return GetRandomByteSlice(0, address.PubKeyLength+address.PubKeyHashLength-1, rng)
}

func MakeTooManyAddressBytes(rng *rand.Rand) []byte {
	const maxInvalidAddressLength = address.PubKeyLength + address.PubKeyHashLength*2
	return GetRandomByteSlice(address.PubKeyLength+address.PubKeyHashLength+1, maxInvalidAddressLength, rng)
}

func MakeTooLongSignature(rng *rand.Rand) wallet.Sig {
	const maxInvalidSigLength = 0x100
	return GetRandomByteSlice(wire.SignatureLength+1, maxInvalidSigLength, rng)
}

func MakeTooShortSignature(rng *rand.Rand) wallet.Sig {
	return GetRandomByteSlice(0, wire.SignatureLength-1, rng)
}
