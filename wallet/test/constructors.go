package test

import (
	"perun.network/perun-cardano-backend/wallet"
	"perun.network/perun-cardano-backend/wallet/address"
)

func MakeRemoteAccount(address address.Address, remote wallet.Remote) wallet.RemoteAccount {
	return wallet.MakeRemoteAccount(address, remote, "")
}

func NewRemoteWallet(remote wallet.Remote) *wallet.RemoteWallet {
	return wallet.NewRemoteWallet(remote, "")
}
