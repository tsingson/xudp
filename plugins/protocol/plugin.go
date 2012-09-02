// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package protocol

import (
	"github.com/jteeuwen/xudp"
	"net"
)

type Plugin struct {
	proto uint32
}

// New creates a new protocolid plugin with the given Id value.
func New(protocolId uint32) xudp.Plugin {
	p := new(Plugin)
	p.proto = protocolId
	return p
}

func (p *Plugin) PayloadSize() int    { return 4 }
func (c *Plugin) Open(port int) error { return nil }
func (c *Plugin) Close() error        { return nil }

func (p *Plugin) Send(addr net.Addr, payload []byte, index int) error {
	n := p.proto
	payload[0] = byte(n >> 24)
	payload[1] = byte(n >> 16)
	payload[2] = byte(n >> 8)
	payload[3] = byte(n)
	return nil
}

func (p *Plugin) Recv(addr net.Addr, payload []byte, index int) error {
	proto := uint32(payload[0])<<24 | uint32(payload[1])<<16 |
		uint32(payload[2])<<8 | uint32(payload[3])

	if proto != p.proto {
		return xudp.ErrDiscard
	}

	return nil
}
