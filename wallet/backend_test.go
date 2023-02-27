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

package wallet_test

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"perun.network/perun-cardano-backend/wallet"
	"perun.network/perun-cardano-backend/wallet/address"
	"perun.network/perun-cardano-backend/wallet/test"
	"perun.network/perun-cardano-backend/wire"
	pkgtest "polycry.pt/poly-go/test"
	"testing"
)

func TestBackend_NewAddress(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := wallet.MakeRemoteBackend(r)
	actualAddress := uut.NewAddress()
	_, ok := actualAddress.(*address.Address)
	require.True(t, ok, "NewAddress() does not return an Address")
}

func TestBackend_DecodeSig(t *testing.T) {
	const maxRandomBytesLength = 128

	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := wallet.MakeRemoteBackend(r)

	readerExact := bytes.NewReader(r.MockSignature)
	actualSig, err := uut.DecodeSig(readerExact)
	require.NoError(t, err, "received an error when decoding signature")
	require.Equal(t, r.MockSignature, actualSig, "decoded signature is incorrect")

	randomBytes := test.GetRandomByteSlice(0, maxRandomBytesLength, rng)
	readerLonger := bytes.NewReader(append(r.MockSignature, randomBytes...))
	actualSig, err = uut.DecodeSig(readerLonger)
	require.NoError(t, err, "received an error when decoding signature")
	require.Equal(t, r.MockSignature, actualSig, "decoded signature is incorrect")
	rest, err := io.ReadAll(readerLonger)
	require.NoErrorf(
		t,
		err,
		"only one signature (%d bytes) should be read from given reader",
		wire.SignatureLength,
	)
	require.Equalf(
		t,
		randomBytes,
		rest,
		"only one signature (%d bytes) should be read from given reader. No more should be read from the reader",
		wire.SignatureLength,
	)

	invalidReader := bytes.NewReader(r.InvalidSignatureShorter)
	_, err = uut.DecodeSig(invalidReader)
	require.Errorf(
		t,
		err,
		"did not error when decoding a shorter signature of length: %d",
		len(r.InvalidSignatureShorter),
	)
}

func TestBackend_VerifySignature(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := wallet.MakeRemoteBackend(r)
	valid, err := uut.VerifySignature(r.MockMessage, r.MockSignature, &r.MockAddress)
	require.NoError(t, err, "received error when verifying a valid signature")
	require.True(t, valid, "did not verify a valid signature as valid")

	valid, err = uut.VerifySignature(r.MockMessage, r.OtherSignature, &r.MockAddress)
	require.NoError(t, err, "received an error when verifying an invalid signature")
	require.False(t, valid, "verified an invalid signature as valid")

	_, err = uut.VerifySignature(r.MockMessage, r.InvalidSignatureShorter, &r.MockAddress)
	require.Errorf(
		t,
		err,
		"failed to error when verifying signature of invalid length: %d",
		len(r.InvalidSignatureShorter),
	)

	_, err = uut.VerifySignature(r.MockMessage, r.InvalidSignatureLonger, &r.MockAddress)
	require.Errorf(
		t,
		err,
		"failed to error when verifying signature of invalid length: %d",
		len(r.InvalidSignatureLonger),
	)
}

func TestRemoteBackend_VerifyChannelStateSignature(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	uut := wallet.MakeRemoteBackend(r)
	valid, err := uut.VerifyChannelStateSignature(r.MockChannelState, r.MockSignature, &r.MockAddress)
	require.NoError(t, err, "received error when verifying a valid signature")
	require.True(t, valid, "did not verify a valid signature as valid")

	valid, err = uut.VerifyChannelStateSignature(r.MockChannelState, r.OtherSignature, &r.MockAddress)
	require.NoError(t, err, "received an error when verifying an invalid signature")
	require.False(t, valid, "verified an invalid signature as valid")

	_, err = uut.VerifyChannelStateSignature(r.MockChannelState, r.InvalidSignatureShorter, &r.MockAddress)
	require.Errorf(
		t,
		err,
		"failed to error when verifying signature of invalid length: %d",
		len(r.InvalidSignatureShorter),
	)

	_, err = uut.VerifyChannelStateSignature(r.MockChannelState, r.InvalidSignatureLonger, &r.MockAddress)
	require.Errorf(
		t,
		err,
		"failed to error when verifying signature of invalid length: %d",
		len(r.InvalidSignatureLonger),
	)
}
