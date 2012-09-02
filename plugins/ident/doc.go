// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

/*
This plugin implements a peer identification layer for an
`xudp.Connection`. It embeds a 32 byte SHA256 hash of the peer's local
IP and port into each packet. Combined with the peer's public IP address,
this allows us to accurately identify multiple peers sending data from the
same NAT router/firewall.

Accessing this identifier is done through a PeerFunc handler
we supply to the plugin. Whenever a new packet arrives, this handler
is called with the unique peer hash and the payload.
*/
package ident
