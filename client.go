// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

// Client represents a single client connection.
type Client struct {
	*ReliableConnection
}

// NewClient creates a new, uninitialized client.
func NewClient(mtu, protocolId uint32) *Client {
	c := new(Client)
	c.ReliableConnection = NewReliableConnection(mtu, protocolId)
	return c
}
