## Ident

This plugin implements a peer identification layer for an
`xudp.Connection`. It embeds a 32 byte SHA256 hash of the peer's local
IP and port into each packet. Combined with the peer's public IP address,
this allows us to accurately identify multiple peers sending data from the
same NAT router/firewall.

Accessing this identifier is done through a PeerFunc handler
we supply to the plugin. Whenever a new packet arrives, this handler
is called with the unique peer hash and the payload.

The hash is generated from the client's internal NAT address and port.
Combined with the public IP, this gives us a reliable key by which to tell
them apart.

The internal port is not strictly necessary for the outside world, but it
is needed when two peers from the same local IP talk to each other.
This happens when two clients are run on the same computer. The only thing
setting them apart is their local port number.

The hash is implemented as follows:

	private := SHA256(private_ip + private_port)
	hash := Base64( SHA256( public_ip + public_port + private ) )

The private part is included in every outgoing packet.


### Usage

    go get github.com/jteeuwen/xudp/plugins/ident


### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

