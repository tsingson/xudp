// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

/*
This application demonstrates a simple Peer to Peer echo loop.

We launch the program twice. Once with just a port. This will serve
as the 'server'. The second gets the target address of the server and
will function as a 'client'. It will initiate the echo loop.

	$ go build
	$ ./echo -p 30000
	$ ./echo -p 30001 -a :30000

From this point on, both will simply bounce a random packet payload
back and forth between them until one of the programs is stopped.

The speed of the packet transfer is limited to our hypothetical game loop.
It is set to 30 frames per second. This means we send/recv data at this
same rate.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/jteeuwen/xudp"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	MTU         = 1400
	ProtocolId  = 0xBADBEEF
	FrameRate   = 30
	PayloadSize = MTU - xudp.UDPHeaderSize - xudp.XUDPHeaderSize
)

func main() {
	port, address := parseArgs()

	conn := initConn(port)
	defer conn.Close()

	loop(conn, address)
}

// parseArgs parses commandline arguments.
func parseArgs() (int, net.Addr) {
	port := flag.Int("p", 30000, "Port to listen on for connections.")
	addr := flag.String("a", "", "The server address to connect to. Only needed for client mode.")
	flag.Parse()

	if *port == 0 && len(*addr) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if len(*addr) > 0 {
		address, err := net.ResolveUDPAddr("udp", *addr)

		if err != nil {
			fmt.Fprintf(os.Stderr, "initConn: %v\n", err)
			os.Exit(1)
		}

		return *port, address
	}

	return *port, nil
}

// initConn initializes our connection.
func initConn(port int) *xudp.Connection {
	conn := xudp.NewConnection(MTU, ProtocolId)
	err := conn.Open(port)

	if err != nil {
		fmt.Fprintf(os.Stderr, "initConn: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Listening on port %d...\n", port)
	return conn
}

// The main 'game' loop.
func loop(c *xudp.Connection, address net.Addr) {
	var prevTime, currTime int64
	var delta float32
	var ok bool

	// Track average sent/ACK'ed bandwidth
	avgSent := make([]float32, 0, 100)
	avgAcked := make([]float32, 0, 100)

	// Frame progression ticker.
	frameTick := time.NewTicker(time.Second / FrameRate)

	// Statistics printing ticker.
	statTick := time.NewTicker(time.Second)

	// Channel which accepts incoming messages.
	// This allows us to do non-blocking reads.
	recv := readLoop(c)

	// If we have an address, we are the 'client' and should
	// initiate the echo loop.
	if address != nil {
		c.Send(address, []byte("Hello"))
	}

	for {
		select {
		case <-frameTick.C:
			// Compute new delta time.
			currTime = time.Now().UnixNano()
			delta = float32(currTime-prevTime) / float32(time.Second)
			prevTime = currTime

			// Receive next message.
			if address, ok = <-recv; !ok {
				break
			}

			// Send a random payload back.
			c.Send(address, make([]byte, rand.Int31n(PayloadSize)))
			c.Update(delta)

		case <-statTick.C:
			stat(c, &avgSent, &avgAcked)
		}
	}
}

// readLoop reads data from the connection and yields it through the
// returned channel. This allows us to make the read a non-blocking operation.
//
// In this particular program, we do not care about the actual payload.
// Just the sender's address.
func readLoop(c *xudp.Connection) <-chan net.Addr {
	ch := make(chan net.Addr)

	go func() {
		defer close(ch)

		for {
			address, _, err := c.Recv()

			if err != nil {
				return
			}

			ch <- address
		}
	}()

	return ch
}

// stat prints connection statistics.
func stat(c *xudp.Connection, sent, acked *[]float32) {
	rt := c.RTT
	sp := c.SentPackets
	ap := c.AckedPackets
	lp := c.LostPackets

	// Update list for average sent bandwidth
	if len(*sent) < cap(*sent) {
		*sent = append(*sent, c.SentBandwidth)
	} else {
		copy((*sent)[1:], *sent)
		(*sent)[0] = c.SentBandwidth
	}

	// Update list for average ACK'ed bandwidth
	if len(*acked) < cap(*acked) {
		*acked = append(*acked, c.AckedBandwidth)
	} else {
		copy((*acked)[1:], *acked)
		(*acked)[0] = c.AckedBandwidth
	}

	var lr float32

	if sp > 0 {
		lr = float32(lp) / float32(sp) * 100.0
	}

	fmt.Printf(
		"rtt %.1fms, sent %d (%.1fkbps), acked %d (%.1fkbps), lost %d (%.1f%%)\n",
		rt*1000.0, sp, avg(*sent), ap, avg(*acked), lp, lr)
}

// avg returns the average of all values in the given list.
func avg(list []float32) float32 {
	switch len(list) {
	case 0:
		return 0

	case 1:
		return list[0]

	default:
		var total float64

		for _, v := range list {
			total += float64(v)
		}

		return float32(total / float64(len(list)))
	}

	panic("unreachable")
}
