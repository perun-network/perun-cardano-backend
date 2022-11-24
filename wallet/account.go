package wallet

import (
	"perun.network/go-perun/wallet"
)

type Account struct {
	addr PaymentPubKeyHash
}

func (a Account) Address() wallet.Address {
	return &a.addr
}

func (a Account) SignData(data []byte) ([]byte, error) {
	// call server to sign data
	return nil, nil
}
