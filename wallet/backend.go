package wallet

import (
	"encoding/hex"
	"fmt"
	"io"
	"perun.network/go-perun/wallet"
)

// SignatureLength is the length of valid Cardano signatures in bytes
const SignatureLength = 64

type RemoteBackend struct {
	walletServer Remote
}

func MakeRemoteBackend(remote Remote) RemoteBackend {
	return RemoteBackend{remote}
}

func (b RemoteBackend) NewAddress() wallet.Address {
	a := new(PubKey)
	return a
}

func (b RemoteBackend) DecodeSig(reader io.Reader) (wallet.Sig, error) {
	sig := make([]byte, SignatureLength)
	if _, err := io.ReadFull(reader, sig); err != nil {
		return nil, err
	}
	return sig, nil
}

func (b RemoteBackend) VerifySignature(msg []byte, sig wallet.Sig, a wallet.Address) (bool, error) {
	a_, ok := a.(*PubKey)
	if !ok {
		return false, fmt.Errorf("invalid Address for signature verification")
	}
	request := VerificationRequest{
		SigWrapper: SignatureWrapper{hex.EncodeToString(sig)},
		Key:        *a_,
		Data:       hex.EncodeToString(msg),
	}
	if len(sig) != SignatureLength {
		return false, fmt.Errorf(
			"signature has incorrect length. expected: %d bytes actual: %d bytes",
			SignatureLength,
			len(sig),
		)
	}

	verificationResponse, err := b.walletServer.CallVerify(request)
	if err != nil {
		return false, fmt.Errorf("wallet server could not verify message: %w", err)
	}
	return verificationResponse, nil
}

var _ wallet.Backend = RemoteBackend{}
