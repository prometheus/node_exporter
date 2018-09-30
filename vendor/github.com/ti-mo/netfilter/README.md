# netfilter [![GoDoc](https://godoc.org/github.com/ti-mo/netfilter?status.svg)](https://godoc.org/github.com/ti-mo/netfilter) [![Build Status](https://semaphoreci.com/api/v1/ti-mo/netfilter/branches/master/shields_badge.svg)](https://semaphoreci.com/ti-mo/netfilter) [![Coverage Status](https://coveralls.io/repos/github/ti-mo/netfilter/badge.svg?branch=master)](https://coveralls.io/github/ti-mo/netfilter?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/ti-mo/netfilter)](https://goreportcard.com/report/github.com/ti-mo/netfilter)

Package netfilter provides encoding and decoding of Netlink messages into Netfilter attributes.
It handles Netfilter-specific nesting of attributes, endianness, and is written around a native
Netlink implementation (https://github.com/mdlayher/netlink). It is purely written in Go,
without any dependency on Cgo or any C library, kernel headers or userspace tools.

The goal of this package is to be used for implementing the Netfilter family of Netlink protocols.
For an example implementation, see https://github.com/ti-mo/conntrack.

## Contributing

Contributions are absolutely welcome! Before starting work on large changes, please create an issue first, or join #networking on Gophers Slack to discuss the design.

If you encounter a problem implementing the library, please open a GitHub issue for help.
