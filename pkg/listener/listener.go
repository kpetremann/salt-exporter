package listener

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/kpetremann/salt-exporter/pkg/event"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
)

type eventParser interface {
	Parse(message map[string]any) (event.SaltEvent, error)
}

const DefaultIPCFilepath = "/var/run/salt/master/master_event_pub.ipc"

// EventListener listens to the salt-master event bus and sends events to the event channel.
type EventListener struct {
	// ctx specificies the context used mainly for cancellation
	ctx context.Context

	// eventChan is the channel to send events to
	eventChan chan event.SaltEvent

	// iPCFilepath is filepath to the salt-master event bus
	iPCFilepath string

	// saltEventBus keeps the connection to the salt-master event bus
	saltEventBus net.Conn

	eventParser eventParser
}

// Open opens the salt-master event bus.
func (e *EventListener) Open() {
	log.Info().Str("file", e.iPCFilepath).Msg("connecting to salt-master event bus")
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
			return
		}
	}
}

// Close closes the salt-master event bus.
func (e *EventListener) Close() error {
	log.Info().Msg("disconnecting from salt-master event bus")
	if e.saltEventBus != nil {
		return e.saltEventBus.Close()
	} else {
		return errors.New("trying to close already closed bus")
	}
}

// Reconnect reconnects to the salt-master event bus.
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
func NewEventListener(ctx context.Context, eventParser eventParser, eventChan chan event.SaltEvent) *EventListener {
	e := EventListener{
		ctx:         ctx,
		eventChan:   eventChan,
		eventParser: eventParser,
		iPCFilepath: DefaultIPCFilepath,
	}
	return &e
}

// SetIPCFilepath sets the filepath to the salt-master event bus
//
// The IPC file must be readable by the user running the exporter.
//
// Default: /var/run/salt/master/master_event_pub.ipc.
func (e *EventListener) SetIPCFilepath(filepath string) {
	e.iPCFilepath = filepath
}

// readFramedEvent reads one length-prefixed msgpack frame from the salt
// event bus and unpacks it into a map with "head" and "body" keys, per
// salt's frame_msg_ipc wire format: 4-byte big-endian length + msgpack payload.
func (e *EventListener) readFramedEvent() (map[string]interface{}, error) {
	lenBuf := make([]byte, 4)
	n, err := io.ReadFull(e.saltEventBus, lenBuf)
	if err != nil {
		return nil, fmt.Errorf("read length prefix: %w", err)
	}

	msgLen := binary.BigEndian.Uint32(lenBuf)
	if msgLen == 0 {
		log.Info().Msg("frame length is zero, skipping read of payload")
		return nil, fmt.Errorf("frame length is zero")
	}

	payload := make([]byte, msgLen)
	n, err = io.ReadFull(e.saltEventBus, payload)
	if err != nil {
		return nil, fmt.Errorf("read payload: %w", err)
	}

	var framed map[string]interface{}
	if err := msgpack.Unmarshal(payload, &framed); err != nil {
		return nil, fmt.Errorf("unmarshal payload: %w", err)
	}

	return framed, nil
}

func (e *EventListener) ListenEvents() {
	e.Open()

	for {
		select {
		case <-e.ctx.Done():
			log.Info().Msg("stop listening events")
			e.Close()
			return
		default:
			framed, err := e.readFramedEvent()
			if err != nil {
				log.Error().Str("error", err.Error()).Msg("unable to read event")
				log.Error().Msg("event bus may be closed, trying to reconnect")

				e.Reconnect()

				continue
			}

			evt, err := e.eventParser.Parse(framed)
			if err != nil {
				log.Debug().
					Str("error", err.Error()).
					Interface("framed", framed).
					Msg("event parser failed")
				continue
			}

			e.eventChan <- evt
		}
	}
}
