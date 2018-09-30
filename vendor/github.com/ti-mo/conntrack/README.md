# conntrack [![GoDoc](https://godoc.org/github.com/ti-mo/conntrack?status.svg)](https://godoc.org/github.com/ti-mo/conntrack) [![Build Status](https://semaphoreci.com/api/v1/ti-mo/conntrack/branches/master/shields_badge.svg)](https://semaphoreci.com/ti-mo/conntrack) [![Coverage Status](https://coveralls.io/repos/github/ti-mo/conntrack/badge.svg?branch=master)](https://coveralls.io/github/ti-mo/conntrack?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/ti-mo/conntrack)](https://goreportcard.com/report/github.com/ti-mo/conntrack)

Package `conntrack` implements the Conntrack subsystem of the Netfilter (Netlink) protocol family.
The package is intended to be clear, user-friendly, thoroughly tested and easy to understand.

It is purely written in Go, without any dependency on Cgo or any C library, kernel headers
or userspace tools.  It uses a native Netlink implementation (https://github.com/mdlayher/netlink)
and does not parse or scrape any output of the `conntrack` command.

It is designed in a way that makes the user acquainted with the structure of the protocol,
with a clean separation between the Conntrack types/attributes and the Netfilter layer (implemented
in https://github.com/ti-mo/netfilter).

All Conntrack attributes known to the kernel up until version 4.17 are implemented. There is experimental
support for manipulating Conntrack 'expectations', beside listening and dumping. The original focus of the
package was receiving Conntrack events over Netlink multicast sockets, but was since expanded to be a full
implementation supporting queries.

## Features

With this library, the user can:

- Interact with conntrack connections and expectations through Flow and Expect types respectively
- Create, get, update and delete Flows in an idiomatic way (and Expects, to an extent)
- Listen for create/update/destroy events
- Flush (empty) and dump (display) the whole conntrack table, optionally filtering on specific connection marks

There are many usage examples in the [godoc](https://godoc.org/github.com/ti-mo/conntrack).

## Contributing

Contributions are absolutely welcome! Before starting work on large changes, please create an issue first,
or join #networking on Gophers Slack to discuss the design.

If you encounter a problem implementing the library, please open a GitHub issue for help.
