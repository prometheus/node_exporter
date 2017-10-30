[![Build Status](https://travis-ci.org/beevik/ntp.svg?branch=master)](https://travis-ci.org/beevik/ntp)
[![GoDoc](https://godoc.org/github.com/beevik/ntp?status.svg)](https://godoc.org/github.com/beevik/ntp)

ntp
===

The ntp package is an implementation of a Simple NTP (SNTP) client based on
[RFC5905](https://tools.ietf.org/html/rfc5905). It allows you to connect to
a remote NTP server and request the current time.

If all you care about is the current time according to a known remote NTP
server, simply use the `Time` function:
```go
time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
```

If you want the time as well as additional metadata about the time, use the
`Query` function instead:
```go
response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
```

To use the NTP pool in your application, please request your own
[vendor zone](http://www.pool.ntp.org/en/vendors.html).  Avoid using 
the `[number].pool.ntp.org` zone names in your applications.
