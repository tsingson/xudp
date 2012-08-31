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
	buf        Packet         // Temporary receive buffer.
	udp        net.PacketConn // Underlying socket.
	protocolId uint32         // Protocol Id identifying our packets.
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
	c.buf = make(Packet, mtu-UDPHeaderSize)
	return c
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

	if len(payload) > len(c.buf)-XUDPHeaderSize {
		return ErrPacketSize
	}

	packet := NewPacket(payload)
	packet.SetHeader(c.protocolId, c.LocalSequence, c.RemoteSequence, c.AckVector())

	size, err := c.udp.WriteTo(packet, addr)

	if err != nil {
		return
	}

	if size < len(packet) {
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

	size, addr, err := c.udp.ReadFrom(c.buf)
	if err != nil {
		return
	}

	if size < XUDPHeaderSize || c.buf.Protocol() != c.protocolId {
		return // Not enough data or not meant for us.
	}

	c.PacketRecv(c.buf.Sequence(), c.buf.Ack(), c.buf.AckVector(), uint32(size))

	size -= XUDPHeaderSize
	if size <= 0 {
		return // No payload data.
	}

	payload = make([]byte, size)
	copy(payload, c.buf[XUDPHeaderSize:XUDPHeaderSize+size])
	return
}
