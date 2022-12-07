package test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"perun.network/perun-cardano-backend/wallet"
	"testing"
)

func TestBackend_NewAddress(t *testing.T) {
	seed := SetSeed()
	r := NewMockRemote()
	uut := wallet.MakeRemoteBackend(r)
	actualAddress := uut.NewAddress()
	_, ok := actualAddress.(*wallet.Address)
	require.Truef(t, ok, "NewAddress() does not return a Address, test-seed: %d", seed)
}

func TestBackend_DecodeSig(t *testing.T) {
	seed := SetSeed()
	r := NewMockRemote()
	uut := wallet.MakeRemoteBackend(r)

	readerExact := bytes.NewReader(r.MockSignature)
	actualSig, err := uut.DecodeSig(readerExact)
	require.NoErrorf(t, err, "received an error when decoding signature, test-seed: %d", seed)
	require.Equalf(t, r.MockSignature, actualSig, "decoded signature is incorrect, test-seed: %d", seed)

	randomBytes := make([]byte, rand.Intn(129))
	rand.Read(randomBytes)
	readerLonger := bytes.NewReader(append(r.MockSignature, randomBytes...))
	actualSig, err = uut.DecodeSig(readerLonger)
	require.NoErrorf(t, err, "received an error when decoding signature, test-seed: %d", seed)
	require.Equalf(t, r.MockSignature, actualSig, "decoded signature is incorrect, test-seed: %d", seed)
	rest, err := io.ReadAll(readerLonger)
	require.NoErrorf(
		t,
		err,
		"only one signature (%d bytes) should be read from given reader, test-seed: %d",
		wallet.SignatureLength,
		seed,
	)
	require.Equalf(
		t,
		randomBytes,
		rest,
		"only one signature (%d bytes) should be read from given reader. No more should be read from the reader, test-seed: %d",
		wallet.SignatureLength,
		seed,
	)

	invalidReader := bytes.NewReader(r.InvalidSignatureShorter)
	fmt.Println(len(r.InvalidSignatureShorter))
	_, err = uut.DecodeSig(invalidReader)
	require.Errorf(
		t,
		err,
		"did not error when decoding a shorter signature of length: %d, test-seed: %d",
		len(r.InvalidSignatureShorter),
		seed,
	)
}

func TestBackend_VerifySignature(t *testing.T) {
	seed := SetSeed()
	r := NewMockRemote()
	uut := wallet.MakeRemoteBackend(r)
	valid, err := uut.VerifySignature(r.MockMessage, r.MockSignature, &r.MockPubKey)
	require.NoErrorf(t, err, "received error when verifying a valid signature, test-seed: %d", seed)
	require.Truef(t, valid, "did not verify a valid signature as valid, test-seed: %d", seed)

	valid, err = uut.VerifySignature(r.MockMessage, r.OtherSignature, &r.MockPubKey)
	require.NoErrorf(t, err, "received an error when verifying an invalid signature, test-seed: %d", seed)
	require.Falsef(t, valid, "verified an invalid signature as valid, test-seed: %d", seed)

	_, err = uut.VerifySignature(r.MockMessage, r.InvalidSignatureShorter, &r.MockPubKey)
	require.Error(
		t,
		err,
		"failed to error when verifying signature of invalid length: %d, test-seed: %d",
		len(r.InvalidSignatureShorter),
		seed,
	)

	_, err = uut.VerifySignature(r.MockMessage, r.InvalidSignatureLonger, &r.MockPubKey)
	require.Error(
		t,
		err,
		"failed to error when verifying signature of invalid length: %d, test-seed: %d",
		len(r.InvalidSignatureLonger),
		seed,
	)
}
