// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package ident

import (
	"crypto/sha256"
	"encoding/base64"
	"net"
)

// Size of a single Peer id in bytes.
const PeerHashSize = 32

// PeerHash represents a unique peer.
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
//
// The hash is implemented as follows:
//
//    private := SHA256(private_ip + private_port)
//    hash := Base64( SHA256( public_ip + public_port + private ) )
//
// The private part is included in every outgoing packet.
type PeerHash string

// NewPeerHash returns a base64 encoded SHA256 hash of the public IP address
// combined with the supplied id. This can be used as a reliable
// identification key for a given peer.
func NewPeerHash(addr net.Addr, id []byte) PeerHash {
	if addr == nil || id == nil {
		return ""
	}

	ua, ok := addr.(*net.UDPAddr)

	if !ok {
		return PeerHash("")
	}

	hm := sha256.New()
	hm.Write(ua.IP.To16())
	hm.Write(id)

	return PeerHash(base64.StdEncoding.EncodeToString(hm.Sum(nil)))
}

// Equals returns true if the two hashes represent the same peer.
//
// A constant time comparison is used to prevent timing attacks from
// being performed. With a normal bytes.Equal(a, b) comparison, an attacker can 
// time how long this function takes to complete. The longer it takes
// to return, the more of the hash he knows will be correct. A constant time
// comparison always runs in the same amount of time, regardless of the
// hash contents; thus eliminating the timing attack vector.
func (a PeerHash) Equals(b PeerHash) bool {
	if len(a) != len(b) {
		return false
	}

	var v byte

	for i := 0; i < len(a); i++ {
		v |= a[i] ^ b[i]
	}

	v = ^v
	v &= v >> 4
	v &= v >> 2
	v &= v >> 1
	return v == 1
}
