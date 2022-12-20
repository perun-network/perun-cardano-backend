package wallet

import (
	"fmt"
	"io"
	"perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/channel/types"
	"perun.network/perun-cardano-backend/wallet/address"
	"perun.network/perun-cardano-backend/wire"
)

// RemoteBackend is a wallet.Backend implementation with a remote server for signing data and verifying signatures.
type RemoteBackend struct {
	walletServer Remote
}

// MakeRemoteBackend returns a new RemoteBackend struct.
func MakeRemoteBackend(remote Remote) RemoteBackend {
	return RemoteBackend{remote}
}

// NewAddress returns a pointer to a new, empty address.
func (b RemoteBackend) NewAddress() wallet.Address {
	return new(address.Address)
}

// DecodeSig reads SignatureLength bytes from the given reader and returns the read signature.
func (b RemoteBackend) DecodeSig(reader io.Reader) (wallet.Sig, error) {
	sig := make([]byte, wire.SignatureLength)
	if _, err := io.ReadFull(reader, sig); err != nil {
		return nil, fmt.Errorf("unable to read signature from reader: %w", err)
	}
	return sig, nil
}

// VerifySignature returns true, iff the given signature is valid for the given message under the public key associated
// with the given address.
func (b RemoteBackend) VerifySignature(msg []byte, sig wallet.Sig, a wallet.Address) (bool, error) {
	addr, ok := a.(*address.Address)
	if !ok {
		return false, fmt.Errorf("invalid PubKey for signature verification")
	}
	if len(sig) != wire.SignatureLength {
		return false, fmt.Errorf(
			"signature has incorrect length. expected: %d bytes actual: %d bytes",
			wire.SignatureLength,
			len(sig),
		)
	}
	request := wire.MakeVerificationRequest(sig, *addr, msg)
	var response wire.VerificationResponse
	err := b.walletServer.CallEndpoint(EndpointVerifyDataSignature, request, &response)
	if err != nil {
		return false, fmt.Errorf("wallet server could not verify message: %w", err)
	}
	return response, nil
}

// VerifyChannelStateSignature returns true, iff the given signature is valid for the given ChannelState under the
// public key associated with the given address.
func (b RemoteBackend) VerifyChannelStateSignature(state types.ChannelState, sig wallet.Sig, a wallet.Address) (bool, error) {
	addr, ok := a.(*address.Address)
	if !ok {
		return false, fmt.Errorf("invalid PubKey for signature verification")
	}
	if len(sig) != wire.SignatureLength {
		return false, fmt.Errorf(
			"signature has incorrect length. expected: %d bytes actual: %d bytes",
			wire.SignatureLength,
			len(sig),
		)
	}
	request := wire.MakeChannelStateVerificationRequest(sig, *addr, state)
	var response wire.VerificationResponse
	err := b.walletServer.CallEndpoint(EndpointVerifyChannelStateSignature, request, &response)
	if err != nil {
		return false, fmt.Errorf("wallet server could not verify message: %w", err)
	}
	return response, nil
}

var _ wallet.Backend = RemoteBackend{}
