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
	ErrSocketOpen   = errors.New("socket is already open.")
	ErrSocketClosed = errors.New("socket is already closed.")
)

// socket handles low level UDP socket mechanics.
//
// For normal use cases, there is no need to use this type directly.
// Refer to one of the higher level types instead.
type socket struct {
	udp        net.PacketConn // Sockets underlying connection.
	protocolId uint32         // Protocol ID identifying our packets.
	mtu        uint32         // Maximum transport unit.
}

// newSocket creates a new, uninitialized socket.
func newSocket(mtu, protocolId uint32) *socket {
	s := new(socket)
	s.protocolId = protocolId
	s.mtu = mtu
	return s
}

// Open opens the socket on the given port number.
func (s *socket) Open(port uint) (err error) {
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
	return
}

// Close closes the socket.
func (s *socket) Close() (err error) {
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

	_, err = s.udp.WriteTo(payload, dest)
	return
}

// Poll returns a channel on which new packets are received.
//
// The received packet data refers to the socket's internal buffer
// and will only remain valid until the next packet is received.
// It is up to the caller to create a copy of the data when needed.
//
// It yields only those packets that match the socket's protocol Id.
func (s *socket) Poll() <-chan Packet {
	c := make(chan Packet)

	go func() {
		var addr net.Addr
		var err error
		var size int

		defer close(c)

		packet := make(Packet, (s.mtu-UDPHeaderSize)+XUDPAddrSize)

		for {
			if s.udp == nil {
				return
			}

			size, addr, err = s.udp.ReadFrom(packet[XUDPAddrSize:])
			if err != nil {
				return
			}

			if size < XUDPHeaderSize {
				continue // Not enough data.
			}

			if packet.Protocol() != s.protocolId {
				continue // Not meant for us.
			}

			packet.setAddr(addr.(*net.UDPAddr))
			c <- packet[:XUDPAddrSize+size]
		}
	}()

	return c
}

// IsOpen returns true if the connection is currently open.
func (s *socket) IsOpen() bool { return s.udp != nil }
