package wallet_test

import (
	"github.com/stretchr/testify/require"
	"perun.network/perun-cardano-backend/wallet"
	"perun.network/perun-cardano-backend/wallet/test"
	"testing"
)

func TestRemoteWallet_Trace(t *testing.T) {
	seed := test.SetSeed()

	r := test.NewMockRemote()
	backend := wallet.MakeRemoteBackend(r)
	address := backend.NewAddress()
	err := address.UnmarshalBinary(r.MockPubKeyBytes)
	require.NoErrorf(t, err, "unable to marshal binary address into Address, test-seed: %d", seed)
	require.Equalf(t, &r.MockPubKey, address, "unmarshalled address is not as expected, test-seed: %d", seed)
	require.NotEqualf(
		t,
		r.UnavailablePubKey,
		address,
		"unmarshalled address is equal to a wrong address, test-seed: %d",
		seed,
	)

	err = backend.NewAddress().UnmarshalBinary(r.InvalidPubKeyBytes)
	require.Errorf(t, err, "unmarshalled invalid binary address into Address, test-seed: %d", seed)

	binaryAddress, err := address.MarshalBinary()
	require.NoErrorf(t, err, "unable to marshal valid address Address into binary, test-seed: %d", seed)
	require.Equalf(
		t,
		r.MockPubKeyBytes,
		binaryAddress,
		"marshalled Address is not as expected, test-seed: %d",
		seed,
	)

	w := wallet.NewRemoteWallet(r)
	account, err := w.Unlock(address)
	require.NoErrorf(t, err, "failed to unlock valid address, test-seed: %d", seed)

	_, err = w.Unlock(&r.UnavailablePubKey)
	require.Errorf(t, err, "unlocked an account for an unavailable public key, test-seed: %d", seed)

	require.Equalf(
		t,
		&r.MockPubKey,
		account.Address(),
		"address has account has unexpected pubic key, test-seed: %d",
		seed,
	)
	signature, err := account.SignData(r.MockMessage)
	require.NoErrorf(t, err, "unable to sign message, test-seed: %d", seed)
	require.Equalf(t, r.MockSignature, signature, "signature is not as expected, test-seed: %d", seed)

	valid, err := backend.VerifySignature(r.MockMessage, signature, address)
	require.NoErrorf(t, err, "failed to verify valid signature, test-seed: %d", seed)
	require.Truef(t, valid, "valid signature verified as invalid, test-seed: %d", seed)

	invalid, err := backend.VerifySignature(r.MockMessage, r.OtherSignature, address)
	require.NoErrorf(t, err, "unable to establish validity of invalid signature, test-seed: %d", seed)
	require.Falsef(t, invalid, "invalid signature verified as valid, test-seed: %d", seed)
}

func TestRemoteWallet_Unlock(t *testing.T) {
	seed := test.SetSeed()
	r := test.NewMockRemote()
	w := wallet.NewRemoteWallet(r)
	account, err := w.Unlock(&r.MockPubKey)
	require.NoErrorf(t, err, "unable to unlock available address, test-seed: %d", seed)
	require.Equalf(t, &r.MockPubKey, account.Address(), "wrong address in account, test-seed: %d", seed)
	require.NotEqualf(
		t,
		r.UnavailablePubKey,
		account.Address(),
		"Wrong address in account. This is probably because of wrong implementation of Address.Equal, test-seed: %d",
		seed,
	)
	_, err = w.Unlock(&r.UnavailablePubKey)
	require.Errorf(
		t,
		err,
		"unlock should fail if the remote wallet does not have the private key to the given address, test-seed: %d",
		seed,
	)
}
