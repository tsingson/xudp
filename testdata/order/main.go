// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	port, address := parseArgs()

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
func initConn(port int) *OrderedConnection {
	conn := NewOrderedConnection(1400)
	err := conn.Open(port)

	if err != nil {
		fmt.Fprintf(os.Stderr, "initConn: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Listening on port %d...\n", port)
	return conn
}

// The main loop.
func loop(c *OrderedConnection, addr net.Addr) {
	var count int
	var err error

	// Initiate echo loop with first message.
	if addr != nil {
		c.Send(addr, []byte(fmt.Sprintf("Message #%d", count)))
	}

	tick := time.NewTicker(time.Millisecond * 100)

	for {
		select {
		case <-tick.C:
			if addr != nil {
				// Client keeps sending new data packages.
				data := []byte(fmt.Sprintf("Message #%d", count))
				count++

				err = c.Send(addr, data)

				if err != nil {
					return
				}

				log.Printf("send: %s", data)
			}

		case data := <-c.Incoming:
			if data == nil {
				return
			}

			if addr == nil {
				// Server only returns empty packets.
				// These are necessary because they contain our ACK data.
				log.Printf("recv: %s", data.Payload)
				err = c.Send(data.Addr, nil)
			}

			if err != nil {
				return
			}
		}
	}
}
