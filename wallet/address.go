package wallet

import (
	"bytes"
	"encoding/json"
	"perun.network/go-perun/wallet"
)

type PaymentPubKeyHash struct {
	PubKeyHash SubType `json:"unPaymentPubKeyHash"`
}

type SubType struct {
	HashString string `json:"getPubKeyHash"`
}

func (a PaymentPubKeyHash) MarshalBinary() ([]byte, error) {
	b, err := json.Marshal(a)
	return b, err
}

func (a *PaymentPubKeyHash) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}

func (a PaymentPubKeyHash) String() string {
	binary, err := a.MarshalBinary()
	if err != nil {
		panic("unable to marshal PaymentPubKeyHash")
	}
	return string(binary)
}

func (a PaymentPubKeyHash) Equal(b wallet.Address) bool {
	var err error
	b_, ok := b.(*PaymentPubKeyHash)
	if !ok {
		return false
	}
	aBinary, err := a.MarshalBinary()
	if err != nil {
		panic("unable to marshal PaymentPubKeyHash")
	}
	bBinary, err := b_.MarshalBinary()
	if err != nil {
		return false
	}
	return bytes.Equal(aBinary, bBinary)
}
