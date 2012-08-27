// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

type ConnectionMode uint8

const (
	ModeClient ConnectionMode = iota
	ModeServer
)

// A connection allows reliable, two-way communication with an end point.
// It can be run as either a client or server.
type Connection struct {
	mode ConnectionMode // Connection mode: Server or client.
}

// NewConnection creates a new connection of the given type.
//
// The mode determines if we are using this connection as a server, or
// a client.
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
// protocolId defines the socket's protocol identifier.
//
// The protocol Id is a numerical identifier for all the packets
// sent and received by our program. It can be any number we want, but
// it is advised to use something relatively unique.
//
// It basically means: if an incoming packet does not start with this
// number, discard it because it is not meant for us.
func NewConnection(mode ConnectionMode, mtu, protocolId uint32) *Connection {
	c := new(Connection)
	c.mode = mode
	return c
}

// Close closes the connection.
func (c *Connection) Close() error {
	return nil
}
