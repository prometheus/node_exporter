// +build linux

package perf

import (
	"golang.org/x/sys/unix"
)

// BPFProfiler is a Profiler that allows attaching a Berkeley
// Packet Filter (BPF) program to an existing kprobe tracepoint event.
// You need CAP_SYS_ADMIN privileges to use this interface. See:
// https://lwn.net/Articles/683504/
type BPFProfiler interface {
	Profiler
	AttachBPF(int) error
}

// AttachBPF is used to attach a BPF program to a profiler by using the file
// descriptor of the BPF program.
func (p *profiler) AttachBPF(fd int) error {
	return unix.IoctlSetInt(p.fd, unix.PERF_EVENT_IOC_SET_BPF, fd)
}
