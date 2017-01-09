// +build linux,!386

package netlink

import (
	"syscall"
	"unsafe"
)

// setsockopt provides access to the setsockopt syscall.
func setsockopt(fd, level, name int, v unsafe.Pointer, l uint32) error {
	_, _, errno := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
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
