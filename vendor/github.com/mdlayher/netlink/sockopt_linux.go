// +build linux,!386

package netlink

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// setsockopt provides access to the setsockopt syscall.
func setsockopt(fd, level, name int, v unsafe.Pointer, l uint32) error {
	_, _, errno := unix.Syscall6(
		unix.SYS_SETSOCKOPT,
		uintptr(fd),
		uintptr(level),
		uintptr(name),
		uintptr(v),
		uintptr(l),
		0,
	)
	if errno != 0 {
		return error(errno)
	}

	return nil
}
