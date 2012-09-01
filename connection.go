// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"errors"
	"fmt"
	"net"
	"time"
)

var (
	ErrConnectionOpen   = errors.New("Connection is already open.")
	ErrConnectionClosed = errors.New("Connection is already closed.")
	ErrPacketSize       = errors.New("Packet size exceeds MTU.")
	ErrShortWrite       = errors.New("Short write: Send was incomplete.")
)

// A connection allows reliable, two-way communication with an end point.
type Connection struct {
	*Reliability
	inbuf      []byte         // Temporary receive buffer.
	outbuf     []byte         // Temporary output buffer.
	udp        net.PacketConn // Underlying socket.
	protocolId uint32         // Protocol Id identifying our packets.
	mtu        uint32         // maximum packet size.
}

// New creates a new connection.
//
// MTU defines the maximum size of a single packet in bytes.
// This includes the UDP and XUDP headers.
// The available payload space can be calculated as:
//
//     payloadSize := MTU - UDPHeaderSize - XUDPHeaderSize
//
// Some commonly used values are as follows:
//
//     1500 - The largest Ethernet packet size. This is the typical setting for
//            non-PPPoE, non-VPN connections. The default value for NETGEAR
//            routers, adapters and switches.
//     1492 - The size PPPoE prefers.
//     1472 - Maximum size to use for pinging (Bigger packets are fragmented).
//     1468 - The size DHCP prefers.
//     1460 - Usable by AOL if you don't have large email attachments, etc.
//     1430 - The size VPN and PPTP prefer.
//     1400 - Maximum size for AOL DSL.
//      576 - Typical value to connect to dial-up ISPs.
//
// The protocol Id is a numerical identifier for all the packets
// sent and received by our program. It can be any number we want, but
// it is advised to use something relatively unique.
// If an incoming packet does not start with this number, we discard it
// because it is not meant for us. A 4 byte hash of the name of your
// program can be a suitable protocol Id.
func New(mtu, protocolId uint32) *Connection {
	c := new(Connection)
	c.Reliability = NewReliability()
	c.protocolId = protocolId
	c.mtu = mtu
	c.inbuf = make([]byte, mtu-UDPHeaderSize)
	c.outbuf = make([]byte, mtu-UDPHeaderSize)
	return c
}

// PayloadSize returns the maximum size in bytes for a single paket payload.
func (c *Connection) PayloadSize() int {
	return int(c.mtu) - UDPHeaderSize - XUDPHeaderSize
}

func (c *Connection) IsOpen() bool { return c.udp != nil }

// Open opens the connection on the given port number.
func (c *Connection) Open(port int) (err error) {
	if c.udp != nil {
		return ErrConnectionOpen
	}

	c.udp, err = net.ListenPacket("udp", fmt.Sprintf(":%d", port))

	if err != nil {
		return
	}

	zero := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	c.udp.SetReadDeadline(zero)
	c.udp.SetWriteDeadline(zero)
	return
}

// Close closes the connection.
func (c *Connection) Close() (err error) {
	if c.udp == nil {
		return ErrConnectionClosed
	}

	err = c.udp.Close()
	c.udp = nil
	return
}

// Send sends the given payload to the specified destination.
func (c *Connection) Send(addr net.Addr, payload []byte) (err error) {
	if c.udp == nil {
		return ErrConnectionClosed
	}

	if len(payload) > c.PayloadSize() {
		return ErrPacketSize
	}

	b := c.outbuf
	packetSize := XUDPHeaderSize + len(payload)

	// Construct a new packet with the appropriate header fields.
	n := c.protocolId
	b[0] = byte(n >> 24)
	b[1] = byte(n >> 16)
	b[2] = byte(n >> 8)
	b[3] = byte(n)

	n = c.LocalSequence
	b[4] = byte(n >> 24)
	b[5] = byte(n >> 16)
	b[6] = byte(n >> 8)
	b[7] = byte(n)

	n = c.RemoteSequence
	b[8] = byte(n >> 24)
	b[9] = byte(n >> 16)
	b[10] = byte(n >> 8)
	b[11] = byte(n)

	n = c.AckVector()
	b[12] = byte(n >> 24)
	b[13] = byte(n >> 16)
	b[14] = byte(n >> 8)
	b[15] = byte(n)

	copy(b[XUDPHeaderSize:], payload)

	size, err := c.udp.WriteTo(b[:packetSize], addr)

	if err != nil {
		return
	}

	if size < packetSize {
		return ErrShortWrite
	}

	c.PacketSent(uint32(size))
	return
}

// Recv receives a new payload. This is a blocking operation.
func (c *Connection) Recv() (addr net.Addr, payload []byte, err error) {
	if c.udp == nil {
		return nil, nil, ErrConnectionClosed
	}

	b := c.inbuf
	size, addr, err := c.udp.ReadFrom(b)

	if err != nil {
		return
	}

	if size < XUDPHeaderSize {
		return // Not enough data.
	}

	proto := uint32(b[0])<<24 | uint32(b[1])<<16 |
		uint32(b[2])<<8 | uint32(b[3])

	if proto != c.protocolId {
		return // Not meant for us.
	}

	sequence := uint32(b[4])<<24 | uint32(b[5])<<16 |
		uint32(b[6])<<8 | uint32(b[7])

	ack := uint32(b[8])<<24 | uint32(b[9])<<16 |
		uint32(b[10])<<8 | uint32(b[11])

	vector := uint32(b[12])<<24 | uint32(b[13])<<16 |
		uint32(b[14])<<8 | uint32(b[15])

	c.PacketRecv(sequence, ack, vector, uint32(size))

	size -= XUDPHeaderSize
	if size <= 0 {
		return // No payload data.
	}

	payload = make([]byte, size)
	copy(payload, b[XUDPHeaderSize:XUDPHeaderSize+size])
	return
}
