// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"net"
	"time"
)

// A reliable connection allows reliable, two-way communication with
// an end point. It functions as both a client and server at the same time.
// It deals with dropped packet retransmission.
type ReliableConnection struct {
	*Connection
	*Reliability
	Timeout    time.Duration // Timeout defines the connection timeout in seconds.
	LastPacket int64         // Timestamp for last received packet.
}

// NewReliableConnection creates a new reliable connection.
//
// Refer to the documentation on NewConnection for details on
// what the parameters mean.
func NewReliableConnection(mtu, protocolId uint32) *ReliableConnection {
	rc := new(ReliableConnection)
	rc.Connection = NewConnection(mtu, protocolId)
	rc.Reliability = NewReliability()
	rc.Timeout = time.Second * 3
	return rc
}

// Recv receives a new packet. This is a blocking operation.
func (rc *ReliableConnection) Recv() (addr net.Addr, packet Packet, err error) {
	addr, packet, err = rc.Connection.Recv()

	if err != nil {
		return
	}

	rc.LastPacket = time.Now().UTC().UnixNano()
	return
}
