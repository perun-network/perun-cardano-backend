// Copyright 2023 - See NOTICE file for copyright holders.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package channel_test

import (
	"math/big"
	"math/rand"
	gpchannel "perun.network/go-perun/channel"
	gptest "perun.network/go-perun/channel/test"
	gpwallet "perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/channel"
	"perun.network/perun-cardano-backend/channel/types"
	"perun.network/perun-cardano-backend/wallet"
	"perun.network/perun-cardano-backend/wallet/address"
	"perun.network/perun-cardano-backend/wallet/test"
	pkgtest "polycry.pt/poly-go/test"
	"testing"
)

func setup(rng *rand.Rand) *gptest.Setup {

	w := test.NewRemoteWallet(GenericTestRemote)
	acc, err := w.Unlock(&GenericTestRemote.AvailableAddresses[0])
	newRandomAddress := func() gpwallet.Address {
		rAddr := test.MakeRandomAddress(rng)
		return &rAddr
	}
	newParamsAndState := func(opts ...gptest.RandomOpt) (*gpchannel.Params, *gpchannel.State) {
		return gptest.NewRandomParamsAndState(
			rng,
			gptest.WithoutApp().
				Append(gptest.WithParts(newRandomAddress(), newRandomAddress())).
				Append(gptest.WithLedgerChannel(true)).
				Append(gptest.WithVirtualChannel(false)).
				Append(gptest.WithAssets(types.Asset)).
				Append(gptest.WithBalancesInRange(new(big.Int).SetUint64(0), types.MaxBalance)).
				Append(opts...),
		)
	}
	p1, s1 := newParamsAndState()
	// We need this because `s1` and `s2` must differ IN EVERY FIELD for the verification tests to work.
	p2, s2 := newParamsAndState(gptest.WithIsFinal(!s1.IsFinal))
	if err != nil {
		panic(err)
	}
	return &gptest.Setup{
		Params:        p1,
		Params2:       p2,
		State:         s1,
		State2:        s2,
		Account:       acc,
		RandomAddress: newRandomAddress,
	}
}

var GenericTestRemote *test.GenericRemote

type main struct{}

func (main) Name() string {
	return "TestMain"
}

func TestMain(m *testing.M) {
	// Setting up the go-perun backend values for the generic tests.
	rng := pkgtest.Prng(main{})
	addr := test.MakeRandomAddress(rng)
	GenericTestRemote = test.NewGenericRemote([]address.Address{addr}, rng)
	wb := wallet.MakeRemoteBackend(GenericTestRemote)
	gpwallet.SetBackend(wb)
	channel.SetWalletBackend(wb)
	gpchannel.SetBackend(channel.Backend)

	m.Run()
}

func TestBackend(t *testing.T) {
	rng := pkgtest.Prng(t)
	gptest.GenericBackendTest(t, setup(rng), gptest.IgnoreApp, gptest.IgnoreAssets)
}
