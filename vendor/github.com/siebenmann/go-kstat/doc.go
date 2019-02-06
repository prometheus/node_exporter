//
// Package kstat provides a Go interface to the Solaris/OmniOS
// kstat(s) system for user-level access to a lot of kernel
// statistics. For more documentation on kstats, see kstat(1) and
// kstat(3kstat).
//
// The package can retrieve what are called 'named' kstat statistics,
// IO statistics, and the most common additional types of 'raw'
// statistics, which covers almost all kstats you will normally find
// in the kernel. You can see the names and types of other kstats, but
// not currently retrieve data for them. Named statistics are the most
// common type for general information; IO statistics are exported by
// disks and some other things. Supported additional raw kstats are
// unix:0:sysinfo, unix:0:vminfo, unix:0:var, and mnt:*:mntinfo.
//
// General usage for named statistics: call Open() to obtain a Token,
// then call GetNamed() on it to obtain Named(s) for specific
// statistics. Note that this always gives you the very latest value
// for the statistic. If you want a number of statistics from the same
// module:inst:name triplet (eg several network counters from the same
// network interface) and you want them to all have been gathered at
// the same time, you need to call .Lookup() to obtain a KStat and
// then repeatedly call its .GetNamed() (this is also slightly more
// efficient).
//
// The short version: a kstat is a collection of some related
// statistics, eg various network counters for a particular network
// interface. A Token is a handle for a collection of kstats. You go
// collection (Token) -> kstat (KStat) -> specific statistic (Named)
// in order to retrieve the value of a specific statistic.
//
// (IO stats are retrieved all at once with GetIO(), because they come
// to us from the kernel as one single struct so that's what you get.)
//
// This is a cgo-based package. Cross compilation is up to you.
// Goroutine safety is in no way guaranteed because the underlying
// C kstat library is probably not thread or goroutine safe (and
// there are some all-Go concurrency races involving .Close()).
//
// This package may leak memory, especially since the Solaris kstat
// manpage is not clear on the requirements here. However I believe
// it's reasonably memory safe. It's possible to totally corrupt
// memory with use-after-free errors if you do operations on kstats
// after calling Token.Close(), although we try to avoid that.
//
// NOTE: this package is quite young. The API may well change as
// I (and other people) gain more experience with it.
//
// PERFORMANCE
//
// In general this is not going to be as lean and mean as calling
// C directly, partly because of intrinsic CGo overheads and partly
// because we do more memory allocation and deallocation than a C
// program would (partly because we prioritize not leaking memory).
//
// SUPPORTED AND UNSUPPORTED KSTAT TYPES
//
// We support named kstats and IO kstats (KSTAT_TYPE_NAMED and
// KSTAT_TYPE_IO / kstat_io_t respectively). kstat(1) also knows about
// a number of magic specific 'raw' stats (which are generally custom
// C structs); of these we support unix:0:sysinfo, unix:0:vminfo,
// unix:0:var, and mnt:*:mntinfo for NFS filesystem mounts.
//
// In theory kstat supports general timer and interrupt stats. In
// practice there is no use of KSTAT_TYPE_TIMER in the current Illumos
// kernel source and very little use of KSTAT_TYPE_INTR (mostly by
// very old hardware drivers, although the vioif driver uses it too).
// Since I can't test KSTAT_TYPE_INTR stats, we don't currently
// support it.
//
// There are also a few additional KSTAT_TYPE_RAW raw stats that we
// don't support, mostly because they seem to be effectively obsolete.
// These specific raw stats can be found listed in the Illumos source
// code in cmd/stat/kstat/kstat.h in the ks_raw_lookup array. See
// cmd/stat/kstat/kstat.c for how they're interpreted. If you need
// access to one of these kstats, the KStat.CopyTo() and KStat.Raw()
// methods give you an escape hatch to roll your own. You'll probably
// need to use cgo to generate an appropriate Go struct that matches
// the C struct you need. My notes on this process may be helpful:
//
// https://utcc.utoronto.ca/~cks/space/blog/programming/GoCGoCompatibleStructs
//
// Author: Chris Siebenmann
// https://github.com/siebenmann/go-kstat
//
// Copyright: standard Go copyright.
//
// (If you're reading this documentation on a non-Solaris platform,
// you're probably not seeing the detailed API documentation for
// constants, types, and so on because of tooling limitations in godoc
// et al.)
//
package kstat

//
// This exists in large part to give non-Solaris systems something
// that makes the package name visible. Otherwise eg goimports thinks
// that this is an empty package and deletes 'import ...' statements
// for it if you process a Go file that imports the package on a
// non-Solaris platform.
//
// Since this file exists, I've put the package-level documentation
// here to increase the odds that tools running on non-Solaris systems
// will be able to show you at least some documentation.
//
// This is a hack, and annoying.
