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
	eventQueue chan gpchannel.AdjudicatorEvent
	connection *websocket.Conn
	lastError  chan error
	ChannelID  types.ID
	close      chan struct{}
	// IsPerunSub specifies whether a subscription yields Perun events or Internal events.
	IsPerunSub        bool
	receivedNilOnNext bool
	receivedError     error
}

func newAdjudicatorSub(contractUrl *url.URL, id types.ID, isPerunSub bool) (*AdjudicatorSub, error) {
	conn, _, err := websocket.DefaultDialer.Dial(contractUrl.String(), nil)
	if err != nil {
		return nil, errors.New("unable to establish connection to PAB")
	}
	a := &AdjudicatorSub{
		eventQueue: make(chan gpchannel.AdjudicatorEvent),
		connection: conn,
		lastError:  make(chan error, 1),
		ChannelID:  id,
		close:      make(chan struct{}),
		IsPerunSub: isPerunSub,
	}
	go receiveEvents(a)
	return a, nil
}

func receiveEvents(a *AdjudicatorSub) {
	closeGracefully := func(err error) {
		a.lastError <- err
		close(a.eventQueue)
		close(a.lastError)
		_ = a.connection.Close()
	}

	var message wire.SubscriptionMessage
	for {
		err := a.connection.ReadJSON(&message)
		if err != nil {
			closeGracefully(err)
			return
		}
		if message.Tag != wire.EventMessageTag {
			continue
		}
		var events []wire.Event
		err = json.Unmarshal(message.Contents, &events)
		if err != nil {
			closeGracefully(fmt.Errorf("malformed event message: %w", err))
			return
		}
		for _, e := range events {
			adjEvent, err := decodeEvent(e, a.ChannelID)
			if err != nil {
				closeGracefully(err)
				return
			}
			if !a.IsPerunSub {
				select {
				case a.eventQueue <- adjEvent:
				case <-a.close:
					closeGracefully(errors.New("subscription closed by user"))
					return
				}
				a.eventQueue <- adjEvent
			}
			perunEvent := adjEvent.ToPerunEvent()
			if perunEvent == nil {
				continue
			}
			select {
			case a.eventQueue <- perunEvent:
			case <-a.close:
				closeGracefully(errors.New("subscription closed by user"))
				return
			}
		}
	}
}

// Next returns the next AdjudicatorEvent. It blocks until the next event is available.
// If the subscription is closed, or there is an error, it returns nil.
// Once Next returns nil, a subsequent call to Err will return the error that caused the subscription to close and all
// subsequent calls to Next will also return nil.
// It is important to only interact with one Subscription from a single go-routine!
// Note: This may return either a types.InternalEvent or an AdjudicatorEvent depending on the type of the subscription.
func (a AdjudicatorSub) Next() gpchannel.AdjudicatorEvent {
	if a.receivedNilOnNext {
		return nil
	}
	ret := <-a.eventQueue
	if ret == nil {
		a.receivedNilOnNext = true
	}
	return ret
}

// Err returns the error after a call to Next returned nil, or nil if there is no error.
// Once Err returns a non-nil error, all subsequent calls to Err will return the same error.
// It is important to only interact with one Subscription from a single go-routine!
func (a AdjudicatorSub) Err() error {
	if a.receivedError != nil {
		return a.receivedError
	}
	select {
	case err := <-a.lastError:
		a.receivedError = err
		return err
	default:
		return nil
	}
}

// Close closes the subscription.
func (a AdjudicatorSub) Close() error {
	close(a.close)
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
