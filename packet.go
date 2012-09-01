// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

const (
	// This is the size of a standard UDP datagram header in bytes.
	// It is part of every packet we send. This header is processed by the
	// operating system's transport layer; we will therefore never see it.
	// Note that this header counts towards the maximum size for a single
	// UDP datagram.
	UDPHeaderSize = 22

	// Minimum size for an XUDP header in bytes. It is part of every
	// packet we send and receive. This counts towards the maximum size for
	// a single UDP datagram.
	XUDPHeaderSize = 16
)

// A PacketHandler is used to notify the host of
// ACK'ed or lost packets by their sequence number.
type PacketHandler func(sequence uint32)
