# Perf
[![GoDoc](https://godoc.org/github.com/hodgesds/perf-utils?status.svg)](https://godoc.org/github.com/hodgesds/perf-utils)

This package is a go library for interacting with the `perf` subsystem in
Linux. It allows you to do things like see how many CPU instructions a function
takes, profile a process for various hardware events, and other interesting
things. The library is by no means finalized and should be considered pre-alpha
at best.

# Use Cases
A majority of the utility methods in this package should only be used for
testing and/or debugging performance issues. Due to the nature of the go
runtime profiling on the goroutine level is extremely tricky, with the
exception of a long running worker goroutine locked to an OS thread. Eventually
this library could be used to implement many of the features of `perf` but in
accessible via Go directly.

## Caveats
* Some utility functions will call
  [`runtime.LockOSThread`](https://golang.org/pkg/runtime/#LockOSThread) for
  you, they will also unlock the thread after profiling. ***Note*** using these
  utility functions will incur significant overhead.
* Overflow handling is not implemented.

# Setup
Most likely you will need to tweak some system settings unless you are running as root. From `man perf_event_open`:

```
   perf_event related configuration files
       Files in /proc/sys/kernel/

           /proc/sys/kernel/perf_event_paranoid
                  The perf_event_paranoid file can be set to restrict access to the performance counters.

                  2   allow only user-space measurements (default since Linux 4.6).
                  1   allow both kernel and user measurements (default before Linux 4.6).
                  0   allow access to CPU-specific data but not raw tracepoint samples.
                  -1  no restrictions.

                  The existence of the perf_event_paranoid file is the official method for determining if a kernel supports perf_event_open().

           /proc/sys/kernel/perf_event_max_sample_rate
                  This sets the maximum sample rate.  Setting this too high can allow users to sample at a rate that impacts overall machine performance and potentially lock up the machine.  The default value is 100000  (samples  per
                  second).

           /proc/sys/kernel/perf_event_max_stack
                  This file sets the maximum depth of stack frame entries reported when generating a call trace.

           /proc/sys/kernel/perf_event_mlock_kb
                  Maximum number of pages an unprivileged user can mlock(2).  The default is 516 (kB).

```

# Example
Say you wanted to see how many CPU instructions a particular function took:

```
package main

import (
	"fmt"
	"log"
	"github.com/hodgesds/perf-utils"
)

func foo() error {
	var total int
	for i:=0;i<1000;i++ {
		total++
	}
	return nil
}

func main() {
	profileValue, err := perf.CPUInstructions(foo)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("CPU instructions: %+v\n", profileValue)
}
```

# Benchmarks
To profile a single function call there is an overhead of ~0.4ms.

```
$ go test  -bench=BenchmarkCPUCycles .
goos: linux
goarch: amd64
pkg: github.com/hodgesds/perf-utils
BenchmarkCPUCycles-8        3000            397924 ns/op              32 B/op          1 allocs/op
PASS
ok      github.com/hodgesds/perf-utils  1.255s
```

The `Profiler` interface has low overhead and suitable for many use cases:

```
$ go test  -bench=BenchmarkProfiler .
goos: linux
goarch: amd64
pkg: github.com/hodgesds/perf-utils
BenchmarkProfiler-8      3000000               488 ns/op              32 B/op          1 allocs/op
PASS
ok      github.com/hodgesds/perf-utils  1.981s
```

# BPF Support
BPF is supported by using the `BPFProfiler` which is available via the
`ProfileTracepoint` function. To use BPF you need to create the BPF program and
then call `AttachBPF` with the file descriptor of the BPF program. This is not
well tested so use at your own peril.

# Misc
Originally I set out to use `go generate` to build Go structs that were
compatible with perf, I found a really good
[article](https://utcc.utoronto.ca/~cks/space/blog/programming/GoCGoCompatibleStructs)
on how to do so. Eventually, after digging through some of the `/x/sys/unix`
code I found pretty much what I was needed. However, I think if you are
interested in interacting with the kernel it is a worthwhile read.
