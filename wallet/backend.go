package wallet

import (
	"fmt"
	"io"
	"perun.network/go-perun/wallet"
)

// SignatureLength is the length of valid Cardano signatures in bytes.
const SignatureLength = 64

// RemoteBackend is a wallet.Backend implementation with a remote server for signing data and verifying signatures.
type RemoteBackend struct {
	walletServer Remote
}

// MakeRemoteBackend returns a new RemoteBackend struct.
func MakeRemoteBackend(remote Remote) RemoteBackend {
	return RemoteBackend{remote}
}

// NewAddress returns a pointer to a new empty PubKey.
func (b RemoteBackend) NewAddress() wallet.Address {
	return new(PubKey)
}

// DecodeSig reads SignatureLength bytes from the given reader and returns the read signature.
func (b RemoteBackend) DecodeSig(reader io.Reader) (wallet.Sig, error) {
	sig := make([]byte, SignatureLength)
	if _, err := io.ReadFull(reader, sig); err != nil {
		return nil, fmt.Errorf("unable to read signature from reader: %w", err)
	}
	return sig, nil
}

// VerifySignature returns true, iff the given signature is valid for the given message under the public key associated
// with the given address.
func (b RemoteBackend) VerifySignature(msg []byte, sig wallet.Sig, a wallet.Address) (bool, error) {
	pubKey, ok := a.(*PubKey)
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
	request := MakeVerificationRequest(sig, *pubKey, msg)
	verificationResponse, err := b.walletServer.CallVerify(request)
	if err != nil {
		return false, fmt.Errorf("wallet server could not verify message: %w", err)
	}
	return verificationResponse, nil
}

var _ wallet.Backend = RemoteBackend{}
