package wire

import (
	"encoding/hex"
	"fmt"
	"perun.network/perun-cardano-backend/wallet/address"
)

// PubKey is a json serializable public key to communicate with cardano apis (see: Ledger.Crypto.PubKey).
type PubKey struct {
	PubKey string `json:"getPubKey"`
}

// MakePubKey returns a PubKey
func MakePubKey(address address.Address) PubKey {
	return PubKey{
		PubKey: hex.EncodeToString(address.GetPubKeySlice()),
	}
}

func (p PubKey) Decode() (address.Address, error) {
	pubKey, err := hex.DecodeString(p.PubKey)
	if err != nil {
		return address.Address{}, fmt.Errorf("unable to decode PubKey hex string: %w", err)
	}
	return address.MakeAddressFromByteSlice(pubKey)
}
