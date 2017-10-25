// Copyright 2017 The Prometheus Authors
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

// +build freebsd dragonfly
// +build !nomeminfo

package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
	"unsafe"
)

// #include <sys/types.h>
import "C"

type bsdSysctlType uint8

// BSD-specific sysctl value types.  There is an impedience mismatch between
// native C types, e.g. int vs long, and the golang unix.Sysctl variables
const (
	// Default to uint32.
	bsdSysctlTypeUint32 bsdSysctlType = iota
	bsdSysctlTypeUint64
	bsdSysctlTypeStructTimeval
	bsdSysctlTypeCLong
)

// Contains all the info needed to map a single bsd-sysctl to a prometheus
// value.
type bsdSysctl struct {
	// Prometheus name
	name string

	// Simple prometheus description
	description string

	// Prometheus type
	valueType prometheus.ValueType

	// Sysctl name
	mib string

	// Sysctl data-type
	dataType bsdSysctlType

	// Post-retrieval conversion hooks
	conversion func(float64) float64
}

func (b bsdSysctl) Value() (float64, error) {
	var tmp32 uint32
	var tmp64 uint64
	var tmpf64 float64
	var err error

	switch b.dataType {
	case bsdSysctlTypeUint32:
		tmp32, err = unix.SysctlUint32(b.mib)
		tmpf64 = float64(tmp32)
	case bsdSysctlTypeUint64:
		tmp64, err = unix.SysctlUint64(b.mib)
		tmpf64 = float64(tmp64)
	case bsdSysctlTypeStructTimeval:
		tmpf64, err = b.getStructTimeval()
	case bsdSysctlTypeCLong:
		tmpf64, err = b.getCLong()
	}

	if err != nil {
		return 0, err
	}

	if b.conversion != nil {
		return b.conversion(tmpf64), nil
	}

	return tmpf64, nil
}

func (b bsdSysctl) getStructTimeval() (float64, error) {
	raw, err := unix.SysctlRaw(b.mib)
	if err != nil {
		return 0, err
	}

	/*
	 * From 10.3-RELEASE sources:
	 *
	 * /usr/include/sys/_timeval.h:47
	 *  time_t      tv_sec
	 *  suseconds_t tv_usec
	 *
	 * /usr/include/sys/_types.h:60
	 *  long __suseconds_t
	 *
	 * ... architecture dependent, via #ifdef:
	 *  typedef __int64_t __time_t;
	 *  typedef __int32_t __time_t;
	 */
	if len(raw) != (C.sizeof_time_t + C.sizeof_suseconds_t) {
		// Shouldn't get here, unless the ABI changes...
		return 0, fmt.Errorf(
			"length of bytes received from sysctl (%d) does not match expected bytes (%d)",
			len(raw),
			C.sizeof_time_t+C.sizeof_suseconds_t,
		)
	}

	secondsUp := unsafe.Pointer(&raw[0])
	susecondsUp := uintptr(secondsUp) + C.sizeof_time_t
	unix := float64(*(*C.time_t)(secondsUp))
	usec := float64(*(*C.suseconds_t)(unsafe.Pointer(susecondsUp)))

	// This conversion maintains the usec precision.  Using the time
	// package did not.
	return (unix + (usec / float64(1000*1000))), nil
}

func (b bsdSysctl) getCLong() (float64, error) {
	raw, err := unix.SysctlRaw(b.mib)
	if err != nil {
		return 0, err
	}

	if len(raw) == C.sizeof_long {
		return float64(*(*C.long)(unsafe.Pointer(&raw[0]))), nil
	}

	if len(raw) == C.sizeof_int {
		// This is valid for at least vfs.bufspace, and the default
		// long handler - which can clamp longs to 32-bits:
		//   https://github.com/freebsd/freebsd/blob/releng/10.3/sys/kern/vfs_bio.c#L338
		//   https://github.com/freebsd/freebsd/blob/releng/10.3/sys/kern/kern_sysctl.c#L1062
		return float64(*(*C.int)(unsafe.Pointer(&raw[0]))), nil
	}

	return 0, fmt.Errorf(
		"length of bytes received from sysctl (%d) does not match expected bytes (long: %d), (int: %d)",
		len(raw),
		C.sizeof_long,
		C.sizeof_int,
	)

}
