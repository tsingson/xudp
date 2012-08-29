// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"net"
)

// Server represents a server connection.
type Server struct {
	*Connection
	clients    []*Client // List of active clients.
	MaxClients uint      // Maximum number of clients we can accept at any given time.
}

// NewServer creates a new, uninitialized server.
func NewServer(mtu, protocolId uint32) *Server {
	s := new(Server)
	s.Connection = NewConnection(mtu, protocolId)
	s.MaxClients = 100
	return s
}

// Recv receives a new packet. This is a blocking operation.
func (s *Server) Recv() (addr net.Addr, packet Packet, err error) {
	addr, packet, err = s.Connection.Recv()

	if err != nil {
		return
	}

	return
}
