Release v0.3.0
==============

There have been no breaking changes or further deprecations since the
previous release.

**Changes**

* Fixed a bug in the calculation of NTP timestamps.

Release v0.2.0
==============

There are no breaking changes or further deprecations in this release.

**Changes**

* Added `KissCode` to the `Response` structure.


Release v0.1.1
==============

**Breaking changes**

* Removed the `MaxStratum` constant.

**Deprecations**

* Officially deprecated the `TimeV` function.

**Internal changes**

* Removed `minDispersion` from the `RootDistance` calculation, since the value
  was arbitrary.
* Moved some validation into main code path so that invalid `TransmitTime` and
  `mode` responses trigger an error even when `Response.Validate` is not
  called.


Release v0.1.0
==============

This is the initial release of the `ntp` package.  Currently it supports the following features:
* `Time()` to query the current time according to a remote NTP server.
* `Query()` to query multiple pieces of time-related information from a remote NTP server.
* `QueryWithOptions()`, which is like `Query()` but with the ability to override default query options.

Time-related information returned by the `Query` functions includes:
* `Time`: the time the server transmitted its response, according to the server's clock.
* `ClockOffset`: the estimated offset of the client's clock relative to the server's clock. You may apply this offset to any local system clock reading once the query is complete.
* `RTT`: an estimate of the round-trip-time delay between the client and the server.
* `Precision`: the precision of the server's clock reading.
* `Stratum`: the "stratum" level of the server, where 1 indicates a server directly connected to a reference clock, and values greater than 1 indicating the number of hops from the reference clock.
* `ReferenceID`: A unique identifier for the NTP server that was contacted.
* `ReferenceTime`: The time at which the server last updated its local clock setting.
* `RootDelay`: The server's round-trip delay to the reference clock.
* `RootDispersion`: The server's total dispersion to the referenced clock.
* `RootDistance`: An estimate of the root synchronization distance.
* `Leap`: The leap second indicator.
* `MinError`: A lower bound on the clock error between the client and the server.
* `Poll`: the maximum polling interval between successive messages on the server.

The `Response` structure returned by the `Query` functions also contains a `Response.Validate()` function that returns an error if any of the fields returned by the server are invalid.
