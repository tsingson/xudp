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
back and forth between them for the duration of one minute.

The -cpu flag lets us write CPU profile data to the given file.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/jteeuwen/xudp"
	"github.com/jteeuwen/xudp/plugins/reliability"
	"math/rand"
	"net"
	"os"
	"runtime/pprof"
	"time"
)

var (
	plugin *reliability.Plugin
	cpu = flag.String("cpu", "", "File name to write CPU profile to.")
)

func main() {
	port, address := parseArgs()

	if len(*cpu) > 0 {
		fd, err := os.Create(*cpu)

		if err != nil {
			fmt.Fprintf(os.Stderr, "os.Create: %v\n", err)
			os.Exit(1)
		}

		defer fd.Close()

		pprof.StartCPUProfile(fd)
		defer pprof.StopCPUProfile()
	}

	conn := initConn(port)
	defer conn.Close()

	go loop(conn, address)
	<-time.After(time.Minute)
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
	plugin = reliability.New(nil, nil, nil, nil, 30).(*reliability.Plugin)

	conn := xudp.New(1400)
	conn.Register(plugin)
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
	// Track average sent/ACK'ed bandwidth
	avgSent := make([]float32, 0, 100)
	avgAcked := make([]float32, 0, 100)

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
		case <-statTick.C:
			stat(&avgSent, &avgAcked)

		default:
			// Receive next message.
			address, ok := <-recv

			if !ok {
				break
			}

			// Send a random payload back.
			err := c.Send(address, make([]byte, rand.Int31n(int32(c.PayloadSize()))))

			if err != nil {
				return
			}
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
func stat(sent, acked *[]float32) {
	rt := plugin.RTT
	sp := plugin.SentPackets
	ap := plugin.AckedPackets
	lp := plugin.LostPackets

	// Update list for average sent bandwidth
	if len(*sent) < cap(*sent) {
		*sent = append(*sent, plugin.SentBandwidth)
	} else {
		copy((*sent)[1:], *sent)
		(*sent)[0] = plugin.SentBandwidth
	}

	// Update list for average ACK'ed bandwidth
	if len(*acked) < cap(*acked) {
		*acked = append(*acked, plugin.AckedBandwidth)
	} else {
		copy((*acked)[1:], *acked)
		(*acked)[0] = plugin.AckedBandwidth
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
