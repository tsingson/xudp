// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

type ConnectionMode uint8

const (
	ModeClient ConnectionMode = iota
	ModeServer
)

// A connection allows reliable, two-way communication with an end point.
// It can be run as either a client or server.
type Connection struct {
	protocolId uint32         // Protocol ID used by this connection.
	mode       ConnectionMode // Connection mode: Server or client.
}

// NewConnection creates a new connection of the given type.
func NewConnection(mode ConnectionMode, protocolId uint32) *Connection {
	c := new(Connection)
	c.mode = mode
	c.protocolId = protocolId
	return c
}

// Close closes the connection.
func (c *Connection) Close() error {
	return nil
}
