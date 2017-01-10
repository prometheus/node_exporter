netlink [![Build Status](https://travis-ci.org/mdlayher/netlink.svg?branch=master)](https://travis-ci.org/mdlayher/netlink) [![GoDoc](https://godoc.org/github.com/mdlayher/netlink?status.svg)](https://godoc.org/github.com/mdlayher/netlink) [![Go Report Card](https://goreportcard.com/badge/github.com/mdlayher/netlink)](https://goreportcard.com/report/github.com/mdlayher/netlink)
=======

Package `netlink` provides low-level access to Linux netlink sockets.
MIT Licensed.

Why?
----

A [number of netlink packages](https://godoc.org/?q=netlink) are already
available for Go, but I wasn't able to find one that aligned with what
I wanted in a netlink package:

- Simple, idiomatic API
- Well tested
- Well documented
- Makes use of Go best practices
- Doesn't need root to work

My goal for this package is to use it as a building block for the creation
of other netlink protocol family packages.
