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

// A connection allows reliable, two-way communication with an end point.
// It can be run as either a client or server.
type Connection struct {
	sock    *socket        // The UDP socket for this connection.
	mode    ConnectionMode // Connection mode: Server or client.
	timeout uint           // Connection timeout in seconds.
	mtu     uint32         // Maximum packet size.
	proto   uint32         // Protocol identifier.
}

// New creates a new connection of the given type.
//
// The mode determines if we are using this connection as a server, or
// a client.
func New(mode ConnectionMode) *Connection {
	c := new(Connection)
	c.mode = mode
	c.mtu = 1400
	c.timeout = 3
	c.proto = 'X'<<24 | 'U'<<16 | 'D'<<8 | 'P'
	return c
}

// Close closes the connection.
func (c *Connection) Close() (err error) {
	if c.sock != nil {
		err = c.sock.Close()
		c.sock = nil
	}

	return
}

// Mode returns the connection mode for this connection.
func (c *Connection) Mode() ConnectionMode { return c.mode }

// SetMTU sets the maximum size of a single packet in bytes.
//
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
func (c *Connection) SetMTU(mtu uint32) { c.mtu = mtu }

// MTU returns the maximum size of a single packet in bytes.
func (c *Connection) MTU() uint32 { return c.mtu }

// SetProtocolId sets the protocol identifier.
//
// The protocol Id is a numerical identifier for all the packets
// sent and received by our program. It can be any number we want, but
// it is advised to use something relatively unique. It basically means:
// if an incoming packet does not start with this number, discard it because
// it is not meant for us.
func (c *Connection) SetProtocolId(id uint32) { c.proto = id }

// ProtocolId returns the protocol identifier.
func (c *Connection) ProtocolId() uint32 { return c.proto }

// SetTimeout sets the connection timeout in seconds.
func (c *Connection) SetTimeout(t uint) { c.timeout = t }

// Timeout returns the connection timeout in seconds.
func (c *Connection) Timeout() uint { return c.timeout }
