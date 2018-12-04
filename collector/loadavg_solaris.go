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
	"fmt"
	"strconv"

	"github.com/siebenmann/go-kstat"
)

// #include <sys/param.h>
import "C"

func kstatToFloat(ks *kstat.KStat, kstat_key string) float64 {
	kstat_value, err := ks.GetNamed(kstat_key)

	if err != nil {
		panic(err)
	}

	kstat_loadavg, err := strconv.ParseFloat(
		fmt.Sprintf("%.2f", float64(kstat_value.UintVal)/C.FSCALE), 64)

	if err != nil {
		panic(err)
	}

	return kstat_loadavg
}

func getLoad() ([]float64, error) {
	tok, err := kstat.Open()
	if err != nil {
		panic(err)
	}

	defer tok.Close()

	ks, err := tok.Lookup("unix", 0, "system_misc")

	if err != nil {
		panic(err)
	}

	loadavg_1min := kstatToFloat(ks, "avenrun_1min")
	loadavg_5min := kstatToFloat(ks, "avenrun_5min")
	loadavg_15min := kstatToFloat(ks, "avenrun_15min")

	return []float64{loadavg_1min, loadavg_5min, loadavg_15min}, nil
}
