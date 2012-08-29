// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

type packetData struct {
	sequence uint32
	size     uint32
	time     float32
}

// packetQueue holds a list of packets, sorted by sequence number.
type packetQueue []packetData

// Clear empties the queue.
func (q *packetQueue) Clear() { *q = (*q)[:0] }

func (q packetQueue) Exists(seq uint32) bool {
	for _, p := range q {
		if p.sequence == seq {
			return true
		}
	}

	return false
}

// RemoveAt removes the entry at the given index.
func (q *packetQueue) RemoveAt(i int) {
	if i < 0 || i >= len(*q) {
		return
	}

	tq := *q
	copy(tq[i:], tq[i+1:])
	tq = tq[:len(tq)-1]
	*q = tq
}

// Insert inserts the given packet into the queue.
// This guarantees the packets remain sorted by sequence number.
func (q *packetQueue) Insert(p packetData) {
	tq := *q

	if len(tq) == 0 {
		*q = append(tq, p)
		return
	}

	pseq := p.sequence

	switch {
	case isMoreRecent(tq[0].sequence, pseq):
		tq = append(tq, p)
		sz := len(tq) - 1
		tq[0], tq[sz] = tq[sz], tq[0]

	case isMoreRecent(pseq, tq[len(tq)-1].sequence):
		tq = append(tq, p)

	default:
		for i := range tq {
			if !isMoreRecent(tq[i].sequence, pseq) {
				continue
			}

			tmp := make(packetQueue, len(tq)+1)
			copy(tmp, tq[:i])
			copy(tmp[i+1:], tq[i:])
			tmp[i] = p
			tq = tmp
			break
		}
	}

	*q = tq
}
