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

package collector

import (
	"os"
	"testing"
)

func TestMemInfoNuma(t *testing.T) {
	file, err := os.Open("fixtures/sys/devices/system/node/node0/meminfo")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	memInfo, err := parseMemInfoNuma(file)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := 707915776.0, memInfo[5].value; want != got {
		t.Errorf("want memory Active(anon) value %f, got %f", want, got)
	}

	if want, got := "Active_anon", memInfo[5].metricName; want != got {
		t.Errorf("want metric Active(anon) metricName %s, got %s", want, got)
	}

	if want, got := 150994944.0, memInfo[25].value; want != got {
		t.Errorf("want memory AnonHugePages %f, got %f", want, got)
	}

	file, err = os.Open("fixtures/sys/devices/system/node/node1/meminfo")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	memInfo, err = parseMemInfoNuma(file)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := 291930112.0, memInfo[6].value; want != got {
		t.Errorf("want memory Inactive(anon) %f, got %f", want, got)
	}

	if want, got := 85585088512.0, memInfo[13].value; want != got {
		t.Errorf("want memory FilePages %f, got %f", want, got)
	}
}

func TestMemInfoNumaStat(t *testing.T) {
	file, err := os.Open("fixtures/sys/devices/system/node/node0/numastat")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	numaStat, err := parseMemInfoNumaStat(file, "0")
	if err != nil {
		t.Fatal(err)
	}

	if want, got := 193460335812.0, numaStat[0].value; want != got {
		t.Errorf("want numa stat numa_hit value %f, got %f", want, got)
	}

	if want, got := "numa_hit_total", numaStat[0].metricName; want != got {
		t.Errorf("want numa stat numa_hit metricName %s, got %s", want, got)
	}

	if want, got := 193454780853.0, numaStat[4].value; want != got {
		t.Errorf("want numa stat local_node %f, got %f", want, got)
	}

	file, err = os.Open("fixtures/sys/devices/system/node/node1/numastat")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	numaStat, err = parseMemInfoNumaStat(file, "1")
	if err != nil {
		t.Fatal(err)
	}

	if want, got := 59858626709.0, numaStat[1].value; want != got {
		t.Errorf("want numa stat numa_miss %f, got %f", want, got)
	}

	if want, got := 59860526920.0, numaStat[5].value; want != got {
		t.Errorf("want numa stat other_node %f, got %f", want, got)
	}
}
