// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"net"
)

// A connection allows reliable, two-way communication with an end point.
// It functions as both a client and server at the same time.
type Connection struct {
	sock    *socket // The underlying UDP socket for this connection.
	Timeout uint    // Timeout defines the connection timeout in seconds.
}

// New creates a new connection of the given type.
//
// MTU defines the maximum size of a single packet in bytes.
// This includes the UDP and XUDP headers.
// The available payload space can be calculated as:
//
//     payloadSize := MTU - UDPHeaderSize - XUDPHeaderSize
//
// Some commonly used values are as follows:
//
//     1500 - The largest Ethernet packet size. This is the typical setting for
//            non-PPPoE, non-VPN connections. The default value for NETGEAR
//            routers, adapters and switches.
//     1492 - The size PPPoE prefers.
//     1472 - Maximum size to use for pinging (Bigger packets are fragmented).
//     1468 - The size DHCP prefers.
//     1460 - Usable by AOL if you don't have large email attachments, etc.
//     1430 - The size VPN and PPTP prefer.
//     1400 - Maximum size for AOL DSL.
//      576 - Typical value to connect to dial-up ISPs.
//
// The protocol Id is a numerical identifier for all the packets
// sent and received by our program. It can be any number we want, but
// it is advised to use something relatively unique.
// If an incoming packet does not start with this number, discard it
// because it is not meant for us. A 4 byte hash of the name of your
// program can be a suitable protocol ID.
func New(mtu, protocolId uint32) *Connection {
	c := new(Connection)
	c.sock = newSocket(mtu, protocolId)
	c.Timeout = 3
	return c
}

// Open opens the connection on the given port.
// The port defines the port on which we listen for incoming connections.
func (c *Connection) Open(port int) (err error) {
	return c.sock.Open(port)
}

// Close closes the connection.
func (c *Connection) Close() (err error) {
	if c.sock != nil {
		err = c.sock.Close()
		c.sock = nil
	}
	return
}

// Send sends the given payload to the specified end point.
//
// The size of the data may exceed the maximum packet size.
// This routine will automatically use packet fragmentation in this case.
func (c *Connection) Send(endpoint *net.UDPAddr, data []byte) (err error) {

	return
}
