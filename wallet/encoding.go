package wallet

type SigningRequest struct {
	Key  PubKey `json:"sPubKey"`
	Data string `json:"sMessage"`
}

type VerificationRequest struct {
	SigWrapper SignatureWrapper `json:"vSignature"`
	Key        PubKey           `json:"vPubKey"`
	Data       string           `json:"vMessage"`
}

type SignatureWrapper struct {
	Signature string `json:"getSignature"`
}

type VerificationResponse = bool
type SigningResponse = SignatureWrapper
