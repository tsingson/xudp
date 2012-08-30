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
	Addr *net.UDPAddr
}

var (
	bob = &TestConn{
		Connection: NewConnection(MTU, ProtocolId),
		Addr:       &net.UDPAddr{Port: 12345},
	}

	jane = &TestConn{
		Connection: NewConnection(MTU, ProtocolId),
		Addr:       &net.UDPAddr{Port: 12346},
	}
)

func TestConnection(t *testing.T) {
	initConnections(t)

	go echo(t, bob)
	go echo(t, jane)

	jane.Send(bob.Addr, []byte("Hello, World"))

	for {
		select {
		case <-time.After(time.Second / 2):
			bob.Close()
			jane.Close()
			return
		}
	}
}

func echo(t *testing.T, c *TestConn) {
	var prevTime, currTime int64
	var sender net.Addr
	var payload []byte
	var delta float32
	var err error

	for {
		currTime = time.Now().UnixNano()
		delta = float32(currTime-prevTime) / float32(time.Second)
		prevTime = currTime

		sender, payload, err = c.Recv()
		if err != nil {
			return
		}

		err = c.Send(sender, payload)
		if err != nil {
			return
		}

		c.Update(delta)
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
