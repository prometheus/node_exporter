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
	"fmt"
	"regexp"
	"strconv"
)

var reMemTotal = regexp.MustCompile(`MemTotal:\s+(\d+)\s+kB`)

// ParseMeminfo extracts MemTotal from sysfs node meminfo content and returns bytes.
// The sysfs format is "Node N MemTotal:   NNNNNN kB"; the regex matches the
// "MemTotal: N kB" suffix regardless of any leading "Node N" prefix.
func ParseMeminfo(content string) (int64, error) {
	m := reMemTotal.FindStringSubmatch(content)
	if m == nil {
		return 0, fmt.Errorf("MemTotal not found in meminfo")
	}
	kib, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing MemTotal: %w", err)
	}
	return kib * 1024, nil
}
