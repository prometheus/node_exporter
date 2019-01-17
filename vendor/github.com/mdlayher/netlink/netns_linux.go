//+build linux

package netlink

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// getThreadNetNS gets the network namespace file descriptor of the thread the current goroutine
// is running on. Make sure to call runtime.LockOSThread() before this so the goroutine does not
// get scheduled on another thread in the meantime.
func getThreadNetNS() (int, error) {
	file, err := os.Open(fmt.Sprintf("/proc/%d/task/%d/ns/net", unix.Getpid(), unix.Gettid()))
	if err != nil {
		return -1, err
	}
	return int(file.Fd()), nil
}

// setThreadNetNS sets the network namespace of the thread of the current goroutine to
// the namespace described by the user-provided file descriptor.
func setThreadNetNS(fd int) error {
	return unix.Setns(fd, unix.CLONE_NEWNET)
}
