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

package address_test

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wallet/address"
	"perun.network/perun-cardano-backend/wallet/test"
	pkgtest "polycry.pt/poly-go/test"
	"testing"
)

func TestAddress_MarshalBinary(t *testing.T) {
	rng := pkgtest.Prng(t)
	uut := test.MakeRandomAddress(rng)
	actualBytes, err := uut.MarshalBinary()
	require.NoError(t, err, "unable to marshal valid address")
	require.Equal(
		t,
		append(uut.GetPubKeySlice(), uut.GetPubKeyHashSlice()...),
		actualBytes,
		"wrong bytes representation of marshalled address")
}

func TestAddress_UnmarshalBinary(t *testing.T) {
	rng := pkgtest.Prng(t)
	validTest := func() func(*testing.T) {
		referenceAddress := test.MakeRandomAddress(rng)

		return func(t *testing.T) {
			t.Parallel()
			uut := address.Address{}
			err := uut.UnmarshalBinary(append(referenceAddress.GetPubKeySlice(), referenceAddress.GetPubKeyHashSlice()...))
			require.NoError(t, err, "unable to unmarshal valid address bytes")
			require.Equal(
				t,
				referenceAddress,
				uut,
				"marshalled address is not as expected",
			)
		}
	}
	invalidTest := func(bytesOfInvalidLength []byte) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			uut := address.Address{}
			err := uut.UnmarshalBinary(bytesOfInvalidLength)
			require.Errorf(
				t,
				err,
				"failed to error when unmarshalling invalid public key bytes with length: %d",
				len(bytesOfInvalidLength),
			)
		}
	}
	for i := 0; i < 100; i++ {
		t.Run("Valid", validTest())
		t.Run("Invalid - too few bytes", invalidTest(test.MakeTooFewAddressBytes(rng)))
		t.Run("Invalid - too many bytes", invalidTest(test.MakeTooManyAddressBytes(rng)))
	}
}

func TestAddress_String(t *testing.T) {
	rng := pkgtest.Prng(t)
	uut := test.MakeRandomAddress(rng)
	require.Equal(
		t,
		hex.EncodeToString(uut.GetPubKeySlice()),
		uut.String(),
		"wrong string representation for public key",
	)
}

func TestAddress_Equal(t *testing.T) {
	rng := pkgtest.Prng(t)
	a := test.MakeRandomAddress(rng)
	equalToA := address.MakeAddressFromPubKeyByteArray(a.GetPubKey())
	equalToA.SetPaymentPubKeyHash(a.GetPubKeyHash())

	// Get an address with a strictly different public key to a's.
	differentToA := test.MakeRandomAddress(rng)
	for differentToA.GetPubKey() == a.GetPubKey() {
		differentToA = test.MakeRandomAddress(rng)
	}

	require.True(t, a.Equal(&equalToA), "addresses that have the same public key should be equal")
	require.True(t, equalToA.Equal(&a), "address equality should be commutative")
	require.False(
		t,
		a.Equal(&differentToA),
		"addresses with different public keys should not be equal",
	)
	require.False(
		t,
		differentToA.Equal(&a),
		"address equality should be commutative",
	)
	c := test.NewOtherAddressImpl(t)
	require.False(
		t,
		a.Equal(c),
		"addresses should not be equal to address of different type",
	)
}

func TestAddress_GetTestnetAddress(t *testing.T) {
	const testPubKeyString = "eb94e8236e2099357fa499bfbc415968691573f25ec77435b7949f5fdfaa5da0"
	const expected = "addr_test1vru2drx33ev6dt8gfq245r5k0tmy7ngqe79va69de9dxkrg09c7d3"
	addrBytes, err := hex.DecodeString(testPubKeyString)
	require.NoErrorf(t, err, "this should not fail!")
	addr, err := address.MakeAddressFromPubKeyByteSlice(addrBytes)
	require.NoErrorf(t, err, "unable to create address from byte slice")
	actual, err := addr.GetTestnetAddressOfPubKey()
	require.NoErrorf(t, err, "unexpected error when deriving address string from address")
	require.Equal(t, expected, actual, "address string is not as expected")
}

func TestAddress_GetMainnetAddress(t *testing.T) {
	const testPubKeyString = "eb94e8236e2099357fa499bfbc415968691573f25ec77435b7949f5fdfaa5da0"
	const expected = "addr1v8u2drx33ev6dt8gfq245r5k0tmy7ngqe79va69de9dxkrg5dvzz5"
	addrBytes, err := hex.DecodeString(testPubKeyString)
	require.NoErrorf(t, err, "this should not fail!")
	addr, err := address.MakeAddressFromPubKeyByteSlice(addrBytes)
	require.NoErrorf(t, err, "unable to create address from byte slice")
	actual, err := addr.GetMainnetAddressOfPubKey()
	require.NoErrorf(t, err, "unexpected error when deriving address string from address")
	require.Equal(t, expected, actual, "address string is not as expected")
}

func TestAddress_Calculate(t *testing.T) {
	const testPubKeyString = "eb94e8236e2099357fa499bfbc415968691573f25ec77435b7949f5fdfaa5da0"
	const expectedPubKeyHashString = "f8a68cd18e59a6ace848155a0e967af64f4d00cf8acee8adc95a6b0d"
	expected, err := hex.DecodeString(expectedPubKeyHashString)
	require.NoErrorf(t, err, "this should not fail!")
	addrBytes, err := hex.DecodeString(testPubKeyString)
	require.NoErrorf(t, err, "this should not fail!")
	addr, err := address.MakeAddressFromPubKeyByteSlice(addrBytes)
	require.NoErrorf(t, err, "unable to create address from byte slice")
	actual, err := address.CalculatePubKeyHash(addr.GetPubKey())
	require.NoErrorf(t, err, "unexpected error when deriving PubKeyHash from address")
	require.Equal(t, expected, actual[:], "PubKeyHash is not as expected")
}
