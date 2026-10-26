// Copyright 2024 The Prometheus Authors
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

package numatopology

import (
	"strconv"
	"strings"
)

// CountCPUList counts the CPUs described by a sysfs cpulist string.
// Examples: "0" → 1, "0-3" → 4, "0-3,8-11" → 8.
func CountCPUList(cpulist string) int {
	total := 0
	for _, part := range strings.Split(cpulist, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if idx := strings.IndexByte(part, '-'); idx != -1 {
			lo, err1 := strconv.Atoi(part[:idx])
			hi, err2 := strconv.Atoi(part[idx+1:])
			if err1 == nil && err2 == nil && hi >= lo {
				total += hi - lo + 1
			}
		} else {
			if _, err := strconv.Atoi(part); err == nil {
				total++
			}
		}
	}
	return total
}
