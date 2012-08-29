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
)

// A connection allows reliable, two-way communication with an end point.
type Connection struct {
	*Reliability
	buf   Packet         // Temporary receive buffer.
	udp   net.PacketConn // Sockets underlying connection.
	proto uint32         // Protocol ID identifying our packets.
}

// NewConnection creates a new connection.
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
func NewConnection(mtu, protocolId uint32) *Connection {
	c := new(Connection)
	c.Reliability = NewReliability()
	c.proto = protocolId
	c.buf = make(Packet, mtu-UDPHeaderSize)
	return c
}

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

// Send sends the given packet to the specified destination.
func (c *Connection) Send(addr net.Addr, packet Packet) (size int, err error) {
	if c.udp == nil {
		return 0, ErrConnectionClosed
	}

	if len(packet) > len(c.buf) {
		return 0, ErrPacketSize
	}

	packet.SetHeader(c.proto, c.LocalSequence, c.RemoteSequence, c.AckVector())

	size, err = c.udp.WriteTo(packet, addr)

	if err != nil {
		return
	}

	c.PacketSent(size)
	return
}

// Recv receives a new packet. This is a blocking operation.
func (c *Connection) Recv() (addr net.Addr, packet Packet, err error) {
	if c.udp == nil {
		return nil, nil, ErrConnectionClosed
	}

	size, addr, err := c.udp.ReadFrom(c.buf)
	if err != nil {
		return
	}

	if size < XUDPHeaderSize || c.buf.Protocol() != c.proto {
		return // Not enough data or not meant for us.
	}

	packet = make(Packet, size)
	copy(packet, c.buf)
	return
}
