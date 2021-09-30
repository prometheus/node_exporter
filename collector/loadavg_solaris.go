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

//go:build !noloadavg
// +build !noloadavg

package collector

import (
	"fmt"
	"strconv"

	"github.com/illumos/go-kstat"
)

// #include <sys/param.h>
import "C"

func kstatToFloat(ks *kstat.KStat, kstatKey string) float64 {
	kstatValue, err := ks.GetNamed(kstatKey)

	if err != nil {
		panic(err)
	}

	kstatLoadavg, err := strconv.ParseFloat(
		fmt.Sprintf("%.2f", float64(kstatValue.UintVal)/C.FSCALE), 64)

	if err != nil {
		panic(err)
	}

	return kstatLoadavg
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

	loadavg1Min := kstatToFloat(ks, "avenrun_1min")
	loadavg5Min := kstatToFloat(ks, "avenrun_5min")
	loadavg15Min := kstatToFloat(ks, "avenrun_15min")

	return []float64{loadavg1Min, loadavg5Min, loadavg15Min}, nil
}
