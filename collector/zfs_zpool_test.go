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

package collector

import (
	"os"
	"testing"
)

func TestZpoolParsing(t *testing.T) {
	zpoolOutput, err := os.Open("fixtures/zfs/zpool_stats_stdout.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer zpoolOutput.Close()

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	pools := make([]string, 2)
	troutSize := float64(-1)
	troutDedupratio := float64(-1)
	zrootCapacity := float64(-1)

	err = c.parseZpoolOutput(zpoolOutput, func(pool, name string, value float64) {
		pools = append(pools, pool)
		if pool == "trout" && name == "size" {
			troutSize = value
		}
		if pool == "trout" && name == "dedupratio" {
			troutDedupratio = value
		}
		if pool == "zroot" && name == "capacity" {
			zrootCapacity = value
		}
	})

	if err != nil {
		t.Fatal(err)
	}

	if pools[0] == "trout" && pools[1] == "zroot" {
		t.Fatal("Did not parse all pools in fixture")
	}

	if troutSize != float64(4294967296) {
		t.Fatal("Unexpected value for pool 'trout's size value")
	}

	if troutDedupratio != float64(1.0) {
		t.Fatal("Unexpected value for pool 'trout's dedupratio value")
	}

	if zrootCapacity != float64(0.5) {
		t.Fatal("Unexpected value for pool 'zroot's capacity value")
	}
}
