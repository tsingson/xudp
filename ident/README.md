## Ident

This package adds a peer identification layer on top of a basic
xudp.Connection. It embeds a 32 byte SHA256 hash of the peer's local
IP and port into each packet. This allows us to accurately identify
multiple peers sending data from the same NAT router/firewall.


### Usage

    go get github.com/jteeuwen/xudp/ident


### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

