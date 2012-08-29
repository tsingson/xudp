// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

import (
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

	jane.Send(bob.Addr, NewPacket(nil))

	for {
		select {
		case <-time.After(time.Microsecond * 10):
			bob.Close()
			jane.Close()
			return
		}
	}
}

func echo(t *testing.T, c *TestConn) {
	var sender net.Addr
	var packet Packet
	var err error

	for {
		sender, packet, err = c.Recv()
		if err != nil {
			return
		}

		c.Count++

		_, err = c.Send(sender, packet)

		if err != nil {
			t.Errorf("%s.Send: %v", c.Name, err)
			return
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
