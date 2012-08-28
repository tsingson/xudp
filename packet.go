// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import "net"

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

	// Size of an address and port in bytes.
	XUDPAddrSize = 20
)

// A packet holds data for a single UDP datagram.
// This includes a 16 byte sender address, a 4 byte port, our own XUDP header
// and the payload.
//
// Note that the address and port are never sent/received.
// We include it in a packet manually to make packet handling simpler.
type packet []byte

// setAddr sets the address and port.
func (p packet) setAddr(addr *net.UDPAddr) {
	copy(p, addr.IP.To16())

	port := addr.Port
	p[16] = byte(port >> 24)
	p[17] = byte(port >> 16)
	p[18] = byte(port >> 8)
	p[19] = byte(port)
}

// setHeader builds the packet header.
// The packet must have a minimum size of XUDPAddrSize+XUDPHeaderSize.
func (p packet) setHeader(protocol, sequence, ack, ackvector uint32) {
	if len(p) < XUDPAddrSize+XUDPHeaderSize {
		return
	}

	p[20] = byte(protocol >> 24)
	p[21] = byte(protocol >> 16)
	p[22] = byte(protocol >> 8)
	p[23] = byte(protocol)

	p[24] = byte(sequence >> 24)
	p[25] = byte(sequence >> 16)
	p[26] = byte(sequence >> 8)
	p[27] = byte(sequence)

	p[28] = byte(ack >> 24)
	p[29] = byte(ack >> 16)
	p[30] = byte(ack >> 8)
	p[31] = byte(ack)

	p[32] = byte(ackvector >> 24)
	p[33] = byte(ackvector >> 16)
	p[34] = byte(ackvector >> 8)
	p[35] = byte(ackvector)
}

// Addr returns the sender's address and port.
func (p packet) Addr() *net.UDPAddr {
	a := new(net.UDPAddr)
	a.IP = net.IP(p[:16])
	a.Port = int(p[16])<<24 | int(p[17])<<16 | int(p[18])<<8 | int(p[19])
	return a
}

// Protocol returns the 32 bit, unsigned protocol Id.
func (p packet) Protocol() uint32 {
	return uint32(p[20])<<24 | uint32(p[21])<<16 | uint32(p[22])<<8 | uint32(p[23])
}

// Sequence returns the 32 bit, unsigned sequence number for this packet.
func (p packet) Sequence() uint32 {
	return uint32(p[24])<<24 | uint32(p[25])<<16 | uint32(p[26])<<8 | uint32(p[27])
}

// Ack returns the 32 bit, unsigned sequence number for an acknowledged packet.
// We incorporate this in the header, so ACKS can piggyback on regular data packets.
func (p packet) Ack() uint32 {
	return uint32(p[28])<<24 | uint32(p[29])<<16 | uint32(p[30])<<8 | uint32(p[31])
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
func (p packet) AckVector() uint32 {
	return uint32(p[32])<<24 | uint32(p[33])<<16 | uint32(p[34])<<8 | uint32(p[35])
}

// Payload returns the packet data.
func (p packet) Payload() []byte { return p[XUDPAddrSize+XUDPHeaderSize:] }
