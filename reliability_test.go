// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"testing"
)

func TestBitIndex(t *testing.T) {
	tests := [][3]uint32{
		{99, 100, 0},
		{0, 1, 0},
		{MaxSequence, 0, 0},
		{MaxSequence, 1, 1},
		{MaxSequence - 1, 1, 2},
		{MaxSequence - 1, 2, 3},
	}

	for _, bt := range tests {
		bit := bitIndex(bt[0], bt[1])

		if bit != bt[2] {
			t.Fatalf("Ack %x, Vector %x: Expected %x, got %x",
				bt[0], bt[1], bt[2], bit)
		}
	}
}

func TestAckVector(t *testing.T) {
	r := NewReliability()

	for i := uint32(0); i < 32; i++ {
		r.recvQueue.Insert(packetData{sequence: uint32(i)})
	}

	tests := [][2]uint32{
		{32, 0xffffffff},
		{31, 0x7fffffff},
		{33, 0xfffffffe},
		{16, 0x0000ffff},
		{48, 0xffff0000},
	}

	for _, bt := range tests {
		r.RemoteSequence = bt[0]
		vector := r.AckVector()

		if vector != bt[1] {
			t.Errorf("Ack %d. Want 0x%08x, Got 0x%08x", bt[0], bt[1], vector)
		}
	}
}

func TestAckVectorWrapped(t *testing.T) {
	r := NewReliability()

	r.recvQueue.Insert(packetData{sequence: MaxSequence - 1})
	r.recvQueue.Insert(packetData{sequence: MaxSequence})
	r.recvQueue.Insert(packetData{sequence: 0})

	tests := [][2]uint32{
		{0, 0x3},
		{MaxSequence, 0x1},
		{1, 0x7},
		{MaxSequence - 1, 0x0},
		{MaxSequence - 2, 0x0},
		{16, 0x00038000},
		{32, 0x80000000},
		{33, 0x0},
	}

	for _, bt := range tests {
		r.RemoteSequence = bt[0]
		vector := r.AckVector()

		if vector != bt[1] {
			t.Errorf("Ack %d. Want 0x%08x, Got 0x%08x", bt[0], bt[1], vector)
		}
	}
}
