// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package reliability

import (
	"math/rand"
	"testing"
	"time"
)

func TestQueueInsertBack(t *testing.T) {
	var q packetQueue

	for i := uint32(0); i < 100; i++ {
		q.Insert(&packetData{sequence: i})

		if !isQueueSorted(q) {
			t.Fatalf("Queue sorting failure at sequence %d", i)
		}
	}
}

func TestQueueInsertFront(t *testing.T) {
	var q packetQueue

	for i := 100; i >= 1; i-- {
		q.Insert(&packetData{sequence: uint32(i)})

		if !isQueueSorted(q) {
			t.Fatalf("Queue sorting failure at sequence %d", i)
		}
	}
}

func TestQueueInsertRandom(t *testing.T) {
	var q packetQueue
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := uint32(100); i >= 1; i-- {
		q.Insert(&packetData{sequence: uint32(r.Int31n(int32(i)))})

		if !isQueueSorted(q) {
			t.Fatalf("Queue sorting failure at sequence %d", i)
		}
	}
}

func TestQueueInsertWrapped(t *testing.T) {
	var q packetQueue

	for i := uint32(MaxSequence - 5); i < MaxSequence; i++ {
		q.Insert(&packetData{sequence: i})

		if !isQueueSorted(q) {
			t.Fatalf("Queue sorting failure at sequence %d", i)
		}
	}

	for i := uint32(0); i <= 5; i++ {
		q.Insert(&packetData{sequence: i})

		if !isQueueSorted(q) {
			t.Fatalf("Queue sorting failure at sequence %d", i)
		}
	}
}

// isQueueSorted tests if the given packet queue is sorted
// by packet sequence number.
func isQueueSorted(q packetQueue) bool {
	if len(q) < 2 {
		return true
	}

	for i, pd := range q {
		if pd.sequence > MaxSequence {
			return false
		}

		if i == 0 {
			continue
		}

		if !isMoreRecent(pd.sequence, q[i-1].sequence) {
			return false
		}
	}

	return true
}
