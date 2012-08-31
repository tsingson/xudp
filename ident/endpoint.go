// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package ident

import (
	"crypto/sha256"
	"encoding/base64"
	"net"
)

// An endpoint identifies a peer. It contains their public address
// along with a hash of their internal NAT address (IP + port).
//
// For UDP traffic this is the only reliable way to keep clients behind a
// NAT apart from each other. We can not rely solely on their public
// address + port, because some NAT routers randomly cycle through ports
// for every outgoing message. This means that two distinct messages from the
// same machine may appear to be from a different source within the NAT system.
//
// The hash is generated from the client's internal NAT address and port.
// Combined with the public IP, this gives us a reliable key by which to tell
// them apart.
//
// The internal port is not strictly necessary for the outside world, but it
// is needed when two peers from the same local IP talk to each other.
// This happens when two clients are run on the same computer. The only thing
// setting them apart is their local port number.
type Endpoint struct {
	// Public IP and port for the peer. This is where we send our data
	// and it may change from one request to the next.
	Addr net.Addr

	// An SHA256 hash of the peer's internal IP address and port.
	// This is included in every packet from the given peer.
	id PeerId
}

// NewEndpoint creates a new endpoint for the given public address
// and peer id. Leave the id empty if you do not have it yet.
func NewEndpoint(addr net.Addr, id PeerId) *Endpoint {
	e := new(Endpoint)
	e.Addr = addr
	e.id = id
	return e
}

// String returns a base64 encoded SHA256 hash of the public IP address
// combined with the endpoint Id. This can be used as a reliable
// identification key for a given peer.
func (e *Endpoint) String() string {
	ua, ok := e.Addr.(*net.UDPAddr)

	if !ok {
		return ""
	}

	hm := sha256.New()
	hm.Write(ua.IP.To16())
	hm.Write(e.id)

	return base64.StdEncoding.EncodeToString(hm.Sum(nil))
}

// Equals returns true if the two endpoints represent the same peer.
// The comparison is done in constant time.
func (e *Endpoint) Equals(dest *Endpoint) bool {
	return compareHash(
		[]byte(e.String()),
		[]byte(dest.String()),
	)
}
