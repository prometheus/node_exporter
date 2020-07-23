// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,386

package netlink

import (
	"syscall"
	"unsafe"
)

const (
	sysSETSOCKOPT = 0xe
)

func socketcall(call int, a0, a1, a2, a3, a4, a5 uintptr) (int, syscall.Errno)

// setsockopt provides access to the setsockopt syscall.
func setsockopt(fd, level, name int, v unsafe.Pointer, l uint32) error {
	_, errno := socketcall(
		sysSETSOCKOPT,
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
