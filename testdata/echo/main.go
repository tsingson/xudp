// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

/*
This application demonstrates a simple Peer to Peer echo loop.

We launch the program twice. Once with just a port. This will serve
as the 'server'. The second gets the target address of the server and
will function as a 'client'. It will initiate the echo loop.

	$ go run main.go -p 12345
	$ go run main.go -p 12346 -a :12345

From this point on, both will simply bounce the same packet payload
back and forth between them until one of the programs is stopped.

The speed of the packet transfer is limited to our hypothetical game loop.
It is set to 60 frames per second. This means we send/recv data at this
same rate. loop() contains timers which govern the progression of each frame.
Play with the ticker timeouts to increase or decrease the transfer speed.

It should be noted that high through put is not the primary goal of the
XUDP package. Reliability and fast access to time sensitive data is.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/jteeuwen/xudp"
	"net"
	"os"
	"time"
)

const (
	MTU        = 1400
	ProtocolId = 0xBADBEEF
)

var (
	address net.Addr
	port = flag.Int("p", 0, "The port number on which to listen on.")
	framerate = flag.Uint("fps", 60, "Frame rate for the game loop simulation.")
	payload = make([]byte, MTU-xudp.UDPHeaderSize-xudp.XUDPHeaderSize)
	units   = []byte{' ', 'K', 'M', 'G', 'T', 'P', 'Y'}
)

func main() {
	parseArgs()

	conn := initConn()
	defer conn.Close()

	if address != nil {
		// If we have a target address, we should
		// initiate the echo loop.
		conn.Send(address, xudp.NewPacket(payload))
	}

	loop(conn)
}

// parseArgs parses commandline arguments.
func parseArgs() {
	straddr := flag.String("a", "", "The server address to connect to. Only needed for client mode.")
	flag.Parse()

	if *port == 0 && len(*straddr) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if len(*straddr) > 0 {
		var err error
		address, err = net.ResolveUDPAddr("udp", *straddr)

		if err != nil {
			fmt.Fprintf(os.Stderr, "initConn: %v\n", err)
			os.Exit(1)
		}
	}
}

// initConn initializes our connection.
func initConn() *xudp.Connection {
	conn := xudp.NewConnection(MTU, ProtocolId)
	err := conn.Open(*port)

	if err != nil {
		fmt.Fprintf(os.Stderr, "initConn: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Listening on port %d...\n", *port)
	return conn
}

// The main 'game' loop.
func loop(c *xudp.Connection) {
	var sender net.Addr
	var packet xudp.Packet
	var err error

	delta := 1.0 / float32(*framerate)
	statTicker := time.NewTicker(time.Second)
	frameTicker := time.NewTicker(time.Second / time.Duration(*framerate))
	start := time.Now().Unix()

	for {
		select {
		case <- statTicker.C:
			stat(c, time.Now().Unix() - start)

		case <- frameTicker.C:
			c.Update(delta)
			sender, packet, err = c.Recv()

			if err != nil {
				fmt.Fprintf(os.Stderr, "Recv: %v\n", err)
				return
			}

			_, err = c.Send(sender, packet)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Send: %v\n", err)
				return
			}
		}
	}
}

// stat prints some connection statistics.
func stat(c *xudp.Connection, delta int64) {
	df := float64(1)
	
	if delta > 0 {
		df = float64(delta)
	}

	rp := c.RecvPackets
	rb := pretty(c.RecvBytes)
	rbps := prettyf(float64(c.RecvBytes)/df)

	sp := c.SentPackets
	sb := pretty(c.SentBytes)
	sbps := prettyf(float64(c.SentBytes)/df)

	fmt.Printf("in: %d @ %s - %s/s     out: %d @ %s - %s/s\n",
		rp, rb, rbps, sp, sb, sbps)
}

// pretty returns a 'pretty' version of the given byte size.
func pretty(b uint64) string {
	var u int

	for b >= 1024 {
		b /= 1024
		u++
	}

	return fmt.Sprintf("%d %cb", b, units[u])
}

// prettyf returns a 'pretty' version of the given byte size.
func prettyf(b float64) string {
	var u int

	for b >= 1024 {
		b /= 1024
		u++
	}

	return fmt.Sprintf("%3.2f %cb", b, units[u])
}
