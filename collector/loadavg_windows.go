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

// +build !noloadavg

package collector

import (
	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/load"
)

// Read loadavg
func getLoad() (loads []float64, err error) {
	data, err := load.Avg()
	if err != nil {
		return nil, err
	}
	loads, err = parseLoad(data)
	if err != nil {
		return nil, err
	}
	return loads, nil
}

// Parse loadavg and return 1m, 5m and 15m.
func parseLoad(data *load.AvgStat) (loads []float64, err error) {
	loads = make([]float64, 0)
	info := make(map[string]float64)
	loadavg, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("could not parse load: %s", err)
	}
	json.Unmarshal(loadavg, &info)
	for _, v := range info {
		loads = append(loads, v)
	}
	return loads, nil
}
