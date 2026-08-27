package listener

import (
	"bufio"
	"bytes"
	"testing"
)

// TestHasLenPrefix pins the exact byte range used to tell the legacy,
// unframed wire format (a bare msgpack map, always starting with a fixmap
// header byte in 0x80-0x8f) apart from the newer 4-byte length-prefixed
// framing (whose first byte is 0x00 for any realistic event size).
func TestHasLenPrefix(t *testing.T) {
	tests := []struct {
		name      string
		firstByte byte
		want      bool
	}{
		{name: "fixmap lower bound (0 pairs)", firstByte: 0x80, want: false},
		{name: "fixmap with 2 pairs (head+body)", firstByte: 0x82, want: false},
		{name: "fixmap upper bound (15 pairs)", firstByte: 0x8f, want: false},
		{name: "just below the fixmap range", firstByte: 0x7f, want: true},
		{name: "just above the fixmap range", firstByte: 0x90, want: true},
		{name: "length prefix high byte for typical event sizes", firstByte: 0x00, want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := bufio.NewReader(bytes.NewReader([]byte{test.firstByte, 0x00, 0x00, 0x00}))

			got, err := hasLenPrefix(reader)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != test.want {
				t.Errorf("hasLenPrefix(0x%02x) = %v, want %v", test.firstByte, got, test.want)
			}
		})
	}
}
