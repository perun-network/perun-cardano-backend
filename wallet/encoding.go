package wallet

import (
	"encoding/hex"
	"fmt"
	"perun.network/go-perun/wallet"
)

type SigningRequest struct {
	AccountPubKey PubKey `json:"sPubKey"`
	Message       string `json:"sMessage"`
}

// MakeSigningRequest returns a new SigningRequest.
func MakeSigningRequest(accountPubKey PubKey, message []byte) SigningRequest {
	return SigningRequest{
		AccountPubKey: accountPubKey,
		Message:       hex.EncodeToString(message),
	}
}

type VerificationRequest struct {
	SigWrapper SignatureWrapper `json:"vSignature"`
	PubKey     PubKey           `json:"vPubKey"`
	Message    string           `json:"vMessage"`
}

// MakeVerificationRequest expects the signature to be of length SignatureLength and returns a new VerificationRequest.
func MakeVerificationRequest(sig wallet.Sig, pubKey PubKey, message []byte) VerificationRequest {
	return VerificationRequest{
		SigWrapper: SignatureWrapper{Signature: hex.EncodeToString(sig)},
		PubKey:     pubKey,
		Message:    hex.EncodeToString(message),
	}
}

type SignatureWrapper struct {
	Signature string `json:"getSignature"`
}

type KeyAvailabilityRequest = PubKey

type SigningResponse = SignatureWrapper

func (sr SigningResponse) Decode() (wallet.Sig, error) {
	sig, err := hex.DecodeString(sr.Signature)
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

type VerificationResponse = bool
type KeyAvailabilityResponse = bool
