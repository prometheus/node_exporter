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

// +build !nomeminfo

package collector

import (
	"encoding/json"
	"regexp"

	"github.com/shirou/gopsutil/mem"
)

var (
	reParens = regexp.MustCompile(`\((.*)\)`)
)

func parseMemInfo(v, s interface{}) error {
	byteMemInfo, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteMemInfo, &s)
	if err != nil {
		return err
	}

	return nil
}

func mergeInfo(infors, s map[string]float64, prefix string) error {
	// Merged values
	for k, v := range infors {
		// change fields name
		if k == "usedPercent" {
			k = prefix + "_used_percent"
		} else {
			k = prefix + "_" + k + "_bytes"
		}

		s[k] = v
	}
	return nil
}

func (c *meminfoCollector) getMemInfo() (map[string]float64, error) {
	info := make(map[string]float64)
	// Get swap memory stats
	swapInfo := make(map[string]float64)
	swapStats, err := mem.SwapMemory()
	if err != nil {
		return info, err
	}
	err = parseMemInfo(swapStats, &swapInfo)
	if err != nil {
		return info, err
	}

	// Get virtual memory stats
	virtualInfo := make(map[string]float64)
	virtualStats, err := mem.VirtualMemory()
	if err != nil {
		return info, err
	}
	err = parseMemInfo(virtualStats, &virtualInfo)
	if err != nil {
		return info, err
	}

	// Merged values
	mergeInfo(swapInfo, info, "swap")
	mergeInfo(virtualInfo, info, "virtual")

	return info, nil
}
