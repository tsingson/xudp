// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package reliability

type packetData struct {
	sequence uint32
	size     uint32
	time     float32
}

// packetQueue holds a list of packets, sorted by sequence number.
type packetQueue []*packetData

// Exists returns true if the given sequence number is present
// in one of the list's elements.
func (q packetQueue) Exists(seq uint32) bool {
	for _, pd := range q {
		if pd.sequence == seq {
			return true
		}
	}

	return false
}

// RemoveAt removes the entry at the given index.
func (q *packetQueue) RemoveAt(i int) {
	tq := *q

	if i < 0 || i >= len(tq) {
		return
	}

	copy(tq[i:], tq[i+1:])
	tq = tq[:len(tq)-1]
	*q = tq
}

// Insert inserts p ensuring the list remains sorted by sequence number.
func (q *packetQueue) Insert(p *packetData) {
	if len(*q) == 0 {
		*q = append(*q, p)
		return
	}

	seq := p.sequence
	tq := *q

	if seq == tq[0].sequence || seq == tq[len(tq)-1].sequence {
		return // Duplicate -- Ignore it.
	}

	// Oldest sequence in list?
	if isMoreRecent(tq[0].sequence, seq) {
		tq = append(tq, p)
		copy(tq[1:], tq)
		tq[0] = p
		*q = tq
		return
	}

	// Newest sequence in list?
	if isMoreRecent(seq, tq[len(tq)-1].sequence) {
		*q = append(tq, p)
		return
	}

	// Somewhere in between -- Find out where.
	for i := 1; i < len(tq)-1; i++ {
		if tq[i].sequence == seq {
			return // Duplicate -- ignore it.
		}

		if isMoreRecent(tq[i].sequence, seq) {
			tmp := make(packetQueue, len(tq)+1)
			copy(tmp, tq[:i])
			tmp[i] = p
			copy(tmp[i+1:], tq[i:])
			*q = tmp
			return
		}
	}
}
