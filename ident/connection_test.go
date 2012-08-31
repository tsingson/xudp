// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package ident

import (
	"net"
	"testing"
	"time"
)

const (
	MTU        = 1400
	ProtocolId = 0xBADBEEF
)

var (
	bob     = New(MTU, ProtocolId, 12345)
	jane    = New(MTU, ProtocolId, 12346)
	bobAddr = &net.UDPAddr{Port: bob.port}
)

func TestConnection(t *testing.T) {
	initConnections(t)

	go echo(t, bob)
	go echo(t, jane)

	jane.Send(bobAddr, []byte("Hello, World"))

	for {
		select {
		case <-time.After(time.Second / 2):
			bob.Close()
			jane.Close()
			return
		}
	}
}

func echo(t *testing.T, c *Connection) {
	var prevTime, currTime int64
	var sender *Endpoint
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

		err = c.Send(sender.Addr, payload)
		if err != nil {
			return
		}

		c.Update(delta)
	}
}

func initConnections(t *testing.T) {
	err := bob.Open()

	if err != nil {
		t.Fatal(err)
	}

	err = jane.Open()

	if err != nil {
		bob.Close()
		t.Fatal(err)
	}
}
