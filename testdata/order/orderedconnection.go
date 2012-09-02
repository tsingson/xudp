// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"github.com/jteeuwen/xudp"
	"github.com/jteeuwen/xudp/plugins/reliability"
	"net"
)

type Data struct {
	Addr    net.Addr
	Payload []byte
}

// An OrderedConnection guarantees in-order reception of data
// and resends lost packets.
type OrderedConnection struct {
	*xudp.Connection
	recvCache map[uint32]*Data
	sendCache map[uint32]*Data
	Incoming  chan *Data
	sequence  uint32
}

// NewOrderedConnection creates a new ordered connection.
func NewOrderedConnection(mtu uint32) *OrderedConnection {
	c := new(OrderedConnection)
	c.Connection = xudp.New(mtu)
	c.recvCache = make(map[uint32]*Data)
	c.sendCache = make(map[uint32]*Data)
	c.Incoming = make(chan *Data)

	c.Register(reliability.New(
		func(seq uint32, addr net.Addr, data []byte) {
			c.sent(seq, addr, data)
		},
		func(seq uint32, addr net.Addr, data []byte) {
			c.recv(seq, addr, data)
		},
		func(seq uint32) { c.acked(seq) },
		func(seq uint32) { c.lost(seq) },
		30,
	))

	return c
}

func (c *OrderedConnection) Open(port int) error {
	err := c.Connection.Open(port)

	if err != nil {
		return err
	}

	go func() {
		for {
			_, _, err := c.Recv()

			if err != nil {
				return
			}
		}
	}()

	return nil
}

func (c *OrderedConnection) sent(seq uint32, addr net.Addr, data []byte) {
	c.sendCache[seq] = &Data{addr, data}
}

func (c *OrderedConnection) recv(seq uint32, addr net.Addr, data []byte) {
	if seq == c.sequence {
		c.Incoming <- &Data{addr, data}
		c.sequence++
		return
	}

	c.recvCache[seq] = &Data{addr, data}

	d, ok := c.recvCache[c.sequence]

	if !ok {
		return
	}

	delete(c.recvCache, c.sequence)
	c.Incoming <- d
	c.sequence++
}

func (c *OrderedConnection) acked(seq uint32) {
	delete(c.sendCache, seq)
}

func (c *OrderedConnection) lost(seq uint32) {
	println("lost", seq)

	data, ok := c.sendCache[seq]

	if !ok {
		return
	}

	delete(c.sendCache, seq)
	c.Send(data.Addr, data.Payload)
}
