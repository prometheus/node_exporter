// Copyright 2015 The Prometheus Authors
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

// +build !noloadavg

package collector

import (
	"syscall"
)

func getLoad() (loads []float64, err error) {
	const scale float64 = 65536 // LINUX_SYSINFO_LOADS_SCALE
	var sysinfo syscall.Sysinfo_t

	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return nil, err
	}

	return []float64{
		float64(sysinfo.Loads[0]) / scale,
		float64(sysinfo.Loads[1]) / scale,
		float64(sysinfo.Loads[2]) / scale,
	}, nil
}
