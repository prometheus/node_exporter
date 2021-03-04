# CHANGELOG

## Unreleased

- n/a

## v1.2.1

- [Bug Fix]
  [commit](https://github.com/mdlayher/netlink/commit/d81418f81b0bfa2465f33790a85624c63d6afe3d):
  `netlink.SetBPF` will no longer panic if an empty BPF filter is set.
- [Improvement]
  [commit](https://github.com/mdlayher/netlink/commit/8014f9a7dbf4fd7b84a1783dd7b470db9113ff36):
  the library now uses https://github.com/josharian/native to provide the
  system's native endianness at compile time, rather than re-computing it many
  times at runtime.

## v1.2.0

**This is the first release of package netlink that only supports Go 1.12+. Users on older versions must use v1.1.1.**

- [Improvement] [#173](https://github.com/mdlayher/netlink/pull/173): support
  for Go 1.11 and below has been dropped. All users are highly recommended to
  use a stable and supported release of Go for their applications.
- [Performance] [#171](https://github.com/mdlayher/netlink/pull/171):
  `netlink.Conn` no longer requires a locked OS thread for the vast majority of
  operations, which should result in a significant speedup for highly concurrent
  callers. Thanks @ti-mo.
- [Bug Fix] [#169](https://github.com/mdlayher/netlink/pull/169): calls to
  `netlink.Conn.Close` are now able to unblock concurrent calls to
  `netlink.Conn.Receive` and other blocking operations.

## v1.1.1

**This is the last release of package netlink that supports Go 1.11.**

- [Improvement] [#165](https://github.com/mdlayher/netlink/pull/165):
  `netlink.Conn` `SetReadBuffer` and `SetWriteBuffer` methods now attempt the
  `SO_*BUFFORCE` socket options to possibly ignore system limits given elevated
  caller permissions. Thanks @MarkusBauer.
- [Note]
  [commit](https://github.com/mdlayher/netlink/commit/c5f8ab79aa345dcfcf7f14d746659ca1b80a0ecc):
  `netlink.Conn.Close` has had a long-standing bug
  [#162](https://github.com/mdlayher/netlink/pull/162) related to internal
  concurrency handling where a call to `Close` is not sufficient to unblock
  pending reads. To effectively fix this issue, it is necessary to drop support
  for Go 1.11 and below. This will be fixed in a future release, but a
  workaround is noted in the method documentation as of now.

## v1.1.0

- [New API] [#157](https://github.com/mdlayher/netlink/pull/157): the
  `netlink.AttributeDecoder.TypeFlags` method enables retrieval of the type bits
  stored in a netlink attribute's type field, because the existing `Type` method
  masks away these bits. Thanks @ti-mo!
- [Performance] [#157](https://github.com/mdlayher/netlink/pull/157): `netlink.AttributeDecoder`
  now decodes netlink attributes on demand, enabling callers who only need a
  limited number of attributes to exit early from decoding loops. Thanks @ti-mo!
- [Improvement] [#161](https://github.com/mdlayher/netlink/pull/161): `netlink.Conn`
  system calls are now ready for Go 1.14+'s changes to goroutine preemption.
  See the PR for details.

## v1.0.0

- Initial stable commit.
