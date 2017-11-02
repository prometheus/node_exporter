[![Build Status](https://travis-ci.org/beevik/ntp.svg?branch=master)](https://travis-ci.org/beevik/ntp)
[![GoDoc](https://godoc.org/github.com/beevik/ntp?status.svg)](https://godoc.org/github.com/beevik/ntp)

ntp
===

The ntp package is an implementation of a Simple NTP (SNTP) client based on
[RFC5905](https://tools.ietf.org/html/rfc5905). It allows you to connect to
a remote NTP server and request information about the current time.


## Querying the current time

If all you care about is the current time according to a remote NTP server,
simply use the `Time` function:
```go
time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
```


## Querying time metadata

To obtain the current time as well as some additional metadata about the time,
use the `Query` function:
```go
response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
time := time.Now().Add(response.ClockOffset)
```

Alternatively, use the `QueryWithOptions` function if you want to change the
default behavior used by the `Query` function:
```go
options := ntp.QueryOptions{ Timeout: 30*time.Second, TTL: 5 }
response, err := ntp.QueryWithOptions("0.beevik-ntp.pool.ntp.org", options)
time := time.Now().Add(response.ClockOffset)
```

The `Response` structure returned by `Query` includes the following
information:
* `Time`: The time the server transmitted its response, according to its own clock.
* `ClockOffset`: The estimated offset of the local system clock relative to the server's clock. For a more accurate time reading, you may add this offset to any subsequent system clock reading.
* `RTT`: An estimate of the round-trip-time delay between the client and the server.
* `Precision`: The precision of the server's clock reading.
* `Stratum`: The server's stratum, which indicates the number of hops from the server to the reference clock. A stratum 1 server is directly attached to the reference clock. If the stratum is zero, the server has responded with the "kiss of death".
* `ReferenceID`: A unique identifier for the consulted reference clock.
* `ReferenceTime`: The time at which the server last updated its local clock setting.
* `RootDelay`: The server's aggregate round-trip-time delay to the stratum 1 server.
* `RootDispersion`: The server's estimated maximum measurement error relative to the reference clock.
* `RootDistance`: An estimate of the root synchronization distance between the client and the stratum 1 server.
* `Leap`: The leap second indicator, indicating whether a second should be added to or removed from the current month's last minute.
* `MinError`: A lower bound on the clock error between the client and the server.
* `KissCode`: A 4-character string describing the reason for a "kiss of death" response (stratum=0).
* `Poll`: The maximum polling interval between successive messages to the server.

The `Response` structure's `Validate` method performs additional sanity checks
to determine whether the response is suitable for time synchronization
purposes.
```go
err := response.Validate()
if err == nil {
    // response data is suitable for synchronization purposes
}
```

## Using the NTP pool

The NTP pool is a shared resource used by people all over the world.
To prevent it from becoming overloaded, please avoid querying the standard
`pool.ntp.org` zone names in your applications.  Instead, consider requesting
your own [vendor zone](http://www.pool.ntp.org/en/vendors.html) or [joining
the pool](http://www.pool.ntp.org/join.html).
