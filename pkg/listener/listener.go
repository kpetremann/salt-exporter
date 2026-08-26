package listener

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
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

	// reader buffers the event bus connection, used both to detect the
	// wire framing and, in legacy mode, as the decoder's source
	reader *bufio.Reader

	// decoder is msgpack decoder for parsing the event bus messages in
	// legacy (unframed) mode
	decoder *msgpack.Decoder

	// tells whether the connected master uses the 4-byte length prefix
	hasLenPrefix bool

	eventParser eventParser

	onParseError func(message map[string]any, err error)
}

// Open opens the salt-master event bus.
func (e *EventListener) Open() {
	log.Info().Str("file", e.iPCFilepath).Msg("connecting to salt-master event bus")

	for {
		select {
		case <-e.ctx.Done():
			return
		default:
		}

		conn, err := net.Dial("unix", e.iPCFilepath)
		if err != nil {
			log.Error().Msg("failed to connect to event bus, retrying in 5 seconds")
			time.Sleep(time.Second * 5)

			continue
		}

		reader := bufio.NewReader(conn)

		hasPrefix, err := hasLenPrefix(reader)
		if err != nil {
			log.Error().Str("error", err.Error()).Msg("failed to detect event bus message format, retrying in 5 seconds")
			conn.Close()
			time.Sleep(time.Second * 5)
			continue
		}
		log.Info().Bool("hasLenPrefix", hasPrefix).Msg("successfully connected to event bus")

		e.saltEventBus = conn
		e.reader = reader
		e.hasLenPrefix = hasPrefix
		e.decoder = msgpack.NewDecoder(reader)

		return
	}
}

// hasLenPrefix checks the first byte to determine whether the message is using legacy or new format (4-byte length prefix).
func hasLenPrefix(reader *bufio.Reader) (bool, error) {
	b, err := reader.Peek(1)
	if err != nil {
		return false, err
	}

	// the first byte of msgpqck is between 0x80 and 0x8F
	// (see https://github.com/msgpack/msgpack/blob/master/spec.md)
	noPrefix := b[0] >= 0x80 && b[0] <= 0x8f
	return !noPrefix, nil
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

// readMessage decodes a Salt message. It handles both legacy format and new format.
func (e *EventListener) readMessage() (map[string]any, error) {
	if !e.hasLenPrefix {
		return e.decoder.DecodeMap()
	}

	var length [4]byte
	if _, err := io.ReadFull(e.reader, length[:]); err != nil {
		return nil, err
	}

	payload := make([]byte, binary.BigEndian.Uint32(length[:]))
	if _, err := io.ReadFull(e.reader, payload); err != nil {
		return nil, err
	}

	var message map[string]any
	if err := msgpack.Unmarshal(payload, &message); err != nil {
		return nil, err
	}

	return message, nil
}

// ListenEvents listens to the salt-master event bus and sends events to the event channel.
func (e *EventListener) ListenEvents() {
	e.Open()

	for {
		select {
		case <-e.ctx.Done():
			log.Info().Msg("stop listening events")
			e.Close()
			return
		default:
			message, err := e.readMessage()
			if err != nil {
				log.Error().Str("error", err.Error()).Msg("unable to read event")
				log.Error().Msg("event bus may be closed, trying to reconnect")

				e.Reconnect()

				continue
			}
			if event, err := e.eventParser.Parse(message); err == nil {
				e.eventChan <- event
			} else if e.onParseError != nil {
				e.onParseError(message, err)
			}
		}
	}
}
