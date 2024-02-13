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

//go:build (freebsd || dragonfly || openbsd || netbsd || darwin) && cgo
// +build freebsd dragonfly openbsd netbsd darwin
// +build cgo

package collector

import (
	"fmt"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
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
