// Copyright 2016 The Prometheus Authors
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

//go:build (darwin || dragonfly || freebsd || netbsd || openbsd) && !noloadavg
// +build darwin dragonfly freebsd netbsd openbsd
// +build !noloadavg

package collector

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func getLoad() ([]float64, error) {
	type loadavg struct {
		load  [3]uint32
		scale int
	}
	b, err := unix.SysctlRaw("vm.loadavg")
	if err != nil {
		return nil, err
	}
	load := *(*loadavg)(unsafe.Pointer((&b[0])))
	scale := float64(load.scale)
	return []float64{
		float64(load.load[0]) / scale,
		float64(load.load[1]) / scale,
		float64(load.load[2]) / scale,
	}, nil
}
