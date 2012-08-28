// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

// A reliable connection allows reliable, two-way communication with
// an end point. It functions as both a client and server at the same time.
// It deals with dropped packet retransmission.
type ReliableConnection struct {
	*Connection
}

// NewReliableConnection creates a new reliable connection.
//
// Refer to the documentation on NewConnection for details on
// what the parameters mean.
func NewReliableConnection(mtu, protocolId uint32) *ReliableConnection {
	c := new(ReliableConnection)
	c.Connection = NewConnection(mtu, protocolId)
	return c
}
