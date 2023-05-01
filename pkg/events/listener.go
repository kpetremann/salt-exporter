package events

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
)

type EventListener struct {
	ctx          context.Context
	eventChan    chan SaltEvent
	saltEventBus net.Conn
	decoder      *msgpack.Decoder
}

func (e *EventListener) Open() net.Conn {
	log.Info().Msg("connecting to salt-master event bus")
	var eventBus net.Conn
	var err error

	for {
		select {
		case <-e.ctx.Done():
			return nil
		default:
		}

		eventBus, err = net.Dial("unix", "/var/run/salt/master/master_event_pub.ipc")
		if err != nil {
			log.Error().Msg("failed to connect to event bus, retrying in 5 seconds")
			time.Sleep(time.Second * 5)
		} else {
			log.Info().Msg("successfully connected to event bus")
			e.decoder = msgpack.NewDecoder(eventBus)
			return eventBus
		}
	}
}

func (e *EventListener) Close() error {
	log.Info().Msg("disconnecting from salt-master event bus")
	if e.saltEventBus != nil {
		return e.saltEventBus.Close()
	} else {
		return errors.New("trying to close already closed bus")
	}
}

func (e *EventListener) Reconnect() {
	select {
	case <-e.ctx.Done():
		return
	default:
		e.Close()
		e.Open()
	}
}

func NewEventListener(ctx context.Context, eventChan chan SaltEvent) *EventListener {
	e := EventListener{ctx: ctx, eventChan: eventChan}
	return &e
}

func (e *EventListener) ListenEvents(keepRawBody bool) {
	e.Open()

	for {
		select {
		case <-e.ctx.Done():
			log.Info().Msg("stop listening events")
			e.Close()
			return
		default:
			message, err := e.decoder.DecodeMap()
			if err != nil {
				log.Error().Str("error", err.Error()).Msg("unable to read event")
				log.Error().Msg("event bus may be closed, trying to reconnect")

				e.Reconnect()

				continue
			}
			ParseEvent(message, e.eventChan, keepRawBody)
		}
	}
}
