// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package ident

import (
	"github.com/jteeuwen/xudp"
	"net"
)

// A connection allows reliable, two-way communication with an end point.
// It offers accurate identification of peers through a unique client id.
type Connection struct {
	*xudp.Connection
	peerId PeerId // This connection's peer id.
}

// New creates a new connection.
func New(mtu, protocolId uint32) *Connection {
	c := new(Connection)
	c.Connection = xudp.New(mtu, protocolId)
	return c
}

// PayloadSize returns the maximum size in bytes for a single paket payload.
func (c *Connection) PayloadSize() int {
	return c.Connection.PayloadSize() - PeerIdSize
}

// Open opens the connection on the connection's port number.
func (c *Connection) Open(port int) error {
	c.peerId = NewPeerId(LocalIP(), port)
	return c.Connection.Open(port)
}

// Send sends the given payload to the specified destination.
func (c *Connection) Send(addr net.Addr, payload []byte) (err error) {
	data := make([]byte, len(payload)+PeerIdSize)
	copy(data, c.peerId)
	copy(data[PeerIdSize:], payload)
	return c.Connection.Send(addr, data)
}

// Recv receives a new payload. This is a blocking operation.
func (c *Connection) Recv() (src *Endpoint, payload []byte, err error) {
	addr, payload, err := c.Connection.Recv()

	if err != nil {
		return
	}

	if len(payload) < PeerIdSize {
		err = xudp.ErrPacketSize
		return
	}

	src = NewEndpoint(addr, payload[:PeerIdSize])
	payload = payload[PeerIdSize:]
	return
}
