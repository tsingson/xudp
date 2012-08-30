// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

// Maximum packet sequence value.	
const MaxSequence = 1<<32 - 1

// isMoreRecent checks if sequence a is newer than sequence b,
// while taking integer overflow into account.
func isMoreRecent(a, b uint32) bool {
	const max = MaxSequence >> 1
	return (a > b) && (a-b <= max) || (b > a) && (b-a > max)
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
// packet queues. In addtion, it tracks bandwidth use and round trip timing.
type Reliability struct {
	OnAcked         PacketHandler // Notify the host when a specific packet is ACK'ed.
	OnLost          PacketHandler // Notify the host when a specific packet is lost.
	sentQueue       packetQueue   // Sent packets used to calculate sent bandwidth.
	pendingAckQueue packetQueue   // Sent packets which have not been acked yet.
	recvQueue       packetQueue   // Received packets used to determine acks to send.
	ackedQueue      packetQueue   // ACK'ed packets.
	SentBytes       uint64        // Number of bytes sent.
	RecvBytes       uint64        // Number of bytes received.
	SentPackets     uint32        // Number of packets sent.
	RecvPackets     uint32        // Number of packets received.
	LostPackets     uint32        // Number of packets lost.
	AckedPackets    uint32        // Number of packets ACK'ed.
	LocalSequence   uint32        // Local sequence number for most recently sent packet.
	RemoteSequence  uint32        // Remote sequence number for most recently received packet.
	SentBandwidth   float32       // Approximate sent bandwidth over the last second.
	AckedBandwidth  float32       // Approximate ACK'ed bandwidth over the last second.
	RTT             float32       // Estimated round trip time.
	RTTMax          float32       // Maximum expected round trip time.
}

// NewReliability creates a new reliability instance.
func NewReliability() *Reliability {
	r := new(Reliability)
	r.reset()
	return r
}

// PacketSent is called whenever a new packet is sent.
func (r *Reliability) PacketSent(size uint32) {
	var pd packetData
	pd.sequence = r.LocalSequence
	pd.size = size

	r.sentQueue = append(r.sentQueue, pd)
	r.pendingAckQueue = append(r.pendingAckQueue, pd)
	r.SentPackets++
	r.LocalSequence++
	r.SentBytes += uint64(size)
}

// PacketRecv is called whenever a new packet is received.
func (r *Reliability) PacketRecv(sequence, ack, vector, size uint32) {
	defer r.processAck(ack, vector)
	r.RecvPackets++
	r.RecvBytes += uint64(size)

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

// Update takes a frame time delta and updates packet timeouts with it.
func (r *Reliability) Update(delta float32) {
	r.advanceQueueTime(delta)
	r.updateQueues()
	r.updateStats()
}

// Reset sets the Reliability system to its initial state.
func (r *Reliability) reset() {
	r.sentQueue = r.sentQueue[:0]
	r.recvQueue = r.recvQueue[:0]
	r.pendingAckQueue = r.pendingAckQueue[:0]
	r.ackedQueue = r.ackedQueue[:0]

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

// processAck handles a single incoming ACK with ACK vector.
func (r *Reliability) processAck(ack, vector uint32) {
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
		r.AckedPackets++
		r.pendingAckQueue.RemoveAt(i)
		i--

		if r.OnAcked != nil {
			r.OnAcked(pd.sequence)
		}
	}
}

// AdvanceQueueTime updates the timestamp for each queued packet.
func (r *Reliability) advanceQueueTime(delta float32) {
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

// updateQueues updates all queues to discard packets that have
// exceeded their timeouts or are otherwise no longer necessary.
func (r *Reliability) updateQueues() {
	const epsilon = 0.001

	if len(r.recvQueue) > 0 {
		var minSeq uint32
		lastSeq := r.recvQueue[len(r.recvQueue)-1].sequence

		if lastSeq >= 34 {
			minSeq = lastSeq - 34
		} else {
			minSeq = MaxSequence - (34 - lastSeq)
		}

		for len(r.recvQueue) > 0 && isMoreRecent(minSeq, r.recvQueue[0].sequence) {
			r.recvQueue = r.recvQueue[1:]
		}
	}

	threshold := r.RTTMax + epsilon
	for len(r.sentQueue) > 0 && r.sentQueue[0].time > threshold {
		r.sentQueue = r.sentQueue[1:]
	}

	for len(r.pendingAckQueue) > 0 && r.pendingAckQueue[0].time > threshold {
		if r.OnLost != nil {
			r.OnLost(r.pendingAckQueue[0].sequence)
		}

		r.pendingAckQueue = r.pendingAckQueue[1:]
		r.LostPackets++
	}

	threshold = r.RTTMax*2 - epsilon
	for len(r.ackedQueue) > 0 && r.ackedQueue[0].time > threshold {
		r.ackedQueue = r.ackedQueue[1:]
	}
}

// updateStats updates bandwidth and timing statistics.
func (r *Reliability) updateStats() {
	var ackedBytesPerSec float32
	var sentBytesPerSec float32
	var pd packetData

	rm := r.RTTMax

	for _, pd = range r.sentQueue {
		sentBytesPerSec += float32(pd.size)
	}

	for _, pd = range r.ackedQueue {
		if pd.time >= rm {
			ackedBytesPerSec += float32(pd.size)
		}
	}

	sentBytesPerSec /= rm
	ackedBytesPerSec /= rm

	r.SentBandwidth = sentBytesPerSec * (8 / 1000.0)
	r.AckedBandwidth = ackedBytesPerSec * (8 / 1000.0)
}
