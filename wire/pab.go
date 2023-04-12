package wire

import (
	"perun.network/go-perun/wallet"
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
	ChannelID ChannelID `json:"contents"`
}

type AdjudicatorSubscriptionActivationBody struct {
	CaID   AdjudicatorSubscriptionActivationID `json:"caID"`
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
		CaID: AdjudicatorSubscriptionActivationID{
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
	Nonce               string              `json:"spNonce"`
	PaymentPubKeyHashes []PaymentPubKeyHash `json:"spPaymentPKs"`
	SigningPubKeys      []PaymentPubKey     `json:"spSigningPKs"`
	TimeLock            int64               `json:"spTimeLock"`
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
	ChannelID    ChannelID  `json:"fpChannelId"`
	ChannelToken AssetClass `json:"fpChannelToken"`
	Index        uint16     `json:"fpIndex"`
}

func MakeFundParams(id ChannelID, token types.ChannelToken, index uint16) FundParams {
	return FundParams{
		ChannelID:    id,
		ChannelToken: MakeAssetClass(token),
		Index:        index,
	}
}

type StateSignatures struct {
	ChannelState ChannelState `json:"aChannelState"`
	Signatures   []Signature  `json:"aSignatures"`
}

func MakeStateSignatures(state types.ChannelState, sigs []wallet.Sig) StateSignatures {
	ws := make([]Signature, len(sigs))
	for i, s := range sigs {
		ws[i] = MakeSignature(s)
	}
	return StateSignatures{
		ChannelState: MakeChannelState(state),
		Signatures:   ws,
	}
}

type DisputeParams struct {
	ChannelID      ChannelID       `json:"dpChannelId"`
	ChannelToken   AssetClass      `json:"dpChannelToken"`
	SignedState    StateSignatures `json:"dpSignedState"`
	SigningPubKeys []PaymentPubKey `json:"dpSigningPKs"`
}

type CloseParams struct {
	SigningPubKeys []PaymentPubKey `json:"cpSigningPKs"`
	SignedState    StateSignatures `json:"cpSignedState"`
	ChannelToken   AssetClass      `json:"cpChannelToken"`
	ChannelID      ChannelID       `json:"cpChannelId"`
}

func MakeCloseParams(id ChannelID, token types.ChannelToken, params types.ChannelParameters, state types.ChannelState, sigs []wallet.Sig) CloseParams {
	wp := MakeChannelParameters(params)

	return CloseParams{
		SigningPubKeys: wp.SigningPubKeys,
		SignedState:    MakeStateSignatures(state, sigs),
		ChannelToken:   MakeAssetClass(token),
		ChannelID:      id,
	}
}

type ForceCloseParams struct {
	ChannelToken AssetClass `json:"fcpChannelToken"`
	ChannelID    ChannelID  `json:"fcpChannelId"`
}
