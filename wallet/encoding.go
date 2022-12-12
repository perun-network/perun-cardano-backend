package wallet

import (
	"encoding/hex"
	"fmt"
	"perun.network/go-perun/wallet"
)

type SigningRequest struct {
	AccountPubKey Address `json:"sPubKey"`
	Message       string  `json:"sMessage"`
}

// MakeSigningRequest returns a new SigningRequest.
func MakeSigningRequest(accountPubKey Address, message []byte) SigningRequest {
	return SigningRequest{
		AccountPubKey: accountPubKey,
		Message:       hex.EncodeToString(message),
	}
}

type VerificationRequest struct {
	SigWrapper SignatureWrapper `json:"vSignature"`
	Address    Address          `json:"vPubKey"`
	Message    string           `json:"vMessage"`
}

// MakeVerificationRequest expects the signature to be of length SignatureLength and returns a new VerificationRequest.
func MakeVerificationRequest(sig wallet.Sig, address Address, message []byte) VerificationRequest {
	return VerificationRequest{
		SigWrapper: SignatureWrapper{Signature: hex.EncodeToString(sig)},
		Address:    address,
		Message:    hex.EncodeToString(message),
	}
}

type SignatureWrapper struct {
	Signature string `json:"getSignature"`
}

type KeyAvailabilityRequest = Address

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
