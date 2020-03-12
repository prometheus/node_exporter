// +build linux

package perf

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	// DebugFS is the filesystem type for debugfs.
	DebugFS = "debugfs"

	// TraceFS is the filesystem type for tracefs.
	TraceFS = "tracefs"

	// ProcMounts is the mount point for file systems in procfs.
	ProcMounts = "/proc/mounts"

	// PerfMaxStack is the mount point for the max perf event size.
	PerfMaxStack = "/proc/sys/kernel/perf_event_max_stack"

	// PerfMaxContexts is a sysfs mount that contains the max perf contexts.
	PerfMaxContexts = "/proc/sys/kernel/perf_event_max_contexts_per_stack"

	// SyscallsDir is a constant of the default tracing event syscalls directory.
	SyscallsDir = "/sys/kernel/debug/tracing/events/syscalls/"

	// TracingDir is a constant of the default tracing directory.
	TracingDir = "/sys/kernel/debug/tracing"
)

var (
	// ErrNoMount is when there is no such mount.
	ErrNoMount = fmt.Errorf("no such mount")
)

// TraceFSMount returns the first found mount point of a tracefs file system.
func TraceFSMount() (string, error) {
	mounts, err := GetFSMount(TraceFS)
	if err != nil {
		return "", err
	}
	if len(mounts) == 0 {
		return "", ErrNoMount
	}
	return mounts[0], nil
}

// DebugFSMount returns the first found mount point of a debugfs file system.
func DebugFSMount() (string, error) {
	mounts, err := GetFSMount(DebugFS)
	if err != nil {
		return "", err
	}
	if len(mounts) == 0 {
		return "", ErrNoMount
	}
	return mounts[0], nil
}

// GetFSMount is a helper function to get a mount file system type.
func GetFSMount(mountType string) ([]string, error) {
	mounts := []string{}
	file, err := os.Open(ProcMounts)
	if err != nil {
		return mounts, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		mountInfo := strings.Split(scanner.Text(), " ")
		if len(mountInfo) > 3 && mountInfo[2] == mountType {
			mounts = append(mounts, mountInfo[1])
		}
	}
	if err := scanner.Err(); err != nil {
		return mounts, err
	}

	return mounts, file.Close()
}

// fileToStrings is a helper method that reads a line line by line and returns
// a slice of strings.
func fileToStrings(path string) ([]string, error) {
	res := []string{}
	f, err := os.Open(path)
	if err != nil {
		return res, err
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return res, err
	}

	return res, nil
}
