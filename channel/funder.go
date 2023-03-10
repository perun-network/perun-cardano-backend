package channel

import (
	"context"
	"errors"
	"fmt"
	"perun.network/go-perun/channel"
	"perun.network/perun-cardano-backend/channel/types"
)

var (
	MismatchingChannelTokenError = errors.New("mismatching channel tokens")
	MismatchingChannelIDError    = errors.New("mismatching channel ids")
)

type Funder struct {
	sub *AdjudicatorSub
	pab *PAB
}

func NewFunder(sub *AdjudicatorSub, pab *PAB) *Funder {
	return &Funder{
		sub: sub,
		pab: pab,
	}
}

func (f Funder) Fund(ctx context.Context, req channel.FundingReq) error {
	params, err := types.MakeChannelParameters(*req.Params)
	if err != nil {
		return fmt.Errorf("unable to convert channel parameters for funding: %w", err)
	}
	state, err := types.ConvertChannelState(*req.State)
	if err != nil {
		return fmt.Errorf("unable to convert channel state for funding: %w", err)
	}

	for i := uint16(0); i < uint16(req.Idx); i++ {
		if i == 0 {
			err = f.ExpectAndHandleStartEvent(req.Params.ID())
			if err != nil {
				return err
			}
		} else {
			err = f.ExpectAndHandleDepositedEvent(req.Params.ID())
			if err != nil {
				return err
			}
		}
	}
	if uint16(req.Idx) == uint16(0) {
		err = f.pab.Start(req.Params.ID(), params, state)
		if err != nil {
			return err
		}
		err = f.ExpectAndHandleStartEvent(req.Params.ID())
		if err != nil {
			return err
		}

	} else {
		err = f.pab.Fund(req.Params.ID(), req.Idx)
		if err != nil {
			return err
		}
		err = f.ExpectAndHandleDepositedEvent(req.Params.ID())
		if err != nil {
			return err
		}
	}

	for i := int(req.Idx); i < len(params.Parties); i++ {
		err = f.ExpectAndHandleDepositedEvent(req.Params.ID())
		if err != nil {
			return err
		}
	}
	return nil
}

func (f Funder) ExpectAndHandleStartEvent(id types.ID) error {
	event := f.sub.Next()
	if event.ID() != id {
		return MismatchingChannelIDError
	}
	start, ok := event.(*types.Started)
	if !ok {
		//TODO: Handle
		return errors.New("expected Started event")
	}
	err := f.pab.SetChannelToken(start.ID(), start.ChannelDatum.ChannelToken)
	if err != nil {
		return fmt.Errorf("unable to set channel token: %w", err)
	}
	return nil
	//TODO: Verify & Check Start event
}

func (f Funder) ExpectAndHandleDepositedEvent(id types.ID) error {
	event := f.sub.Next()
	if event.ID() != id {
		return MismatchingChannelIDError
	}
	deposited, ok := event.(*types.Deposited)
	if !ok {
		//TODO: Handle
		return errors.New("expected Deposited event")
	}
	token, err := f.pab.GetChannelToken(deposited.ID())
	if err != nil {
		return err
	}
	if token != deposited.ChannelDatum.ChannelToken {
		return MismatchingChannelTokenError
	}
	//TODO: Verify & Check Deposit event
	return nil
}
