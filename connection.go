// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

// ConnectionMode defines the purpose of a given connection.
type ConnectionMode uint8

// Available connection modes.
const (
	Client ConnectionMode = iota
	Server
)

type ConnectionHandler func(*Connection) error

// A connection allows reliable, two-way communication with an end point.
// It can be run as either a client or server.
type Connection struct {
	sock    *socket           // The underlying UDP socket for this connection.
	OnOpen  ConnectionHandler // Handler called when a connection is established.
	OnClose ConnectionHandler // Handler called when the connection closes.
	Timeout uint              // Timeout defines the connection timeout in seconds.
	mode    ConnectionMode    // Connection mode: Server or client.
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
// it is advised to use something relatively unique. It basically means:
// if an incoming packet does not start with this number, discard it
// because it is not meant for us.
func New(mode ConnectionMode, mtu, protocolId uint32) *Connection {
	c := new(Connection)
	c.sock = newSocket(mtu, protocolId)
	c.mode = mode
	c.Timeout = 3
	return c
}

// Open opens the connection on the given port.
func (c *Connection) Open(port uint) (err error) {
	if err = c.sock.Open(port); err != nil {
		return
	}

	if c.OnOpen != nil {
		err = c.OnOpen(c)
	}

	return
}

// Close closes the connection.
func (c *Connection) Close() (err error) {
	if c.sock != nil {
		err = c.sock.Close()
		c.sock = nil

		if err != nil {
			return
		}
	}

	if c.OnClose != nil {
		err = c.OnClose(c)
	}

	c.OnClose = nil
	c.OnOpen = nil
	return
}
