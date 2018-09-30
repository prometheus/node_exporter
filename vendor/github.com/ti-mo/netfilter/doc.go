// Package netfilter provides encoding and decoding of Netlink messages into Netfilter attributes.
// It handles Netfilter-specific nesting of attributes, endianness, and is written around a native
// Netlink implementation (https://github.com/mdlayher/netlink). It is purely written in Go,
// without any dependency on Cgo or any C library, kernel headers or userspace tools.
//
// The goal of this package is to be used for implementing the Netfilter family of Netlink protocols.
// For an example implementation, see https://github.com/ti-mo/conntrack.
package netfilter
