// Package conntrack implements the Conntrack subsystem of the Netfilter (Netlink) protocol family.
// The package is intended to be clear, user-friendly, thoroughly tested and easy to understand.
//
// It is purely written in Go, without any dependency on Cgo or any C library, kernel headers
// or userspace tools.  It uses a native Netlink implementation (https://github.com/mdlayher/netlink)
// and does not parse or scrape any output of the `conntrack` command.
//
// It is designed in a way that makes the user acquainted with the structure of the protocol,
// with a clean separation between the Conntrack types/attributes and the Netfilter layer (implemented
// in https://github.com/ti-mo/netfilter).
//
// All Conntrack attributes known to the kernel up until version 4.17 are implemented. There is experimental
// support for manipulating Conntrack 'expectations', beside listening and dumping. The original focus of the
// package was receiving Conntrack events over Netlink multicast sockets, but was since expanded to be a full
// implementation supporting queries.
package conntrack
