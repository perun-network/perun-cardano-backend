package channel

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	gpchannel "perun.network/go-perun/channel"
	"perun.network/perun-cardano-backend/channel/types"
	"perun.network/perun-cardano-backend/wire"
)

type AdjudicatorSub struct {
	eventQueue chan wire.Event
	connection *websocket.Conn
	lastError  chan error
	ChannelID  types.ID
}

func NewAdjudicatorSub(pabUrl url.URL, id types.ID) (*AdjudicatorSub, error) {
	conn, _, err := websocket.DefaultDialer.Dial(pabUrl.String(), nil)
	if err != nil {
		return nil, errors.New("unable to establish connection to PAB")
	}
	a := &AdjudicatorSub{
		eventQueue: make(chan wire.Event),
		connection: conn,
		lastError:  make(chan error),
		ChannelID:  id,
	}
	go receiveEvents(a)
	return a, nil
}

func receiveEvents(a *AdjudicatorSub) {
	var event wire.Event
	for {
		err := a.connection.ReadJSON(&event)
		if err != nil {
			a.lastError <- err
			close(a.eventQueue)
			close(a.lastError)
			_ = a.connection.Close()
			return
		}
		a.eventQueue <- event
	}
}

func (a AdjudicatorSub) Next() gpchannel.AdjudicatorEvent {
	event, ok := <-a.eventQueue
	if !ok {
		return nil
	}
	adjEvent, err := decodeEvent(event, a.ChannelID)
	if err != nil {
		a.lastError <- err
		return nil
	}
	return adjEvent
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

func decodeEvent(event wire.Event, id types.ID) (gpchannel.AdjudicatorEvent, error) {
	switch event.Tag {
	case types.DepositedTag:
		if len(event.ChannelDatum) != 1 {
			return nil, fmt.Errorf("invalid amout of ChannelDatums received in event. Amout: %d", len(event.ChannelDatum))
		}
		datum, err := event.ChannelDatum[0].Decode()
		if err != nil {
			return nil, err
		}
		return types.Deposited{
			ChannelID:    id,
			ChannelDatum: datum,
		}, nil
	case types.DisputedTag:
		if len(event.ChannelDatum) != 1 {
			return nil, fmt.Errorf("invalid amout of ChannelDatums received in event. Amout: %d", len(event.ChannelDatum))
		}
		datum, err := event.ChannelDatum[0].Decode()
		if err != nil {
			return nil, err
		}
		return types.Disputed{
			ChannelID:    id,
			ChannelDatum: datum,
		}, nil
	case types.ProgressedTag:
		// TODO: Figure out what Progressed Events do?
		if len(event.ChannelDatum) != 1 {
			return nil, fmt.Errorf("invalid amout of ChannelDatums received in event. Amout: %d", len(event.ChannelDatum))
		}
		datum, err := event.ChannelDatum[0].Decode()
		if err != nil {
			return nil, err
		}
		return types.Progressed{
			ChannelID:    id,
			ChannelDatum: datum,
		}, nil
	case types.ConcludedTag:
		return types.Concluded{ChannelID: id}, nil
	case types.WithdrawnTag:
		return types.Withdrawn{ChannelID: id}, nil
	default:
		return nil, fmt.Errorf("invalid event tag: %s", event.Tag)
	}
}
