[![Build Status](https://travis-ci.org/beevik/ntp.svg?branch=master)](https://travis-ci.org/beevik/ntp)
[![GoDoc](https://godoc.org/github.com/beevik/ntp?status.svg)](https://godoc.org/github.com/beevik/ntp)

ntp
===

The ntp package is an implementation of a simple NTP client. It allows you
to connect to a remote NTP server and request the current time.

To request the current time, simply do the following:
```go
time, err := ntp.Time("pool.ntp.org")
```

To request the current time along with additional metadata, use the Query
function:
```go
response, err := ntp.Query("pool.ntp.org")
```
