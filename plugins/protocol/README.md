## Protocol

The protocol plugin adds a single 32 bit, unsigned integer to the
packet. It represents our protocol identifier and is used to determine
of a given received packet should be processed by our application or not.

Any packet we send/receive contains this value as part of its header.
Any received packet which does not carry this exact value, is simply
discarded.

This protocol id can be any number you want, but it is advised to use
something relatively unique. A 4 byte hash of the name of your program
can be a suitable id.

For best results, we recommend that this is the first plugin you register.
This ensures it is executed first for every packet. If a packet does not
match the given protocol ID, we do not have to waste time and resources
on other plugins being run.


### Usage

    go get github.com/jteeuwen/xudp/plugins/protocol


### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.



