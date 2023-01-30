package wallet_test

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	gptest "perun.network/go-perun/wallet/test"
	"perun.network/perun-cardano-backend/wallet"
	"perun.network/perun-cardano-backend/wallet/test"
	pkgtest "polycry.pt/poly-go/test"
	"testing"
)

func TestRemoteWallet_Trace(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	backend := wallet.MakeRemoteBackend(r)
	address := backend.NewAddress()
	err := address.UnmarshalBinary(r.MockPubKeyBytes[:])
	require.NoError(t, err, "unable to marshal binary address into Address")
	require.Equal(t, &r.MockAddress, address, "unmarshalled address is not as expected")
	require.NotEqual(
		t,
		r.UnavailableAddress,
		address,
		"unmarshalled address is equal to a wrong address",
	)

	err = backend.NewAddress().UnmarshalBinary(r.InvalidPubKeyBytes)
	require.Error(t, err, "unmarshalled invalid binary address into Address")

	binaryAddress, err := address.MarshalBinary()
	require.NoError(t, err, "unable to marshal valid Address into binary")
	require.Equal(
		t,
		r.MockPubKeyBytes[:],
		binaryAddress,
		"marshalled Address is not as expected",
	)

	w := wallet.NewRemoteWallet(r)
	account, err := w.Unlock(address)
	require.NoError(t, err, "failed to unlock valid address")

	_, err = w.Unlock(&r.UnavailableAddress)
	require.Error(t, err, "unlocked an account for an unavailable public key")

	require.Equal(
		t,
		&r.MockAddress,
		account.Address(),
		"account has address with unexpected public key",
	)
	signature, err := account.SignData(r.MockMessage)
	require.NoError(t, err, "unable to sign message")
	require.Equal(t, r.MockSignature, signature, "signature is not as expected")

	valid, err := backend.VerifySignature(r.MockMessage, signature, address)
	require.NoError(t, err, "failed to verify valid signature")
	require.True(t, valid, "valid signature verified as invalid")

	invalid, err := backend.VerifySignature(r.MockMessage, r.OtherSignature, address)
	require.NoError(t, err, "unable to establish validity of invalid signature")
	require.False(t, invalid, "invalid signature verified as valid")
}

func TestRemoteWallet_Unlock(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewMockRemote(rng)
	w := wallet.NewRemoteWallet(r)
	account, err := w.Unlock(&r.MockAddress)
	require.NoError(t, err, "unable to unlock available address")
	require.Equal(t, &r.MockAddress, account.Address(), "wrong address in account")
	require.NotEqual(
		t,
		r.UnavailableAddress,
		account.Address(),
		"Wrong address in account. This is probably because of wrong implementation of Address.Equal",
	)
	_, err = w.Unlock(&r.UnavailableAddress)
	require.Error(
		t,
		err,
		"unlock should fail if the remote wallet does not have the private key to the given address",
	)
}

func setup(rng *rand.Rand) *gptest.Setup {
	r := test.NewGenericRemote(test.MakeRandomAddress(rng), rng)
	w := wallet.NewRemoteWallet(r)
	b := wallet.MakeRemoteBackend(r)
	marshalledAddress, err := test.MakeRandomAddress(rng).MarshalBinary()
	if err != nil {
		panic(err)
	}
	zero := b.NewAddress()
	return &gptest.Setup{
		Backend:           b,
		Wallet:            w,
		AddressInWallet:   &r.AvailableAddress,
		ZeroAddress:       zero,
		AddressMarshalled: marshalledAddress,
	}
}

func TestAddress(t *testing.T) {
	gptest.TestAddress(t, setup(pkgtest.Prng(t)))
}

func TestGenericSignatureSize(t *testing.T) {
	gptest.GenericSignatureSizeTest(t, setup(pkgtest.Prng(t)))
}

func TestAccountWithWalletAndBackend(t *testing.T) {
	gptest.TestAccountWithWalletAndBackend(t, setup(pkgtest.Prng(t)))
}
