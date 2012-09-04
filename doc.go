// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

/*
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
the host. These can be accessed by assrting the `xudp.Plugin` type to its
concrete implementation type. Refer to each plugin's documentation for
information on this.

Example for setup and use of a connection:

	conn := xudp.New(MTU)
	conn.Register(protocol.New(...))
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
*/
package xudp
