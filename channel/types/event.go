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

package types

import (
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wallet"
	"perun.network/perun-cardano-backend/wire"
)

const (
	CreatedTag   = "Created"
	DepositedTag = "Deposited"
	DisputedTag  = "Disputed"
	ConcludedTag = "Concluded"
)

// TODO: Figure out what to return on AdjudicatorEvent.Version(), if there is no concept of a state version for that
// event. E.g.: The `Concluded` event only has access to the on-chain state prior to conclusion, which might not be
// the version that is eventually concluded.

type InternalEvent interface {
	channel.AdjudicatorEvent
	ToPerunEvent() channel.AdjudicatorEvent
}

type (
	Created struct {
		ChannelID ID
		NewDatum  ChannelDatum
	}
	Deposited struct {
		ChannelID ID
		OldDatum  ChannelDatum
		NewDatum  ChannelDatum
	}
	Disputed struct {
		ChannelID  ID
		OldDatum   ChannelDatum
		NewDatum   ChannelDatum
		Signatures []wallet.Sig
	}
	Concluded struct {
		ChannelID ID
		OldDatum  ChannelDatum
	}
)

func (c Concluded) ID() channel.ID {
	return c.ChannelID
}

func (c Concluded) Timeout() channel.Timeout {
	return nil
}

func (c Concluded) Version() uint64 {
	return 0
}

func (c Concluded) ToPerunEvent() channel.AdjudicatorEvent {
	return channel.NewConcludedEvent(c.ID(), c.Timeout(), c.Version())
}

func (c Concluded) FromEvent(id ID, ev wire.Event) (Concluded, error) {
	if len(ev.DatumList) != 1 {
		return c, NewDecodeEventError(ConcludedTag, 1, len(ev.DatumList))
	}
	oldDatum, err := ev.DatumList[0].Decode()
	if err != nil {
		return c, err
	}
	return Concluded{
		ChannelID: id,
		OldDatum:  oldDatum,
	}, nil
}

func (d Disputed) ID() channel.ID {
	return d.ChannelID
}

func (d Disputed) Timeout() channel.Timeout {
	return &channel.TimeTimeout{Time: d.NewDatum.Time.Add(d.NewDatum.ChannelParameters.Timeout)}
}

func (d Disputed) Version() uint64 {
	return d.NewDatum.ChannelState.Version
}

func (d Disputed) ToPerunEvent() channel.AdjudicatorEvent {
	return channel.NewRegisteredEvent(
		d.ID(),
		d.Timeout(),
		d.Version(),
		nil,          // state is only needed for virtual channels, which we currently do not support anyway
		d.Signatures, // signatures are only needed for virtual channels, but it can not hurt to include them here
	)
}

func (d Disputed) FromEvent(id ID, ev wire.Event) (Disputed, error) {
	if len(ev.DatumList) != 2 {
		return d, DecodeEventError{DisputedTag, 2, len(ev.DatumList)}
	}
	oldDatum, err := ev.DatumList[0].Decode()
	if err != nil {
		return d, err
	}
	newDatum, err := ev.DatumList[1].Decode()
	if err != nil {
		return d, err
	}
	return Disputed{
		ChannelID: id,
		OldDatum:  oldDatum,
		NewDatum:  newDatum,
	}, nil
}

func (d Deposited) ID() channel.ID {
	return d.ChannelID
}

func (d Deposited) Timeout() channel.Timeout {
	return nil
}

func (d Deposited) Version() uint64 {
	return 0
}

func (d Deposited) ToPerunEvent() channel.AdjudicatorEvent {
	return nil
}

func (d Deposited) FromEvent(id ID, ev wire.Event) (Deposited, error) {
	if len(ev.DatumList) != 2 {
		return d, NewDecodeEventError(DepositedTag, 2, len(ev.DatumList))
	}
	oldDatum, err := ev.DatumList[0].Decode()
	if err != nil {
		return d, err
	}
	newDatum, err := ev.DatumList[1].Decode()
	if err != nil {
		return d, err
	}
	return Deposited{
		ChannelID: id,
		OldDatum:  oldDatum,
		NewDatum:  newDatum,
	}, nil
}

func (c Created) ID() channel.ID {
	return c.ChannelID
}

func (c Created) Timeout() channel.Timeout {
	return nil
}

func (c Created) Version() uint64 {
	return 0
}

func (c Created) ToPerunEvent() channel.AdjudicatorEvent {
	return nil
}

func (c Created) FromEvent(id ID, ev wire.Event) (Created, error) {
	if len(ev.DatumList) != 1 {
		return c, NewDecodeEventError(CreatedTag, 1, len(ev.DatumList))
	}
	datum, err := ev.DatumList[0].Decode()
	if err != nil {
		return c, err
	}
	return Created{
		ChannelID: id,
		NewDatum:  datum,
	}, nil
}
