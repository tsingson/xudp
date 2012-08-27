// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

// PacketQueue holds a list of packets, sorted by sequence number.
type PacketQueue []Packet

func (q PacketQueue) Exists(seq uint32) bool {
	for _, p := range q {
		if p.Sequence() == seq {
			return true
		}
	}

	return false
}

// Insert inserts the given packet into the queue.
// This guarantees the packets remain sorted by sequence number.
func (q *PacketQueue) Insert(p Packet, maxSequence uint32) {
	if len(*q) == 0 {
		*q = append(*q, p)
		return
	}

	tq := *q
	pseq := p.Sequence()

	switch {
	case isMoreRecent(tq[0].Sequence(), pseq, maxSequence):
		tq = append(tq, p)
		sz := len(tq) - 1
		tq[0], tq[sz] = tq[sz], tq[0]

	case isMoreRecent(pseq, tq[len(tq)-1].Sequence(), maxSequence):
		tq = append(tq, p)

	default:
		for i := range tq {
			if !isMoreRecent(tq[i].Sequence(), pseq, maxSequence) {
				continue
			}

			tmp := make(PacketQueue, len(tq)+1)
			copy(tmp, tq[:i])
			copy(tmp[i+1:], tq[i:])
			tmp[i] = p
			tq = tmp
			break
		}
	}

	*q = tq
}

func (q PacketQueue) VerifySorted(max uint32) bool {
	for i := range q {
		if q[i].Sequence() > max {
			return false
		}

	}

	return true
}

func isMoreRecent(a, b, max uint32) bool {
	return (a > b) && (a-b <= max) || (b > a) && (b-a > max)
}
