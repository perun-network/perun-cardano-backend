package wire

import (
	"encoding/hex"
	"perun.network/go-perun/channel"
	"perun.network/perun-cardano-backend/channel/types"
	"perun.network/perun-cardano-backend/wallet/address"
)

// ChannelParameters reflects the Haskell type `Channel` of the Channel Smart Contract in respect to its json encoding.
type ChannelParameters struct {
	Nonce               channel.Nonce       `json:"pNonce"`
	PaymentPubKeyHashes []PaymentPubKeyHash `json:"pPaymentPKs"`
	SigningPubKeys      []PaymentPubKey     `json:"pSigningPKs"`
	TimeLock            int64               `json:"pTimeLock"`
}

type PaymentPubKey struct {
	PubKey PubKey `json:"unPaymentPubKey"`
}

type PaymentPubKeyHash struct {
	PubKeyHash PubKeyHash `json:"unPaymentPubKeyHash"`
}

type PubKeyHash struct {
	Hex string `json:"getPubKeyHash"`
}

func MakePaymentPubKey(address address.Address) PaymentPubKey {
	return PaymentPubKey{
		PubKey: MakePubKey(address),
	}
}

func MakeChannelParameters(parameters types.ChannelParameters) ChannelParameters {
	pubKeyHashes := make([]PaymentPubKeyHash, len(parameters.Parties))
	signingPubKeys := make([]PaymentPubKey, len(parameters.Parties))
	for i, addr := range parameters.Parties {
		pubKeyHashes[i] = MakePaymentPubKeyHash(addr)
		signingPubKeys[i] = MakePaymentPubKey(addr)
	}
	return ChannelParameters{
		Nonce:               parameters.Nonce,
		PaymentPubKeyHashes: pubKeyHashes,
		SigningPubKeys:      signingPubKeys,
		TimeLock:            parameters.Timeout.Milliseconds(),
	}

}

func MakePaymentPubKeyHash(address address.Address) PaymentPubKeyHash {
	hash, err := address.GetPubKeyHash()
	if err != nil {
		panic(err)
	}
	return PaymentPubKeyHash{
		PubKeyHash: PubKeyHash{
			Hex: hex.EncodeToString(hash[:]),
		},
	}
}
