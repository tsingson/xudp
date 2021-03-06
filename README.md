## xudp

XUDP offers an API for extended UDP networking.

The basic `xudp.Connection` is nothing more than a wrapper around a
standard go UDP connection. Any useful features should be loaded into
the connection through any of the available plugins, or by writing your
own plugin.

A plugin has a hook into the `Send` and `Recv` methods of the connection.
Each send/received packet is passed through every registered plugin, which
can then perform any necessary operations. For most plugins, this means
appending certain plugin specific metadata to the packet payload in the form
of one or more header fields.

Writing your own plugin is simple. It need only implement the `xudp.Plugin`
interface. It can then be used by calling `Connection.Register()` with an
instance of this plugin. From this point on, the connection will
automatically use your plugin and no other code is necessary to make
it all work.

During `Send` or `Recv` calls, each plugin is passed the sender/receiver's
address and the packet payload, starting at the byte offset for that specific
plugin. It can then access the first byte of data simply at `payload[0]`.

While some plugins can be re-used by multiple connections, it is not
recommended to do so. Some plugins retain internal state on a per-connection
basis. Re-using the same instance in other connections will mess up the
internals. It is therefore advised to give each connection their own,
new instance of a given plugin.

Individual plugins may expose additional fields and methods, useful for
the host. These can be accessed by asserting the `xudp.Plugin` type to its
concrete implementation. Refer to each plugin's documentation for
information on what it has to offer.


### NAT Punch-through

This package does not handle NAT punch-through.

The `ident` plugin is not meant for this purpose. It merely identifies
multiple clients from behind the same NAT through their unique Peer Id.
It does not care whether or not a peer is actually behind a NAT. Nor does
identify the type of NAT.

There are existing protocols around for this purpose (e.g.: [STUN][stun]).
Apart from that, it involves a full server implementation which handles NAT
identification for connecting peers. It can not be implemented merely as
a plugin for the `xudp.Connection` type.

Many UDP based programs out there either use [STUN][stun] directly, or have
a custom implementation of the protocol specific to their application
needs.

It should be noted that NAT punch-through is only necessary when
any of the endpoints use IPv4. If the entire stack uses IPv6, NAT becomes
a non-issue and no special punch-through mechanisms have to be employed.

[stun]: http://tools.ietf.org/html/rfc5389


### Usage

    go get github.com/jteeuwen/xudp
    go get github.com/jteeuwen/xudp/plugin/<name>

Example for setup and use of a connection:

	conn := xudp.New(MTU)
	conn.Register(protocol.New(ProtocolId))
	...

Open the connection for incoming data:

	err := conn.Open(port)
	...
	defer conn.Close()

Sending & receiving data:

	for {
		addr, payload, err := conn.Recv()
		...
		err := conn.Send(addr, payload)
	}

### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

