package wallet_test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
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
	require.True(t, ok, "NewAddress() does not return an PubKey")
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

	randomBytes := make([]byte, rand.Intn(maxRandomBytesLength+1))
	rand.Read(randomBytes)
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
	fmt.Println(len(r.InvalidSignatureShorter))
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
