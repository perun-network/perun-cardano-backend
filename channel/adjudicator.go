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

package channel

import (
	"context"
	"perun.network/go-perun/channel"
	"perun.network/perun-cardano-backend/channel/types"
)

type Adjudicator struct {
	pab *PAB
}

func NewAdjudicator(pab *PAB) *Adjudicator {
	return &Adjudicator{
		pab: pab,
	}
}

func (a Adjudicator) Register(ctx context.Context, req channel.AdjudicatorReq, states []channel.SignedState) error {
	//TODO implement dishonest case
	panic("implement me")
}

func (a Adjudicator) Withdraw(ctx context.Context, req channel.AdjudicatorReq, stateMap channel.StateMap) error {
	params, err := types.MakeChannelParameters(*req.Params.Clone())
	if err != nil {
		return err
	}
	state, err := types.ConvertChannelState(*req.Tx.State.Clone())
	if err != nil {
		return err
	}
	// Note: This assumes the channel-close endpoint to behave like "try-close".
	return a.pab.Close(req.Params.ID(), params, state, req.Tx.Sigs)
}

func (a Adjudicator) Progress(ctx context.Context, req channel.ProgressReq) error {
	//TODO implement dishonest case
	panic("implement me")
}

func (a Adjudicator) Subscribe(ctx context.Context, id channel.ID) (channel.AdjudicatorSubscription, error) {
	return a.pab.NewPerunEventSubscription(id)
}
