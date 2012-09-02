// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package reliability

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

	for i := 0; i < 32; i++ {
		r.recvQueue.Insert(&packetData{sequence: uint32(i)})
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
		vector := r.ackVector()

		if vector != bt[1] {
			t.Errorf("Ack %d. Want 0x%08x, Got 0x%08x", bt[0], bt[1], vector)
		}
	}
}

func TestAckVectorWrapped(t *testing.T) {
	r := NewReliability()

	r.recvQueue.Insert(&packetData{sequence: MaxSequence - 1})
	r.recvQueue.Insert(&packetData{sequence: MaxSequence})
	r.recvQueue.Insert(&packetData{sequence: 0})

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
		vector := r.ackVector()

		if vector != bt[1] {
			t.Errorf("Ack %d. Want 0x%08x, Got 0x%08x", bt[0], bt[1], vector)
		}
	}
}

func TestprocessAck1(t *testing.T) {
	r := NewReliability()

	for i := 0; i < 33; i++ {
		r.pendingAckQueue.Insert(&packetData{sequence: uint32(i)})
	}

	r.RTT = 0
	r.AckedPackets = 0
	r.processAck(32, 0xffffffff)

	if r.AckedPackets != 33 {
		t.Fatalf("AckedPackets mismatch: Want 33, got %d", r.AckedPackets)
	}

	if len(r.ackedQueue) != 33 {
		t.Fatalf("ackedQueue size mismatch: Want 33, got %d", len(r.ackedQueue))
	}

	if len(r.pendingAckQueue) != 0 {
		t.Fatalf("pendingAckQueue size mismatch: Want 0, got %d", len(r.pendingAckQueue))
	}

	if !isQueueSorted(r.ackedQueue) {
		t.Fatalf("ackedQueueis not sorted.")
	}

	for i := range r.ackedQueue {
		if r.ackedQueue[i].sequence != uint32(i) {
			t.Fatalf("ackedQueue %d mismatch.", i)
		}
	}
}

func TestprocessAck2(t *testing.T) {
	r := NewReliability()

	for i := 0; i < 33; i++ {
		r.pendingAckQueue.Insert(&packetData{sequence: uint32(i)})
	}

	r.RTT = 0
	r.AckedPackets = 0
	r.processAck(32, 0x0000ffff)

	n := 17
	if r.AckedPackets != uint32(n) {
		t.Fatalf("AckedPackets mismatch: Want %d, got %d", n, r.AckedPackets)
	}

	if len(r.ackedQueue) != n {
		t.Fatalf("ackedQueue size mismatch: Want %d, got %d", n, len(r.ackedQueue))
	}

	n = 33 - 17
	if len(r.pendingAckQueue) != n {
		t.Fatalf("pendingAckQueue size mismatch: Want %d, got %d", n, len(r.pendingAckQueue))
	}

	if !isQueueSorted(r.ackedQueue) {
		t.Fatalf("ackedQueueis not sorted.")
	}

	for i := range r.pendingAckQueue {
		if r.pendingAckQueue[i].sequence != uint32(i) {
			t.Fatalf("pendingAckQueue %d mismatch.", i)
		}
	}

	for i := range r.ackedQueue {
		if r.ackedQueue[i].sequence != uint32(i)+16 {
			t.Fatalf("ackedQueue %d mismatch.", i)
		}
	}
}

func TestprocessAck3(t *testing.T) {
	r := NewReliability()

	for i := 0; i < 32; i++ {
		r.pendingAckQueue.Insert(&packetData{sequence: uint32(i)})
	}

	r.RTT = 0
	r.AckedPackets = 0
	r.processAck(48, 0xffff0000)

	n := 16
	if r.AckedPackets != uint32(n) {
		t.Fatalf("AckedPackets mismatch: Want %d, got %d", n, r.AckedPackets)
	}

	if len(r.ackedQueue) != n {
		t.Fatalf("ackedQueue size mismatch: Want %d, got %d", n, len(r.ackedQueue))
	}

	if len(r.pendingAckQueue) != n {
		t.Fatalf("pendingAckQueue size mismatch: Want %d, got %d", n, len(r.pendingAckQueue))
	}

	if !isQueueSorted(r.ackedQueue) {
		t.Fatalf("ackedQueueis not sorted.")
	}

	for i := range r.pendingAckQueue {
		if r.pendingAckQueue[i].sequence != uint32(i) {
			t.Fatalf("pendingAckQueue %d mismatch.", i)
		}
	}

	for i := range r.ackedQueue {
		if r.ackedQueue[i].sequence != uint32(i)+16 {
			t.Fatalf("ackedQueue %d mismatch.", i)
		}
	}
}
