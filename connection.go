// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"errors"
	"fmt"
	"net"
	"time"
)

// This is the size of a standard UDP datagram header in bytes.
// It is part of every packet we send. This header is processed by the
// operating system's transport layer; we will therefore never see it.
// Note that this header counts towards the maximum size for a single
// UDP datagram.
const UDPHeaderSize = 22

var (
	ErrConnectionectionOpen   = errors.New("Connectionection is already open.")
	ErrConnectionectionClosed = errors.New("Connectionection is already closed.")
	ErrPacketSize             = errors.New("Packet size exceeds MTU.")
	ErrShortWrite             = errors.New("Short write: Send was incomplete.")
	ErrDiscard                = errors.New("Packet is not meant for us.")
)

// A connection allows two-way communication with an end point.
//
// Without any registered plugins, this is really nothing more than
// a wrapper around the standard Go UDP connection tools. 
//
// The plugins you can add, give it added complexity and usefulness
// in a way specific to your application. By using the stackable
// building block approach, you are never stuck with a connection that
// does more than you need it to do.
type Connection struct {
	PluginList
	inbuf  []byte         // Temporary receive buffer.
	outbuf []byte         // Temporary output buffer.
	udp    net.PacketConn // Underlying socket.
	mtu    uint32         // maximum packet size.
}

// New creates a new connection.
//
// MTU defines the maximum size of a single packet in bytes.
// This includes the UDP and XUDP plugin headers.
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
func New(mtu uint32) *Connection {
	c := new(Connection)
	c.mtu = mtu
	c.inbuf = make([]byte, mtu-UDPHeaderSize)
	c.outbuf = make([]byte, mtu-UDPHeaderSize)
	return c
}

// PayloadSize returns the maximum size in bytes for a single packet payload.
// This is the MTU minus the UDP header and the space required by all
// registered plugins.
func (c *Connection) PayloadSize() int {
	size := int(c.mtu) - UDPHeaderSize

	for _, plg := range c.PluginList {
		size -= plg.PayloadSize()
	}

	return size
}

// Open opens the connection on the given port number.
func (c *Connection) Open(port int) (err error) {
	if c.udp != nil {
		return ErrConnectionectionOpen
	}

	c.udp, err = net.ListenPacket("udp", fmt.Sprintf(":%d", port))

	if err != nil {
		return
	}

	zero := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	c.udp.SetReadDeadline(zero)
	c.udp.SetWriteDeadline(zero)

	for _, plg := range c.PluginList {
		err = plg.Open(port)

		if err != nil {
			c.Close()
			return
		}
	}

	return
}

// Close closes the connection.
func (c *Connection) Close() (err error) {
	if c.udp == nil {
		return ErrConnectionectionClosed
	}

	err = c.udp.Close()
	c.udp = nil

	for _, plg := range c.PluginList {
		plg.Close()
	}

	c.PluginList.Clear()
	c.PluginList = nil
	return
}

// Send sends the given payload to the specified destination.
func (c *Connection) Send(addr net.Addr, payload []byte) (err error) {
	if c.udp == nil {
		return ErrConnectionectionClosed
	}

	if len(payload) > c.PayloadSize() {
		return ErrPacketSize
	}

	var index int

	b := c.outbuf
	header := c.PluginList.PayloadSize()
	total := header + len(payload)

	copy(b[header:], payload)

	for _, plg := range c.PluginList {
		err = plg.Send(addr, b[index:total], header)

		if err != nil {
			return
		}

		index += plg.PayloadSize()
	}

	size, err := c.udp.WriteTo(b[:total], addr)
	if err != nil {
		return
	}

	if size < total {
		err = ErrShortWrite
	}

	return
}

// Recv receives a new payload. This is a blocking operation.
func (c *Connection) Recv() (addr net.Addr, payload []byte, err error) {
	if c.udp == nil {
		return nil, nil, ErrConnectionectionClosed
	}

	b := c.inbuf
	size, addr, err := c.udp.ReadFrom(b)

	if err != nil {
		return
	}

	header := c.PluginList.PayloadSize()
	if size < header {
		return // Not enough data.
	}

	var packetSize int
	for _, plg := range c.PluginList {
		err = plg.Recv(addr, b[packetSize:size], header)

		if err != nil {
			if err == ErrDiscard {
				err = nil // No need to propagate this.
			}
			return
		}

		packetSize += plg.PayloadSize()
	}

	size -= packetSize
	if size <= 0 {
		return // No payload data.
	}

	payload = make([]byte, size)
	copy(payload, b[packetSize:packetSize+size])
	return
}
