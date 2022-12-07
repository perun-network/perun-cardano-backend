package wallet

import (
	"fmt"
	"perun.network/go-perun/wallet"
)

// RemoteAccount represents a cardano account. The secrets are stored on the associated remote walletServer.
type RemoteAccount struct {
	AccountPubKey PubKey
	walletServer  Remote
}

// MakeRemoteAccount returns a new RemoteAccount instance.
func MakeRemoteAccount(pubKey PubKey, r Remote) RemoteAccount {
	return RemoteAccount{
		AccountPubKey: pubKey,
		walletServer:  r,
	}
}

// Address returns the PubKey associated with this account.
func (a RemoteAccount) Address() wallet.Address {
	return &a.AccountPubKey
}

// SignData signs arbitrary data with this account.
func (a RemoteAccount) SignData(data []byte) (wallet.Sig, error) {
	request := MakeSigningRequest(a.AccountPubKey, data)

	signatureResponse, err := a.walletServer.CallSign(request)
	if err != nil {
		return nil, fmt.Errorf("wallet server could not sign message: %w", err)
	}

	// Extract and decode the signature from SigningResponse.
	sig, err := signatureResponse.Decode()
	if err != nil {
		return nil, fmt.Errorf("unable to decode signature from SignatureResponse: %w", err)
	}
	return sig, nil
}

var _ wallet.Account = RemoteAccount{}
