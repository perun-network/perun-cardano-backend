package wallet

import (
	"fmt"
	"perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/wallet/address"
	"perun.network/perun-cardano-backend/wire"
)

// RemoteAccount represents a cardano account. The secrets are stored on the associated remote walletServer.
type RemoteAccount struct {
	AccountAddress address.Address
	walletServer   Remote
}

// MakeRemoteAccount returns a new RemoteAccount instance.
func MakeRemoteAccount(addr address.Address, r Remote) RemoteAccount {
	return RemoteAccount{
		AccountAddress: addr,
		walletServer:   r,
	}
}

// Address returns the Address associated with this account.
func (a RemoteAccount) Address() wallet.Address {
	return &a.AccountAddress
}

// SignData signs arbitrary data with this account.
func (a RemoteAccount) SignData(data []byte) (wallet.Sig, error) {
	request := wire.MakeSigningRequest(a.AccountAddress, data)

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
