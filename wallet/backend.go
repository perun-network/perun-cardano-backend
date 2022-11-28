package wallet

import (
	"encoding/hex"
	"fmt"
	"io"
	"perun.network/go-perun/wallet"
)

const SignatureSize = 64

type Backend struct {
	walletServer Remote
}

func (b Backend) NewAddress() wallet.Address {
	a := new(PubKey)
	return a
}

func (b Backend) DecodeSig(reader io.Reader) (wallet.Sig, error) {
	sig := make([]byte, SignatureSize)
	if _, err := io.ReadFull(reader, sig); err != nil {
		return nil, err
	}
	return sig, nil
}

func (b Backend) VerifySignature(msg []byte, sig wallet.Sig, a wallet.Address) (bool, error) {
	a_, ok := a.(*PubKey)
	if !ok {
		return false, fmt.Errorf("invalid Address for signature verification")
	}
	request := VerificationRequest{
		SigWrapper: SignatureWrapper{hex.EncodeToString(sig)},
		Key:        *a_,
		Data:       hex.EncodeToString(msg),
	}
	if len(sig) != SignatureSize {
		return false, fmt.Errorf(
			"signature has incorrect length. expected: %d bytes actual: %d bytes",
			SignatureSize,
			len(sig),
		)
	}

	verificationResponse, err := b.walletServer.CallVerify(request)
	if err != nil {
		return false, fmt.Errorf("wallet server could not verify message: %w", err)
	}
	return verificationResponse, nil
}

var _ wallet.Backend = Backend{}
