// Copyright 2020 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"golang.org/x/sys/unix"
	"syscall"
	"unsafe"
)

func int8ToString(a []int8) string {
	buf := make([]byte, len(a))
	for i, v := range a {
		if byte(v) == 0 {
			buf = buf[:i]
			break
		}
		buf[i] = byte(v)
	}
	return string(buf)
}

// unix._C_int
type _C_int int32

var _zero uintptr

func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case unix.EAGAIN:
		return syscall.EAGAIN
	case unix.EINVAL:
		return syscall.EINVAL
	case unix.ENOENT:
		return syscall.ENOENT
	}
	return e
}

func _sysctl(mib []_C_int, old *byte, oldlen *uintptr, new *byte, newlen uintptr) (err error) {
	var _p0 unsafe.Pointer
	if len(mib) > 0 {
		_p0 = unsafe.Pointer(&mib[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	for {
		_, _, e1 := unix.Syscall6(unix.SYS___SYSCTL, uintptr(_p0), uintptr(len(mib)), uintptr(unsafe.Pointer(old)), uintptr(unsafe.Pointer(oldlen)), uintptr(unsafe.Pointer(new)), uintptr(newlen))
		if e1 != 0 {
			err = errnoErr(e1)
		}
		if err != unix.EINTR {
			return
		}
	}
	return
}

func sysctl(mib []_C_int) ([]byte, error) {
	n := uintptr(0)
	if err := _sysctl(mib, nil, &n, nil, 0); err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}

	buf := make([]byte, n)
	if err := _sysctl(mib, &buf[0], &n, nil, 0); err != nil {
		return nil, err
	}
	return buf[:n], nil
}
