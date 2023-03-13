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

import "perun.network/go-perun/channel"

const (
	CreatedTag   = "Created"
	DepositedTag = "Deposited"
	DisputedTag  = "Disputed"
	ConcludedTag = "Concluded"
)

// TODO: Figure out what to return on AdjudicatorEvent.Version(), if there is no concept of a state version for that event

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
		ChannelID ID
		OldDatum  ChannelDatum
		NewDatum  ChannelDatum
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

func (d Disputed) ID() channel.ID {
	return d.ChannelID
}

func (d Disputed) Timeout() channel.Timeout {
	return &channel.TimeTimeout{Time: d.NewDatum.Time.Add(d.NewDatum.ChannelParameters.Timeout)}
}

func (d Disputed) Version() uint64 {
	return d.NewDatum.ChannelState.Version
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

func (s Created) ID() channel.ID {
	return s.ChannelID
}

func (s Created) Timeout() channel.Timeout {
	return nil
}

func (s Created) Version() uint64 {
	return 0
}
