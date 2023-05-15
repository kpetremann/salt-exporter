package events

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
)

const defaultIPCFilepath = "/var/run/salt/master/master_event_pub.ipc"

type EventListener struct {
	ctx          context.Context
	eventChan    chan SaltEvent
	iPCFilepath  string
	saltEventBus net.Conn
	decoder      *msgpack.Decoder
}

func (e *EventListener) Open() {
	log.Info().Msg("connecting to salt-master event bus")
	var err error

	for {
		select {
		case <-e.ctx.Done():
			return
		default:
		}

		e.saltEventBus, err = net.Dial("unix", e.iPCFilepath)
		if err != nil {
			log.Error().Msg("failed to connect to event bus, retrying in 5 seconds")
			time.Sleep(time.Second * 5)
		} else {
			log.Info().Msg("successfully connected to event bus")
			e.decoder = msgpack.NewDecoder(e.saltEventBus)
			return
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
	e := EventListener{ctx: ctx, eventChan: eventChan, iPCFilepath: defaultIPCFilepath}
	return &e
}

// SetIPCFilepath sets the filepath to the salt-master event bus
//
// The IPC file must be readable by the user running the exporter
// Default: /var/run/salt/master/master_event_pub.ipc
func (e *EventListener) SetIPCFilepath(filepath string) {
	e.iPCFilepath = filepath
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
