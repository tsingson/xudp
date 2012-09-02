// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import "net"

// A Plugin can be registered with a connection to add
// a unique feature to a packet connection.
type Plugin interface {
	// Returns the size in bytes of any data added to the packet
	// by the given plugin.
	PayloadSize() int

	Open(port int) error
	Close() error

	// Called when a new packet is being sent.
	//
	// It accepts the target address, the full packet which includes the
	// plugin/header data and the payload and the index at which the actual
	// payload starts.
	Send(net.Addr, []byte, int) error

	// Called when a new packet is received.
	//
	// It accepts the source address, the full packet which includes the
	// plugin/header data and the payload and the index at which the actual
	// payload starts.
	Recv(net.Addr, []byte, int) error
}
