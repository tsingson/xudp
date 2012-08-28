// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

const (
	// This is the size of a standard UDP datagram header in bytes.
	// It is part of every packet we send. This header is processed by the
	// operating system's transport layer; we will therefore never see it.
	// Note that this header counts towards the maximum size for a single
	// UDP datagram.
	UDPHeaderSize = 22

	// Minimum size for an XUDP header in bytes. It is part of every
	// packet we send and receive. This counts towards the maximum size for
	// a single UDP datagram.
	XUDPHeaderSize = 16
)

// A packet holds data for a single UDP datagram.
// This includes our own XUDP header and the payload.
type Packet []byte

// NewPacket creates a new packet with the given payload and
// enough header space to fit all header fields.
// The header fields have yet to be set.
func NewPacket(payload []byte) Packet {
	p := make(Packet, XUDPHeaderSize+len(payload))
	copy(p[XUDPHeaderSize:], payload)
	return p
}

// Protocol returns the 32 bit, unsigned protocol id.
func (p Packet) Protocol() uint32 {
	return uint32(p[0])<<24 | uint32(p[1])<<16 | uint32(p[2])<<8 | uint32(p[3])
}

// SetProtocol sets the 32 bit, unsigned protocol id.
func (p Packet) SetProtocol(proto uint32) {
	p[0] = byte(proto >> 24)
	p[1] = byte(proto >> 16)
	p[2] = byte(proto >> 8)
	p[3] = byte(proto)
}

// Sequence returns the 32 bit, unsigned sequence number for this packet.
func (p Packet) Sequence() uint32 {
	return uint32(p[4])<<24 | uint32(p[5])<<16 | uint32(p[6])<<8 | uint32(p[7])
}

func (p Packet) SetSequence(sequence uint32) {
	p[4] = byte(sequence >> 24)
	p[5] = byte(sequence >> 16)
	p[6] = byte(sequence >> 8)
	p[7] = byte(sequence)
}

// Ack returns the 32 bit, unsigned sequence number for an acknowledged packet.
// We incorporate this in the header, so ACKS can piggyback on regular data packets.
func (p Packet) Ack() uint32 {
	return uint32(p[8])<<24 | uint32(p[9])<<16 | uint32(p[10])<<8 | uint32(p[11])
}

func (p Packet) SetAck(ack uint32) {
	p[8] = byte(ack >> 24)
	p[9] = byte(ack >> 16)
	p[10] = byte(ack >> 8)
	p[11] = byte(ack)
}

// AckVector returns a 32 bit, unsigned bitset for additional ACKs.
// This allows us to encode up to 33 simultaneous ACKs in a single packet,
// using only one ACK sequence number and a 32 bit bitset.
//
// This approach allows for a large amount of redundancy in ACK handling
// when we are dealing with two peers sending data at different rates.
// Packets can start to pile up on one side, or get lost all together, 
// messing up a 1:1 ACK approach. The redundancy this bitfield yields,
// solves that problem.
//
// To illustrate: If Ack() == 100, we are acknowledging reception of the
// 100th packet. We can ACK packets 99, 98, 97, ..., 68 in one go, by setting
// each individual bit in this bitfield. If bit 1 is set, then we
// ACK packet 99. Bit 2 ACKs packet 98, etc.
func (p Packet) AckVector() uint32 {
	return uint32(p[12])<<24 | uint32(p[13])<<16 | uint32(p[14])<<8 | uint32(p[15])
}

func (p Packet) SetAckVector(vector uint32) {
	p[12] = byte(vector >> 24)
	p[13] = byte(vector >> 16)
	p[14] = byte(vector >> 8)
	p[15] = byte(vector)
}

// Payload returns the packet data.
func (p Packet) Payload() []byte { return p[XUDPHeaderSize:] }
