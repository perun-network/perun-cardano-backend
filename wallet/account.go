package wallet

import (
	"encoding/hex"
	"fmt"
	"perun.network/go-perun/wallet"
)

// RemoteAccount represents a cardano account. The secrets are stored on the associated remote walletServer
type RemoteAccount struct {
	Addr         PubKey
	walletServer Remote
}

// Address returns the PubKey associated with this account
func (a RemoteAccount) Address() wallet.Address {
	return &a.Addr
}

// SignData signs arbitrary data through the remote wallet server communicating with the local client using
// SigningRequest and SigningResponse
func (a RemoteAccount) SignData(data []byte) ([]byte, error) {
	// prepare SigningRequest for wallet server
	request := SigningRequest{
		Key:  a.Addr,
		Data: hex.EncodeToString(data),
	}

	signatureResponse, err := a.walletServer.CallSign(request)
	if err != nil {
		return nil, fmt.Errorf("wallet server could not sign message: %w", err)
	}

	// extract and decode signature from SigningResponse
	res, err := hex.DecodeString(signatureResponse.Signature)
	if err != nil {
		return nil, fmt.Errorf("unable to decode signature as hex, %w", err)
	}
	if len(res) != SignatureLength {
		return nil, fmt.Errorf(
			"signature has incorrect length. expected: %d bytes actual: %d bytes",
			SignatureLength,
			len(res),
		)
	}
	return res, nil
}

var _ wallet.Account = RemoteAccount{}
