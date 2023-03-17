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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	gpchannel "perun.network/go-perun/channel"
	"perun.network/perun-cardano-backend/channel/types"
	"perun.network/perun-cardano-backend/wire"
)

// AdjudicatorSub is a subscription to the Adjudicator events.
// Instances should only be created using PAB.NewSubscription.
type AdjudicatorSub struct {
	eventQueue chan wire.Event
	connection *websocket.Conn
	lastError  chan error
	ChannelID  types.ID
	IsPerunSub bool
}

func newAdjudicatorSub(contractUrl *url.URL, id types.ID, isPerunSub bool) (*AdjudicatorSub, error) {
	conn, _, err := websocket.DefaultDialer.Dial(contractUrl.String(), nil)
	if err != nil {
		return nil, errors.New("unable to establish connection to PAB")
	}
	a := &AdjudicatorSub{
		eventQueue: make(chan wire.Event),
		connection: conn,
		lastError:  make(chan error),
		ChannelID:  id,
		IsPerunSub: isPerunSub,
	}
	go receiveEvents(a)
	return a, nil
}

func receiveEvents(a *AdjudicatorSub) {
	var message wire.SubscriptionMessage
	for {
		err := a.connection.ReadJSON(&message)
		if err != nil {
			a.lastError <- err
			close(a.eventQueue)
			close(a.lastError)
			_ = a.connection.Close()
			return
		}
		if message.Tag != wire.EventMessageTag {
			continue
		}
		var events []wire.Event
		err = json.Unmarshal(message.Contents, &events)
		if err != nil {
			a.lastError <- fmt.Errorf("malformed event message: %w", err)
			close(a.eventQueue)
			close(a.lastError)
			_ = a.connection.Close()
			return
		}
		for _, e := range events {
			a.eventQueue <- e
		}
	}
}

func (a AdjudicatorSub) Next() gpchannel.AdjudicatorEvent {
	for {
		event, ok := <-a.eventQueue

		if !ok {
			return nil
		}
		adjEvent, err := decodeEvent(event, a.ChannelID)
		if err != nil {
			a.lastError <- err
			return nil
		}
		if !a.IsPerunSub {
			return adjEvent
		}
		perunEvent := adjEvent.ToPerunEvent()
		if perunEvent != nil {
			return perunEvent
		}
	}
}

func (a AdjudicatorSub) Err() error {
	select {
	case err := <-a.lastError:
		return err
	default:
		return nil
	}
}

func (a AdjudicatorSub) Close() error {
	return a.connection.Close()
}

func decodeEvent(event wire.Event, id types.ID) (types.InternalEvent, error) {
	const errorFormat = "invalid amount of ChannelDatums received in %s event. Expected: %d, Actual: %d"
	switch event.Tag {
	case types.CreatedTag:
		if len(event.DatumList) != 1 {
			return nil, fmt.Errorf(errorFormat, types.CreatedTag, 1, len(event.DatumList))
		}
		datum, err := event.DatumList[0].Decode()
		if err != nil {
			return nil, err
		}
		return types.Created{
			ChannelID: id,
			NewDatum:  datum,
		}, nil
	case types.DepositedTag:
		if len(event.DatumList) != 2 {
			return nil, fmt.Errorf(errorFormat, types.DepositedTag, 2, len(event.DatumList))
		}
		oldDatum, err := event.DatumList[0].Decode()
		if err != nil {
			return nil, err
		}
		newDatum, err := event.DatumList[1].Decode()
		if err != nil {
			return nil, err
		}
		return types.Deposited{
			ChannelID: id,
			OldDatum:  oldDatum,
			NewDatum:  newDatum,
		}, nil
	case types.DisputedTag:
		if len(event.DatumList) != 2 {
			return nil, fmt.Errorf(errorFormat, types.DisputedTag, 2, len(event.DatumList))
		}
		oldDatum, err := event.DatumList[0].Decode()
		if err != nil {
			return nil, err
		}
		newDatum, err := event.DatumList[1].Decode()
		if err != nil {
			return nil, err
		}
		return types.Disputed{
			ChannelID: id,
			OldDatum:  oldDatum,
			NewDatum:  newDatum,
		}, nil
	case types.ConcludedTag:
		if len(event.DatumList) != 1 {
			return nil, fmt.Errorf(errorFormat, types.ConcludedTag, 1, len(event.DatumList))
		}
		datum, err := event.DatumList[0].Decode()
		if err != nil {
			return nil, err
		}
		return types.Concluded{
			ChannelID: id,
			OldDatum:  datum,
		}, nil
	default:
		return nil, fmt.Errorf("invalid event tag: %s", event.Tag)
	}
}
