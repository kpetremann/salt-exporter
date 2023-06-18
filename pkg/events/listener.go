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

// EventListener listens to the salt-master event bus and sends events to the event channel
type EventListener struct {
	// ctx specificies the context used mainly for cancellation
	ctx context.Context

	// eventChan is the channel to send events to
	eventChan chan SaltEvent

	// iPCFilepath is filepath to the salt-master event bus
	iPCFilepath string

	// saltEventBus keeps the connection to the salt-master event bus
	saltEventBus net.Conn

	// decoder is msgpack decoder for parsing the event bus messages
	decoder *msgpack.Decoder
}

// Open opens the salt-master event bus
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

// Close closes the salt-master event bus
func (e *EventListener) Close() error {
	log.Info().Msg("disconnecting from salt-master event bus")
	if e.saltEventBus != nil {
		return e.saltEventBus.Close()
	} else {
		return errors.New("trying to close already closed bus")
	}
}

// Reconnect reconnects to the salt-master event bus
func (e *EventListener) Reconnect() {
	select {
	case <-e.ctx.Done():
		return
	default:
		e.Close()
		e.Open()
	}
}

// NewEventListener creates a new EventListener
//
// The events will be sent to eventChan.
func NewEventListener(ctx context.Context, eventChan chan SaltEvent) *EventListener {
	e := EventListener{ctx: ctx, eventChan: eventChan, iPCFilepath: defaultIPCFilepath}
	return &e
}

// SetIPCFilepath sets the filepath to the salt-master event bus
//
// The IPC file must be readable by the user running the exporter.
//
// Default: /var/run/salt/master/master_event_pub.ipc
func (e *EventListener) SetIPCFilepath(filepath string) {
	e.iPCFilepath = filepath
}

// ListenEvents listens to the salt-master event bus and sends events to the event channel
//
// if keepRawBody is true, the raw event body will be kept in the event struct.
// It can be useful for debugging or post-processing.
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
			if event, err := ParseEvent(message, keepRawBody); err == nil {
				e.eventChan <- event
			}
		}
	}
}
