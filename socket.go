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
	ErrSocketOpen   = errors.New("Socket is already open.")
	ErrSocketClosed = errors.New("Socket is already closed.")
	ErrSendFailed   = errors.New("Send was incomplete.")
)

// socket handles low level UDP socket mechanics.
type socket struct {
	Recv  chan packet    // Channel for incoming packets.
	udp   net.PacketConn // Sockets underlying connection.
	proto uint32         // Protocol ID identifying our packets.
	mtu   uint32         // Maximum transport unit.
}

// newSocket creates a new, uninitialized socket.
func newSocket(mtu, proto uint32) *socket {
	s := new(socket)
	s.proto = proto
	s.mtu = mtu
	s.Recv = make(chan packet)
	return s
}

// Open opens the socket on the given port number.
func (s *socket) Open(port int) (err error) {
	if s.udp != nil {
		return ErrSocketOpen
	}

	s.udp, err = net.ListenPacket("udp", fmt.Sprintf(":%d", port))

	if err != nil {
		return
	}

	zero := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	s.udp.SetReadDeadline(zero)
	s.udp.SetWriteDeadline(zero)

	go s.poll()
	return
}

// Close closes the socket.
func (s *socket) Close() (err error) {
	close(s.Recv)

	if s.udp == nil {
		return ErrSocketClosed
	}

	err = s.udp.Close()
	s.udp = nil
	return
}

// Send sends the given data to the specified destination.
func (s *socket) Send(dest net.Addr, payload []byte) (err error) {
	if len(payload) == 0 {
		return
	}

	if s.udp == nil {
		return ErrSocketClosed
	}

	max := s.mtu - UDPHeaderSize

	if uint32(len(payload)) > max {
		payload = payload[:max]
	}

	sent, err := s.udp.WriteTo(payload, dest)

	if err != nil {
		return
	}

	if sent != len(payload) {
		err = ErrSendFailed
	}

	return
}

// poll checks for incoming data.
func (s *socket) poll() {
	buf := make(packet, (s.mtu-UDPHeaderSize)+XUDPAddrSize)

	for {
		if s.udp == nil {
			return
		}

		size, addr, err := s.udp.ReadFrom(buf[XUDPAddrSize:])
		if err != nil {
			return
		}

		if size < XUDPHeaderSize {
			continue // Not enough data.
		}

		if buf.Protocol() != s.proto {
			continue // Not meant for us.
		}

		p := make(packet, XUDPAddrSize+size)
		p.SetAddr(addr.(*net.UDPAddr))
		copy(p[XUDPAddrSize:], buf[XUDPAddrSize:])

		s.Recv <- p
	}
}

// IsOpen returns true if the connection is currently open.
func (s *socket) IsOpen() bool { return s.udp != nil }
