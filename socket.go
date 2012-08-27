// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"errors"
	"fmt"
	"net"
)

var (
	ErrSocketOpen   = errors.New("Socket is already open.")
	ErrSocketClosed = errors.New("Socket is already closed.")
)

// Socket handles low level UDP socket mechanics.
type Socket struct {
	udp *net.UDPConn
}

func NewSocket() *Socket {
	s := new(Socket)
	return s
}

// Open opens the socket on the given port number.
func (s *Socket) Open(port uint) (err error) {
	if s.udp != nil {
		return ErrSocketOpen
	}

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}

	s.udp, err = net.ListenUDP("udp", addr)
	return
}

// Close closes the socket.
func (s *Socket) Close() (err error) {
	if s.udp == nil {
		return ErrSocketClosed
	}

	err = s.udp.Close()
	s.udp = nil
	return
}

// IsOpen returns true if the connection is currently open.
func (s *Socket) IsOpen() bool {
	return false
}

// Send sends the given data to the specified destination.
func (s *Socket) Send(dest *net.UDPAddr, payload []byte) (err error) {

	return
}
