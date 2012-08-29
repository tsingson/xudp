// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"testing"
)

const (
	PacketCount = 100
	PacketSize  = 100
	DeltaTime   = 0.1
)

func TestReliableSend(t *testing.T) {
	r := NewReliability()

	for i := 0; i < PacketCount; i++ {
		r.PacketSent(PacketSize)
	}

	validate(t, r, PacketCount, 0, PacketCount, 0, 0, 0)
}

func TestReliableRecv(t *testing.T) {
	r := NewReliability()

	for i := 0; i < PacketCount; i++ {
		r.PacketRecv(uint32(i), PacketSize)
	}

	validate(t, r, 0, PacketCount-1, 0, PacketCount, 0, 0)
}

func TestReliableQueueTime(t *testing.T) {
	r := NewReliability()

	for i := 0; i < PacketCount; i++ {
		r.PacketSent(PacketSize)
	}

	for i := 0; i < PacketCount; i++ {
		r.AdvanceQueueTime(DeltaTime)
	}

	const epsilon = 0.001

	time := float32(DeltaTime * PacketCount)
	min := time - epsilon
	max := time + epsilon

	for _, pd := range r.sentQueue {
		if pd.time < min || pd.time > max {
			t.Errorf("sentQueue time mismatch: Want %f, have %f",
				time, pd.time)
		}
	}

	for _, pd := range r.recvQueue {
		if pd.time < min || pd.time > max {
			t.Errorf("recvQueue time mismatch: Want %f, have %f",
				time, pd.time)
		}
	}

	for _, pd := range r.pendingAckQueue {
		if pd.time < min || pd.time > max {
			t.Errorf("pendingAckQueue time mismatch: Want %f, have %f",
				time, pd.time)
		}
	}

	for _, pd := range r.ackedQueue {
		if pd.time < min || pd.time > max {
			t.Errorf("ackedQueue time mismatch: Want %f, have %f",
				time, pd.time)
		}
	}
}

func TestReliableUpdateQueue(t *testing.T) {
	r := NewReliability()

	for i := 0; i < PacketCount; i++ {
		r.PacketSent(PacketSize)
		r.AdvanceQueueTime(DeltaTime)
	}

	r.UpdateQueues()

	if sz := len(r.sentQueue); sz != 10 {
		t.Errorf("sentQueue mismatch: Want 10, have %d", sz)
	}

	if sz := len(r.recvQueue); sz != 0 {
		t.Errorf("recvQueue mismatch: Want 0, have %d", sz)
	}

	if sz := len(r.pendingAckQueue); sz != 10 {
		t.Errorf("pendingAckQueue mismatch: Want 10, have %d", sz)
	}

	if sz := len(r.ackedQueue); sz != 0 {
		t.Errorf("ackedQueue mismatch: Want 0, have %d", sz)
	}
}

func validate(t *testing.T, r *Reliability, ls, rs uint32, sp, rp, ap, lp uint64) {
	if r.LocalSequence != ls {
		t.Errorf("LocalSequence mismatch: Want %d, have %d", ls, r.LocalSequence)
	}

	if r.RemoteSequence != rs {
		t.Errorf("LocalSequence mismatch: Want %d, have %d", rs, r.RemoteSequence)
	}

	if r.SentPackets != sp {
		t.Errorf("SentPackets mismatch: Want %d, have %d", sp, r.SentPackets)
	}

	if r.RecvPackets != rp {
		t.Errorf("RecvPackets mismatch: Want %d, have %d", rp, r.RecvPackets)
	}

	if r.AckedPackets != ap {
		t.Errorf("AckedPackets mismatch: Want %d, have %d", ap, r.AckedPackets)
	}

	if r.LostPackets != lp {
		t.Errorf("LostPackets mismatch: Want %d, have %d", lp, r.LostPackets)
	}
}
