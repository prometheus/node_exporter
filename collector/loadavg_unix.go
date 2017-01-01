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

// +build darwin dragonfly netbsd openbsd solaris
// +build !noloadavg

package collector

import (
	"errors"
	"runtime"
)

// #include <stdlib.h>
import "C"

func getLoad() ([]float64, error) {
	var loadavg [3]C.double
	samples := C.getloadavg(&loadavg[0], 3)
	numCpu := float64(runtime.NumCPU())
	if samples > 0 {
		return []float64{float64(loadavg[0]), float64(loadavg[1]), float64(loadavg[2]), float64(loadavg[0]) / numCpu, float64(loadavg[1]) / numCpu, float64(loadavg[2]) / numCpu}, nil
	} else {
		return nil, errors.New("failed to get load average")
	}
}
