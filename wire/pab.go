package wire

import (
	"perun.network/go-perun/channel"
	"perun.network/perun-cardano-backend/channel/types"
)

const (
	PerunContractTag = "PerunContract"
	AdjudicatorTag   = "AdjudicatorContract"
)

type ContractInstanceID struct {
	ID string `json:"unContractInstanceId"`
}

func (id ContractInstanceID) Decode() string {
	return id.ID
}

type ContractActivationWallet struct {
	WalletID string `json:"getWalletId"`
}

type PerunContractActivationID struct {
	Tag string `json:"tag"`
}

type PerunActivationBody struct {
	Tag    PerunContractActivationID `json:"caID"`
	Wallet ContractActivationWallet  `json:"caWallet"`
}

type AdjudicatorSubscriptionActivationID struct {
	Tag       string    `json:"tag"`
	ChannelID ChannelID `json:"channelId"`
}

type AdjudicatorSubscriptionActivationBody struct {
	Tag    AdjudicatorSubscriptionActivationID `json:"caID"`
	Wallet ContractActivationWallet            `json:"caWallet"`
}

func MakePerunActivationBody(walletId string) PerunActivationBody {
	return PerunActivationBody{
		Tag: PerunContractActivationID{
			Tag: PerunContractTag,
		},
		Wallet: ContractActivationWallet{
			WalletID: walletId,
		},
	}
}

func MakeAdjudicatorSubscriptionActivationBody(channelId ChannelID, walletId string) AdjudicatorSubscriptionActivationBody {
	return AdjudicatorSubscriptionActivationBody{
		Tag: AdjudicatorSubscriptionActivationID{
			Tag:       AdjudicatorTag,
			ChannelID: channelId,
		},
		Wallet: ContractActivationWallet{
			WalletID: walletId,
		},
	}
}

type OpenParams struct {
	Balances            []uint64            `json:"spBalances"`
	ChannelID           ChannelID           `json:"spChannelId"`
	Nonce               channel.Nonce       `json:"spNonce"`
	PaymentPubKeyHashes []PaymentPubKeyHash `json:"pPaymentPKs"`
	SigningPubKeys      []PaymentPubKey     `json:"pSigningPKs"`
	TimeLock            int64               `json:"pTimeLock"`
}

func MakeOpenParams(id ChannelID, p types.ChannelParameters, s types.ChannelState) OpenParams {
	wp := MakeChannelParameters(p)
	return OpenParams{
		Balances:            s.Balances,
		ChannelID:           id,
		Nonce:               wp.Nonce,
		PaymentPubKeyHashes: wp.PaymentPubKeyHashes,
		SigningPubKeys:      wp.SigningPubKeys,
		TimeLock:            wp.TimeLock,
	}
}

type FundParams struct {
	ChannelID    ChannelID    `json:"fpChannelId"`
	ChannelToken ChannelToken `json:"fpChannelToken"`
	Index        uint16       `json:"fpIndex"`
}

func MakeFundParams(id ChannelID, token types.ChannelToken, index uint16) FundParams {
	return FundParams{
		ChannelID:    id,
		ChannelToken: MakeChannelToken(token),
		Index:        index,
	}
}

type StateSignatures struct {
	ChannelState ChannelState `json:"aChannelState"`
	Signatures   []Signature  `json:"aSignatures"`
}

type DisputeParams struct {
	ChannelID      ChannelID       `json:"dpChannelId"`
	ChannelToken   ChannelToken    `json:"dpChannelToken"`
	SignedState    StateSignatures `json:"dpSignedState"`
	SigningPubKeys []PaymentPubKey `json:"dpSigningPKs"`
}

type CloseParams struct {
	SigningPubKeys []PaymentPubKey `json:"cpSigningPKs"`
	SignedState    StateSignatures `json:"cpSignedState"`
	ChannelToken   ChannelToken    `json:"cpChannelToken"`
	ChannelID      ChannelID       `json:"cpChannelId"`
}

type ForceCloseParams struct {
	ChannelToken ChannelToken `json:"fcpChannelToken"`
	ChannelID    ChannelID    `json:"fcpChannelId"`
}
