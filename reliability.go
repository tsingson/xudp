// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

// Maximum packet sequence value before overflow.	
const MaxSequence = 1<<32 - 1

// isMoreRecent checks if sequence a is newer than sequcen b,
// while taking integer overflow into account.
func isMoreRecent(a, b uint32) bool {
	return (a > b) && (a-b <= MaxSequence) || (b > a) && (b-a > MaxSequence)
}

// bitIndex finds the ack vector bit index for the given sequence number.
func bitIndex(sequence, ack uint32) uint32 {
	if sequence > ack {
		return ack + (MaxSequence - sequence)
	}
	return ack - 1 - sequence
}

// Reliability implements the algorithms needed to make a reliable connection
// reliable. This means it manages sent, received, pending ACKs and ACK'ed
// packet queues. In addtion, it tracks bandwidth use and rount trip timing.
type Reliability struct {
	sentQueue       packetQueue // Sent packets used to calculate sent bandwidth.
	pendingAckQueue packetQueue // Sent packets which have not been acked yet.
	recvQueue       packetQueue // Received packets used to determine acks to send.
	ackedQueue      packetQueue // ACK'ed packets.
	acks            []uint32    // ACK'ed packets from last set of packet receives.
	SentPackets     uint64      // Number of packets sent.
	RecvPackets     uint64      // Number of packets received.
	LostPackets     uint64      // Number of packets lost.
	AckedPackets    uint64      // Number of packets ACK'ed.
	LocalSequence   uint32      // Local sequence number for most recently sent packet.
	RemoteSequence  uint32      // Remote sequence number for most recently received packet.
	SentBandwidth   float32     // Approximate sent bandwidth over the last second.
	AckedBandwidth  float32     // Approximate ACK'ed bandwidth over the last second.
	RTT             float32     // Estimated round trip time.
	RTTMax          float32     // Maximum expected round trip time.
}

// NewReliability creates a new reliability instance.
func NewReliability() *Reliability {
	r := new(Reliability)
	r.Reset()
	return r
}

// PacketSent is called whenever a new packet is sent.
func (r *Reliability) PacketSent(size uint) {
	var pd packetData
	pd.sequence = r.LocalSequence
	pd.size = size

	r.sentQueue = append(r.sentQueue, pd)
	r.pendingAckQueue = append(r.pendingAckQueue, pd)
	r.SentPackets++
	r.LocalSequence++
}

// PacketRecv is called whenever a new packet is received.
func (r *Reliability) PacketRecv(sequence uint32, size uint) {
	r.RecvPackets++

	if r.recvQueue.Exists(sequence) {
		return
	}

	r.recvQueue = append(r.recvQueue, packetData{
		sequence: sequence,
		size:     size,
	})

	if isMoreRecent(sequence, r.RemoteSequence) {
		r.RemoteSequence = sequence
	}
}

// AckVector generates the ACK vector which should be included in an
// outgoing packet header.
func (r *Reliability) AckVector() uint32 {
	var vector, bit uint32

	ack := r.RemoteSequence

	for _, pd := range r.recvQueue {
		if pd.sequence == ack || isMoreRecent(pd.sequence, ack) {
			break
		}

		bit = bitIndex(pd.sequence, ack)

		if bit <= 31 {
			vector |= 1 << bit
		}
	}

	return vector
}

// ProcessAck handles a single incoming ACK with ACK vector.
func (r *Reliability) ProcessAck(ack, vector uint32) {
	if len(r.pendingAckQueue) == 0 {
		return
	}

	var pd packetData
	var acked bool
	var bit uint32

	for i := 0; i < len(r.pendingAckQueue); i++ {
		pd = r.pendingAckQueue[i]
		acked = false

		if pd.sequence == ack {
			acked = true

		} else if isMoreRecent(ack, pd.sequence) {
			bit = bitIndex(pd.sequence, ack)

			if bit <= 31 {
				acked = (vector>>bit)&1 != 0
			}
		}

		if !acked {
			continue
		}

		r.RTT += (pd.time - r.RTT) * 0.1
		r.ackedQueue.Insert(pd)
		r.acks = append(r.acks, pd.sequence)
		r.AckedPackets++
		r.pendingAckQueue.RemoveAt(i)
		i--
	}
}

// AdvanceQueueTime updates the timestamp for each queued packet.
func (r *Reliability) AdvanceQueueTime(delta float32) {
	var i int

	for i = range r.sentQueue {
		r.sentQueue[i].time += delta
	}

	for i = range r.recvQueue {
		r.recvQueue[i].time += delta
	}

	for i = range r.pendingAckQueue {
		r.pendingAckQueue[i].time += delta
	}

	for i = range r.ackedQueue {
		r.ackedQueue[i].time += delta
	}
}

// UpdateQueues updates all queues to discard packets that have
// exceeded their timeouts or are otherwise no longer necessary.
func (r *Reliability) UpdateQueues() {
	const epsilon = 0.001

	threshold := r.RTTMax + epsilon

	for len(r.sentQueue) > 0 && r.sentQueue[0].time > threshold {
		r.sentQueue = r.sentQueue[1:]
	}

	if sz := len(r.recvQueue); sz > 0 {
		var minSeq uint32
		lastSeq := r.recvQueue[sz-1].sequence

		if lastSeq >= 34 {
			minSeq = lastSeq - 34
		} else {
			minSeq = MaxSequence - (34 - lastSeq)
		}

		for len(r.recvQueue) > 0 && isMoreRecent(minSeq, r.recvQueue[0].sequence) {
			r.recvQueue = r.recvQueue[1:]
		}
	}
}

// Reset sets the Reliability system to its initial state.
func (r *Reliability) Reset() {
	r.sentQueue.Clear()
	r.recvQueue.Clear()
	r.pendingAckQueue.Clear()
	r.ackedQueue.Clear()

	r.LocalSequence = 0
	r.RemoteSequence = 0
	r.SentPackets = 0
	r.RecvPackets = 0
	r.LostPackets = 0
	r.AckedPackets = 0
	r.SentBandwidth = 0
	r.AckedBandwidth = 0
	r.RTT = 0
	r.RTTMax = 1
}