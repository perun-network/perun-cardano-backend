package channel

import (
	"context"
	"errors"
	"fmt"
	"math"
	"perun.network/go-perun/channel"
	"perun.network/perun-cardano-backend/channel/types"
	"time"
)

var (
	MismatchingChannelTokenError = errors.New("mismatching channel tokens")
	MismatchingChannelIDError    = errors.New("mismatching channel ids")
)

type Funder struct {
	pab *PAB
}

func NewFunder(pab *PAB) *Funder {
	return &Funder{
		pab: pab,
	}
}

func (f Funder) Fund(_ context.Context, req channel.FundingReq) error {
	// TODO implement funding abort (reclamation of funds on peer misbehaviour)
	sub, err := f.pab.NewInternalSubscription(req.Params.ID())
	if err != nil {
		return fmt.Errorf("unable to create subscription: %w", err)
	}
	defer sub.Close()

	params, err := types.MakeChannelParameters(*req.Params.Clone())
	if err != nil {
		return fmt.Errorf("unable to convert channel parameters for funding: %w", err)
	}
	if len(params.Parties) >= math.MaxUint16 {
		return fmt.Errorf("too many parties: max: %d, actual: %d", math.MaxUint16, len(params.Parties))
	}
	state, err := types.ConvertChannelState(*req.State.Clone())
	if err != nil {
		return fmt.Errorf("unable to convert channel state for funding: %w", err)
	}

	for i := uint16(0); i < uint16(req.Idx); i++ {
		if i == 0 {
			err = f.ExpectAndHandleStartEvent(req.Params.ID(), sub, state)
			if err != nil {
				return err
			}
		} else {
			err = f.ExpectAndHandleDepositedEvent(req.Params.ID(), sub, i)
			if err != nil {
				return err
			}
		}
	}
	// Unfortunately, this sleep is necessary to avoid a race in the Adjudicator Subscription due to a slow chain index.
	time.Sleep(5 * time.Second)
	if uint16(req.Idx) == uint16(0) {
		err = f.pab.Start(req.Params.ID(), params, state)
		if err != nil {
			return err
		}
		err = f.ExpectAndHandleStartEvent(req.Params.ID(), sub, state)
		if err != nil {
			return err
		}

	} else {
		err = f.pab.Fund(req.Params.ID(), req.Idx)
		if err != nil {
			return err
		}
		err = f.ExpectAndHandleDepositedEvent(req.Params.ID(), sub, uint16(req.Idx))
		if err != nil {
			return err
		}
	}

	for i := int(req.Idx) + 1; i < len(params.Parties); i++ {
		// Narrowing is safe, because we already checked that the number of parties is smaller than math.MaxUint16
		err = f.ExpectAndHandleDepositedEvent(req.Params.ID(), sub, uint16(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func (f Funder) ExpectAndHandleStartEvent(id types.ID, sub *AdjudicatorSub, state types.ChannelState) error {
	event := sub.Next()
	if event.ID() != id {
		return MismatchingChannelIDError
	}
	start, ok := event.(types.Created)
	if !ok {
		return fmt.Errorf("expected Created event, got type %T, value: %v", event, event)
	}
	err := f.pab.SetChannelToken(start.ID(), start.NewDatum.ChannelToken)
	if err != nil {
		return fmt.Errorf("unable to set channel token: %w", err)
	}

	return verifyStartEvent(start.NewDatum, state)
}

func (f Funder) ExpectAndHandleDepositedEvent(id types.ID, sub *AdjudicatorSub, idx uint16) error {
	event := sub.Next()
	if event.ID() != id {
		return MismatchingChannelIDError
	}
	deposited, ok := event.(types.Deposited)
	if !ok {
		return fmt.Errorf("expected Deposited event, got type %T, value: %v", event, event)
	}
	token, err := f.pab.GetChannelToken(deposited.ID())
	if err != nil {
		return err
	}
	if token != deposited.NewDatum.ChannelToken {
		return MismatchingChannelTokenError
	}
	return verifyFundedEvent(deposited.NewDatum, idx)
}

func verifyStartEvent(outputDatum types.ChannelDatum, state types.ChannelState) error {
	if !outputDatum.ChannelState.Equal(state) {
		return errors.New("on-chain channel state does not match channel state in funding request")
	}
	return verifyFundedEvent(outputDatum, 0)
}

func verifyFundedEvent(outputDatum types.ChannelDatum, idx uint16) error {
	if outputDatum.FundingBalances[idx] != outputDatum.ChannelState.Balances[idx] {
		return fmt.Errorf("party %d did not fund the channel correctly", idx)
	}
	return nil
}
