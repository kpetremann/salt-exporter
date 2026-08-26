package listener_test

import (
	"context"
	"encoding/binary"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/kpetremann/salt-exporter/pkg/event"
	"github.com/kpetremann/salt-exporter/pkg/listener"
	"github.com/kpetremann/salt-exporter/pkg/parser"
	"github.com/vmihailenco/msgpack/v5"
)

// buildFrame encodes one salt-master event bus message: a msgpack map with
// "head" and "body" keys, "body" being the raw tag followed by the
// msgpack-encoded event data, exactly as
// salt.transport.frame.frame_msg[_ipc] produces it.
func buildFrame(t *testing.T, tag string, data any) []byte {
	t.Helper()

	dataBytes, err := msgpack.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal event data: %v", err)
	}

	body := append([]byte(tag+"\n\n"), dataBytes...)

	frame, err := msgpack.Marshal(map[string]any{"head": map[string]any{}, "body": body})
	if err != nil {
		t.Fatalf("failed to marshal frame: %v", err)
	}

	return frame
}

// withLengthPrefix prefixes a frame with the 4-byte big-endian length used
// by newer Salt masters (3006.x since d4e2e075aa3, and all 3008.x) to frame
// IPC messages.
func withLengthPrefix(frame []byte) []byte {
	buf := make([]byte, 4+len(frame))
	binary.BigEndian.PutUint32(buf, uint32(len(frame))) //nolint:gosec // test frames are always small
	copy(buf[4:], frame)

	return buf
}

// serveFakeEventBus starts a fake salt-master event bus at path and writes
// raw to the first accepted connection.
func serveFakeEventBus(t *testing.T, path string, raw []byte) {
	t.Helper()

	ln, err := net.Listen("unix", path)
	if err != nil {
		t.Fatalf("failed to listen on %s: %v", path, err)
	}
	t.Cleanup(func() { ln.Close() })

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		if _, err := conn.Write(raw); err != nil {
			return
		}

		// keep the connection open until the test is done reading, so the
		// listener doesn't see EOF and reconnect mid-test
		<-time.After(2 * time.Second)
	}()
}

// TestEventListener_WireFormats ensures events are read correctly whether
// the salt-master event bus uses the legacy, unframed msgpack stream or the
// newer 4-byte length-prefixed framing.
func TestEventListener_WireFormats(t *testing.T) {
	tag1 := "salt/job/20231009113602182345/new"
	data1 := map[string]any{
		"fun":     "test.ping",
		"jid":     "20231009113602182345",
		"minions": []string{"host1.example.com"},
	}

	tag2 := "salt/job/20231009113602182345/ret/host1.example.com"
	data2 := map[string]any{
		"fun":     "test.ping",
		"jid":     "20231009113602182345",
		"id":      "host1.example.com",
		"success": true,
		"return":  true,
	}

	frames := [][]byte{
		buildFrame(t, tag1, data1),
		buildFrame(t, tag2, data2),
	}

	tests := []struct {
		name string
		raw  func() []byte
	}{
		{
			name: "legacy unframed msgpack stream",
			raw: func() []byte {
				raw := make([]byte, 0, len(frames[0])+len(frames[1]))
				for _, frame := range frames {
					raw = append(raw, frame...)
				}
				return raw
			},
		},
		{
			name: "4-byte length-prefixed framing",
			raw: func() []byte {
				raw := make([]byte, 0, len(frames[0])+len(frames[1])+2*4)
				for _, frame := range frames {
					raw = append(raw, withLengthPrefix(frame)...)
				}
				return raw
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			socketPath := filepath.Join(t.TempDir(), "master_event_pub.ipc")
			serveFakeEventBus(t, socketPath, test.raw())

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			eventChan := make(chan event.SaltEvent, 2)
			eventListener := listener.NewEventListener(ctx, parser.NewEventParser(false), eventChan)
			eventListener.SetIPCFilepath(socketPath)

			go eventListener.ListenEvents()

			var got []event.SaltEvent
			for i := 0; i < 2; i++ {
				select {
				case ev := <-eventChan:
					got = append(got, ev)
				case <-time.After(3 * time.Second):
					t.Fatalf("timed out waiting for event %d", i+1)
				}
			}

			if got[0].Tag != tag1 {
				t.Errorf("event 1: got tag %q, want %q", got[0].Tag, tag1)
			}
			if got[0].Data.Fun != "test.ping" {
				t.Errorf("event 1: got fun %q, want %q", got[0].Data.Fun, "test.ping")
			}
			if got[0].TargetNumber != 1 {
				t.Errorf("event 1: got target number %d, want 1", got[0].TargetNumber)
			}

			if got[1].Tag != tag2 {
				t.Errorf("event 2: got tag %q, want %q", got[1].Tag, tag2)
			}
			if got[1].Data.ID != "host1.example.com" {
				t.Errorf("event 2: got id %q, want %q", got[1].Data.ID, "host1.example.com")
			}
			if got[1].StateModuleSuccess != nil {
				t.Errorf("event 2: got state module success %v, want nil (not a state return)", got[1].StateModuleSuccess)
			}
		})
	}
}
