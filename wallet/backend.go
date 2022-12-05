package wallet

import (
	"encoding/hex"
	"fmt"
	"io"
	"perun.network/go-perun/wallet"
)

// SignatureLength is the length of valid Cardano signatures in bytes
const SignatureLength = 64

// RemoteBackend is a wallet.Backend implementation with a remote server for signing data and verifying signatures
type RemoteBackend struct {
	walletServer Remote
}

// MakeRemoteBackend a new RemoteBackend struct, setting the wallet server appropriately
func MakeRemoteBackend(remote Remote) RemoteBackend {
	return RemoteBackend{remote}
}

// NewAddress returns a pointer to a new empty PubKey
func (b RemoteBackend) NewAddress() wallet.Address {
	return new(PubKey)
}

// DecodeSig tries to read SignatureLength bytes from the reader and returns an error if the given reader does not
// supply enough bytes. The signature or the wrapped error is returned
func (b RemoteBackend) DecodeSig(reader io.Reader) (wallet.Sig, error) {
	sig := make([]byte, SignatureLength)
	if _, err := io.ReadFull(reader, sig); err != nil {
		return nil, fmt.Errorf("unable to read signature from reader: %w", err)
	}
	return sig, nil
}

// VerifySignature first checks whether the give address is a PubKey and the given signature is SignatureLength bytes
// long. It then forms a VerificationRequest and uses it to call the verification endpoint of the remote wallet server.
// The VerificationResponse is decoded and the validity is returned
func (b RemoteBackend) VerifySignature(msg []byte, sig wallet.Sig, a wallet.Address) (bool, error) {
	a_, ok := a.(*PubKey)
	if !ok {
		return false, fmt.Errorf("invalid Address for signature verification")
	}
	if len(sig) != SignatureLength {
		return false, fmt.Errorf(
			"signature has incorrect length. expected: %d bytes actual: %d bytes",
			SignatureLength,
			len(sig),
		)
	}
	request := VerificationRequest{
		SigWrapper: SignatureWrapper{hex.EncodeToString(sig)},
		Key:        *a_,
		Data:       hex.EncodeToString(msg),
	}

	verificationResponse, err := b.walletServer.CallVerify(request)
	if err != nil {
		return false, fmt.Errorf("wallet server could not verify message: %w", err)
	}
	return verificationResponse, nil
}

var _ wallet.Backend = RemoteBackend{}
