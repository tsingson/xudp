// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package reliability

import (
	"github.com/jteeuwen/xudp"
	"net"
	"time"
)

type PacketFunc func(sequence uint32, payload []byte)

type Plugin struct {
	*Reliability
	onSent    PacketFunc // Notify the host when a specific packet is sent.
	onRecv    PacketFunc // Notify the host when a specific packet is received.
	frequency uint       // Polling frequency for reliability updates.
}

// New creates a new reliability plugin.
//
// The given handlers optionally notify the host of sent, lost and ACK'ed
// packets by their sequence number.
// 
// The sent/recv handlers tell the host what sequence number belongs to a given
// packet. This sequence can be used as a key in a map, keeping track
// of packet data. The lost and acked handlers only yield this sequence number.
//
// frequency denotes the number of times per second we should update the
// internal packet queues. This is necessary to account for packet timeouts
// and thus accurately determine when a packet is lost or not. 30 times per
// second is the usual value, but you can adjust it towhatever you want.
func New(sent, recv PacketFunc, acked, lost SequenceFunc, frequency uint) xudp.Plugin {
	p := new(Plugin)
	p.Reliability = NewReliability()
	p.onSent = sent
	p.onRecv = recv
	p.onAcked = acked
	p.onLost = lost
	p.frequency = frequency
	go p.poll()
	return p
}

func (c *Plugin) Open(port int) error { return nil }
func (c *Plugin) Close() error        { return nil }

// poll regularly calls update() on the reliability system.
// This keeps the ACK handling synchronized, regardless of how often
// we send/recv data.
func (p *Plugin) poll() {
	var prev, curr int64
	var delta float32

	tick := time.NewTicker(time.Second / time.Duration(p.frequency))

	for {
		select {
		case <-tick.C:
			curr = time.Now().UnixNano()
			delta = float32(curr-prev) / float32(time.Second)
			prev = curr

			p.update(delta)
		}
	}
}

func (p *Plugin) PayloadSize() int { return 12 }

func (p *Plugin) Send(addr net.Addr, payload []byte, index int) error {
	n := p.LocalSequence
	payload[0] = byte(n >> 24)
	payload[1] = byte(n >> 16)
	payload[2] = byte(n >> 8)
	payload[3] = byte(n)

	n = p.RemoteSequence
	payload[4] = byte(n >> 24)
	payload[5] = byte(n >> 16)
	payload[6] = byte(n >> 8)
	payload[7] = byte(n)

	n = p.ackVector()
	payload[8] = byte(n >> 24)
	payload[9] = byte(n >> 16)
	payload[10] = byte(n >> 8)
	payload[11] = byte(n)

	if p.onSent != nil {
		p.onSent(p.LocalSequence, payload[index:])
	}

	p.packetSent(uint32(len(payload[index:])))
	return nil
}

func (p *Plugin) Recv(addr net.Addr, payload []byte, index int) error {
	sequence := uint32(payload[0])<<24 | uint32(payload[1])<<16 |
		uint32(payload[2])<<8 | uint32(payload[3])

	ack := uint32(payload[4])<<24 | uint32(payload[5])<<16 |
		uint32(payload[6])<<8 | uint32(payload[7])

	vector := uint32(payload[8])<<24 | uint32(payload[9])<<16 |
		uint32(payload[10])<<8 | uint32(payload[11])

	p.packetRecv(sequence, ack, vector, uint32(len(payload[index:])))

	if p.onRecv != nil {
		p.onRecv(sequence, payload[index:])
	}

	return nil
}
