// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
	"fmt"
	"net"
	"testing"
	"time"
)

const (
	MTU        = 1400
	ProtocolId = 0xBADBEEF
)

type TestConn struct {
	*Connection
	Addr  *net.UDPAddr
	Name  string
	Count uint64
}

var (
	bob = &TestConn{
		Connection: NewConnection(MTU, ProtocolId),
		Name:       "bob",
		Addr:       &net.UDPAddr{Port: 12345},
	}

	jane = &TestConn{
		Connection: NewConnection(MTU, ProtocolId),
		Name:       "jane",
		Addr:       &net.UDPAddr{Port: 12346},
	}
)

func TestConnection(t *testing.T) {
	initConnections(t)

	go echo(t, bob)
	go echo(t, jane)

	jane.Send(bob.Addr, NewPacket([]byte("Hi")))

	for {
		select {
		case <-time.After(time.Second * 5):
			bob.Close()
			jane.Close()
			stat(bob)
			stat(jane)
			return
		}
	}
}

func stat(c *TestConn) {
	fmt.Printf("Sent: %d, Recv: %d, ACK: %d, Lost: %d, RX: %f, AX: %f, RTT: %f\n",
		c.SentPackets, c.RecvPackets, c.AckedPackets, c.LostPackets,
		c.SentBandwidth, c.AckedBandwidth, c.RTT)
}

func echo(t *testing.T, c *TestConn) {
	const delta = 1.0 / 30.0

	var sender net.Addr
	var packet Packet
	var err error

	tick := time.NewTicker(time.Second / 30)

	for {
		select {
		case <-tick.C:
			c.Update(delta)

			sender, packet, err = c.Recv()
			if err != nil {
				return
			}

			c.Count++

			_, err = c.Send(sender, packet)

			if err != nil {
				return
			}
		}
	}
}

func initConnections(t *testing.T) {
	err := bob.Open(bob.Addr.Port)

	if err != nil {
		t.Fatal(err)
	}

	err = jane.Open(jane.Addr.Port)

	if err != nil {
		bob.Close()
		t.Fatal(err)
	}
}
