package wire

import (
	"encoding/hex"
	"fmt"
	"perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/wallet/address"
)

// SignatureLength is the length of valid Cardano signatures in bytes.
const SignatureLength = 64

// SigningRequest is the json serializable request for signing via the perun-cardano-wallet api.
type SigningRequest struct {
	PubKey  PubKey `json:"sPubKey"`
	Message string `json:"sMessage"`
}

// MakeSigningRequest returns a new SigningRequest.
func MakeSigningRequest(address address.Address, message []byte) SigningRequest {
	return SigningRequest{
		PubKey:  MakePubKey(address),
		Message: hex.EncodeToString(message),
	}
}

// VerificationRequest is the json serializable request for verifying via the perun-cardano-wallet api.
type VerificationRequest struct {
	Signature Signature `json:"vSignature"`
	PubKey    PubKey    `json:"vPubKey"`
	Message   string    `json:"vMessage"`
}

// MakeVerificationRequest returns a new VerificationRequest.
func MakeVerificationRequest(sig wallet.Sig, address address.Address, message []byte) VerificationRequest {
	return VerificationRequest{
		Signature: MakeSignature(sig),
		PubKey:    MakePubKey(address),
		Message:   hex.EncodeToString(message),
	}
}

// Signature is the json serialization for the cardano signature type (see: Ledger.Crypto.Signature).
type Signature struct {
	Hex string `json:"getSignature"`
}

// MakeSignature returns a new Signature. Note that this does not check the length of the received wallet.Sig.
func MakeSignature(sig wallet.Sig) Signature {
	return Signature{Hex: hex.EncodeToString(sig)}
}

// KeyAvailabilityRequest is the json serializable request for key-availability via the perun-cardano-wallet api.
type KeyAvailabilityRequest = PubKey

// MakeKeyAvailabilityRequest returns a new KeyAvailabilityRequest for the given address.
func MakeKeyAvailabilityRequest(address address.Address) KeyAvailabilityRequest {
	return MakePubKey(address)
}

// SigningResponse is the json serializable response when signing via the perun-cardano-wallet api.
type SigningResponse = Signature

// Decode decodes the siganture from a SigningResponse.
func (sr SigningResponse) Decode() (wallet.Sig, error) {
	sig, err := hex.DecodeString(sr.Hex)
	if err != nil {
		return nil, fmt.Errorf("unable to decode Signature from hex string: %w", err)
	}
	if len(sig) != SignatureLength {
		return nil, fmt.Errorf(
			"signature has incorrect length. expected: %d bytes, actual: %d bytes",
			SignatureLength,
			len(sig),
		)
	}
	return sig, nil
}

// VerificationResponse is the json serializable response when verifying signatures via the perun-cardano-wallet api.
type VerificationResponse = bool

// KeyAvailabilityResponse is json serializable response when requesting key-availability via the
// perun-cardano-wallet api.
type KeyAvailabilityResponse = bool
