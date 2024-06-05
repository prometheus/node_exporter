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

//go:build !nomeminfo_procfs
// +build !nomeminfo_procfs

package collector

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

const (
	testMetrics = `# HELP node_memory_bytes Value in bytes for the labeled field in /proc/meminfo.
# TYPE node_memory_bytes gauge
node_memory_bytes{field="ActiveAnon"} 2.068484096e+09
node_memory_bytes{field="Active"} 2.287017984e+09
node_memory_bytes{field="ActiveFile"} 2.18533888e+08
node_memory_bytes{field="AnonHugePages"} 0
node_memory_bytes{field="AnonPages"} 2.298032128e+09
node_memory_bytes{field="Bounce"} 0
node_memory_bytes{field="Buffers"} 2.256896e+07
node_memory_bytes{field="Cached"} 9.53229312e+08
node_memory_bytes{field="CmaFree"} 0
node_memory_bytes{field="CmaTotal"} 0
node_memory_bytes{field="CommitLimit"} 6.210940928e+09
node_memory_bytes{field="CommittedAS"} 8.023486464e+09
node_memory_bytes{field="DirectMap1G"} 0
node_memory_bytes{field="DirectMap2M"} 3.787456512e+09
node_memory_bytes{field="DirectMap4k"} 1.9011584e+08
node_memory_bytes{field="Dirty"} 1.077248e+06
node_memory_bytes{field="HardwareCorrupted"} 0
node_memory_bytes{field="Hugepagesize"} 2.097152e+06
node_memory_bytes{field="InactiveAnon"} 9.04245248e+08
node_memory_bytes{field="Inactive"} 1.053417472e+09
node_memory_bytes{field="InactiveFile"} 1.49172224e+08
node_memory_bytes{field="KernelStack"} 5.9392e+06
node_memory_bytes{field="Mapped"} 2.4496128e+08
node_memory_bytes{field="MemAvailable"} 0
node_memory_bytes{field="MemFree"} 2.30883328e+08
node_memory_bytes{field="MemTotal"} 3.831959552e+09
node_memory_bytes{field="Mlocked"} 32768
node_memory_bytes{field="NFSUnstable"} 0
node_memory_bytes{field="PageTables"} 7.7017088e+07
node_memory_bytes{field="Percpu"} 0
node_memory_bytes{field="SReclaimable"} 4.5846528e+07
node_memory_bytes{field="SUnreclaim"} 5.545984e+07
node_memory_bytes{field="Shmem"} 6.0809216e+08
node_memory_bytes{field="ShmemHugePages"} 0
node_memory_bytes{field="ShmemPmdMapped"} 0
node_memory_bytes{field="Slab"} 1.01306368e+08
node_memory_bytes{field="SwapCached"} 1.97124096e+08
node_memory_bytes{field="SwapFree"} 3.23108864e+09
node_memory_bytes{field="SwapTotal"} 4.2949632e+09
node_memory_bytes{field="Unevictable"} 32768
node_memory_bytes{field="VmallocChunk"} 3.5183963009024e+13
node_memory_bytes{field="VmallocTotal"} 3.5184372087808e+13
node_memory_bytes{field="VmallocUsed"} 3.6130816e+08
node_memory_bytes{field="Writeback"} 0
node_memory_bytes{field="WritebackTmp"} 0
`
)

type testMeminfoProcfsCollector struct {
	mc Collector
}

func (c testMeminfoProcfsCollector) Collect(ch chan<- prometheus.Metric) {
	c.mc.Update(ch)
}

func (c testMeminfoProcfsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func NewTestMeminfoProcfsCollector(logger log.Logger) (prometheus.Collector, error) {
	mc, err := NewMeminfoProcfsCollector(logger)
	if err != nil {
		return testMeminfoProcfsCollector{}, err
	}
	return testMeminfoProcfsCollector{
		mc: mc,
	}, err
}

func TestMemInfoProcfs(t *testing.T) {
	*procPath = "fixtures/proc"
	logger := log.NewLogfmtLogger(os.Stderr)

	collector, err := NewMeminfoProcfsCollector(logger)
	if err != nil {
		panic(err)
	}
	c, err := NewTestMeminfoProcfsCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	sink := make(chan prometheus.Metric)
	go func() {
		err = collector.Update(sink)
		if err != nil {
			panic(fmt.Errorf("failed to update collector: %s", err))
		}
		close(sink)
	}()

	err = testutil.GatherAndCompare(reg, strings.NewReader(testMetrics))
	if err != nil {
		t.Fatal(err)
	}
}
