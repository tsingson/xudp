// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package ident

import (
	"crypto/sha256"
	"github.com/jteeuwen/xudp"
	"net"
)

type PeerFunc func(hash PeerHash, payload []byte)

type Plugin struct {
	id     []byte   // Connection's peer id.
	onRecv PeerFunc // Receive handler.
}

// New creates a new peer id plugin.
//
// The given receive handler is fired whenever a new packet arrives.
// It associates the payload with a unique Peer hash, which can
// be used to identify the client.
func New(recv PeerFunc) xudp.Plugin {
	p := new(Plugin)
	p.onRecv = recv
	return p
}

func (c *Plugin) PayloadSize() int { return PeerHashSize }
func (c *Plugin) Close() error     { return nil }

func (c *Plugin) Open(port int) error {
	ip := localIP()
	data := make([]byte, len(ip)+4)

	data[0] = byte(port >> 24)
	data[1] = byte(port >> 16)
	data[2] = byte(port >> 8)
	data[3] = byte(port)
	copy(data[4:], ip)

	hm := sha256.New()
	hm.Write(data)
	c.id = hm.Sum(nil)
	return nil
}

func (p *Plugin) Send(addr net.Addr, payload []byte, index int) error {
	copy(payload, p.id)
	return nil
}

func (p *Plugin) Recv(addr net.Addr, payload []byte, index int) error {
	if p.onRecv == nil {
		return nil
	}

	p.onRecv(NewPeerHash(addr, payload[:PeerHashSize]), payload[index:])
	return nil
}

// localIP returns the first available local IP address.
// This is the subnet address if the host is located in a subnet.
func localIP() net.IP {
	// Connect to a random machine somewhere. It's irrelevant
	// where to, as long as it's not the loopback address.

	addr, err := net.ResolveUDPAddr("udp", "192.168.0.0:0")

	if err != nil {
		return nil
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		return nil
	}

	addr = conn.LocalAddr().(*net.UDPAddr)
	conn.Close()

	return addr.IP
}
