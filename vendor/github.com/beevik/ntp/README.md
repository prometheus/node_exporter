[![Build Status](https://travis-ci.org/beevik/ntp.svg?branch=master)](https://travis-ci.org/beevik/ntp)
[![GoDoc](https://godoc.org/github.com/beevik/ntp?status.svg)](https://godoc.org/github.com/beevik/ntp)

ntp
===

The ntp package is a small implementation of a limited NTP client. It
requests the current time from a remote NTP server according to
selected version of the NTP protocol. Client uses version 4 of the NTP
protocol RFC5905 by default.

The approach was inspired by a post to the go-nuts mailing list by
Michael Hofmann:

https://groups.google.com/forum/?fromgroups#!topic/golang-nuts/FlcdMU5fkLQ
