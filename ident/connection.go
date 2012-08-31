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
	port   int    // Connection's listen port.
}

// New creates a new connection.
func New(mtu, protocolId uint32, port int) *Connection {
	c := new(Connection)
	c.port = port
	c.Connection = xudp.New(mtu, protocolId)
	c.peerId = NewPeerId(LocalIP(), port)
	return c
}

// Open opens the connection on the connection's port number.
func (c *Connection) Open() (err error) { return c.Connection.Open(c.port) }

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
