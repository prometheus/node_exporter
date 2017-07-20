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
	"errors"
)

/*
// Define "__stack_chk_fail" and "__stack_chk_guard" symbols.
#cgo LDFLAGS: -fno-stack-protector -lssp
// Ensure "hrtime_t" is defined for sys/loadavg.h
#include <sys/time.h>
#include <sys/loadavg.h>
*/
import "C"

func getLoad() ([]float64, error) {
	var loadavg [3]C.double
	samples := C.getloadavg(&loadavg[0], 3)
	if samples != 3 {
		return nil, errors.New("failed to get load average")
	}
	return []float64{float64(loadavg[0]), float64(loadavg[1]), float64(loadavg[2])}, nil
}
