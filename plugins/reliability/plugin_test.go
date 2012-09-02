// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package reliability

import (
	"github.com/jteeuwen/xudp"
	"net"
	"testing"
	"time"
)

var Payload = []byte("Hello, world!")

func TestConn(t *testing.T) {
	ca := initConn(t, 10011)
	cb := initConn(t, 10012)

	defer ca.Close()
	defer cb.Close()

	ca.Send(&net.UDPAddr{Port: 10012}, Payload)

	go loop(t, ca)
	go loop(t, cb)

	<-time.After(time.Second / 2)
}

func loop(t *testing.T, c *xudp.Connection) {
	tick := time.NewTicker(time.Second / 30)

	for {
		select {
		case <-tick.C:
			addr, payload, err := c.Recv()

			if err != nil {
				return
			}

			if len(payload) != len(Payload) {
				t.Errorf("Payload size mismatch: Want %d, have %d",
					len(Payload), len(payload))
				return
			}

			for i := range payload {
				if Payload[i] != payload[i] {
					t.Errorf("Payload mismatch at %d: Want %d, have %d",
						i, Payload[i] != payload[i])
					return
				}
			}

			err = c.Send(addr, payload)

			if err != nil {
				return
			}
		}
	}
}

func initConn(t *testing.T, port int) *xudp.Connection {
	c := xudp.New(1400)
	c.Register(New(sent, recv, acked, lost, 60))

	err := c.Open(port)

	if err != nil {
		t.Fatal(err)
	}

	return c
}

func sent(seq uint32, addr net.Addr, payload []byte) {
	//println("sent", seq)
}

func recv(seq uint32, addr net.Addr, payload []byte) {
	//println("recv", seq)
}

func lost(seq uint32) {
	//println("lost", seq)
}

func acked(seq uint32) {
	//println("ACK'ed", seq)
}
