package wallet

import (
	"io"
	"perun.network/go-perun/wallet"
)

type Backend struct{}

func (b Backend) NewAddress() wallet.Address {
	a := new(PaymentPubKeyHash)
	return a
}

func (b Backend) DecodeSig(reader io.Reader) (wallet.Sig, error) {
	//TODO implement me
	panic("implement me")

}

func (b Backend) VerifySignature(msg []byte, sign wallet.Sig, a wallet.Address) (bool, error) {
	//TODO implement me
	panic("implement me")
}

var _ wallet.Backend = Backend{}
