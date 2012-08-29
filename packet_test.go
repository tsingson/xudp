// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"testing"
)

func TestPacket(t *testing.T) {
	const (
		proto    = 0x12345678
		sequence = 0x2
		ack      = 0x1
		vector   = 0x0000ffff
	)

	payload := []byte{1, 2, 3, 4, 5}

	p := NewPacket(payload)
	p.SetHeader(proto, sequence, ack, vector)

	if p.Protocol() != proto {
		t.Fatalf("Protocol mismatch: Want %x, have %x",
			proto, p.Protocol())
	}

	if p.Sequence() != sequence {
		t.Fatalf("Sequence mismatch: Want %x, have %x",
			sequence, p.Sequence())
	}

	if p.Ack() != ack {
		t.Fatalf("Ack mismatch: Want %x, have %x",
			ack, p.Ack())
	}

	if p.AckVector() != vector {
		t.Fatalf("AckVector mismatch: Want %x, have %x",
			vector, p.AckVector())
	}

	if len(p.Payload()) != len(payload) {
		t.Fatalf("Payload size mismatch: Want %x, have %x",
			len(p.Payload()), len(payload))
	}

	for i, b := range p.Payload() {
		if b != payload[i] {
			t.Fatalf("Payload mismatch at %d: Want %x, have %x",
				i, payload[i], b)
		}
	}
}
