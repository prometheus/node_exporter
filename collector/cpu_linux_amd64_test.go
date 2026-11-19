// Copyright 2026 The Prometheus Authors
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

//go:build !nocpu

package collector

import (
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// The cpuinfo fixture is x86-specific, and procfs selects its parser by GOARCH.
func TestCPUInfoMetrics(t *testing.T) {
	oldProcPath := *procPath
	oldSysPath := *sysPath
	oldEnableCPUInfo := *enableCPUInfo
	oldFlagsInclude := *flagsInclude
	oldBugsInclude := *bugsInclude
	t.Cleanup(func() {
		*procPath = oldProcPath
		*sysPath = oldSysPath
		*enableCPUInfo = oldEnableCPUInfo
		*flagsInclude = oldFlagsInclude
		*bugsInclude = oldBugsInclude
	})

	*procPath = "fixtures/proc"
	*sysPath = "fixtures/sys"
	*enableCPUInfo = false
	*flagsInclude = "^(aes|avx.?|constant_tsc)$"
	*bugsInclude = "^(cpu_meltdown|spectre_.*|mds)$"

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	collector, err := NewCPUCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	if !*enableCPUInfo {
		t.Fatal("CPU info was not enabled by its include flags")
	}

	const expected = `# HELP node_cpu_bug_info The ` + "`bugs`" + ` field of CPU information from /proc/cpuinfo taken from the first core.
# TYPE node_cpu_bug_info gauge
node_cpu_bug_info{bug="cpu_meltdown"} 1
node_cpu_bug_info{bug="mds"} 1
node_cpu_bug_info{bug="spectre_v1"} 1
node_cpu_bug_info{bug="spectre_v2"} 1
# HELP node_cpu_flag_info The ` + "`flags`" + ` field of CPU information from /proc/cpuinfo taken from the first core.
# TYPE node_cpu_flag_info gauge
node_cpu_flag_info{flag="aes"} 1
node_cpu_flag_info{flag="avx"} 1
node_cpu_flag_info{flag="avx2"} 1
node_cpu_flag_info{flag="constant_tsc"} 1
# HELP node_cpu_info CPU information from /proc/cpuinfo.
# TYPE node_cpu_info gauge
node_cpu_info{cachesize="8192 KB",core="0",cpu="0",family="6",microcode="0xb4",model="142",model_name="Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz",package="0",stepping="10",vendor="GenuineIntel"} 1
node_cpu_info{cachesize="8192 KB",core="0",cpu="4",family="6",microcode="0xb4",model="142",model_name="Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz",package="0",stepping="10",vendor="GenuineIntel"} 1
node_cpu_info{cachesize="8192 KB",core="1",cpu="1",family="6",microcode="0xb4",model="142",model_name="Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz",package="0",stepping="10",vendor="GenuineIntel"} 1
node_cpu_info{cachesize="8192 KB",core="1",cpu="5",family="6",microcode="0xb4",model="142",model_name="Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz",package="0",stepping="10",vendor="GenuineIntel"} 1
node_cpu_info{cachesize="8192 KB",core="2",cpu="2",family="6",microcode="0xb4",model="142",model_name="Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz",package="0",stepping="10",vendor="GenuineIntel"} 1
node_cpu_info{cachesize="8192 KB",core="2",cpu="6",family="6",microcode="0xb4",model="142",model_name="Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz",package="0",stepping="10",vendor="GenuineIntel"} 1
node_cpu_info{cachesize="8192 KB",core="3",cpu="3",family="6",microcode="0xb4",model="142",model_name="Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz",package="0",stepping="10",vendor="GenuineIntel"} 1
node_cpu_info{cachesize="8192 KB",core="3",cpu="7",family="6",microcode="0xb4",model="142",model_name="Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz",package="0",stepping="10",vendor="GenuineIntel"} 1
# HELP node_scrape_collector_success node_exporter: Whether a collector succeeded.
# TYPE node_scrape_collector_success gauge
node_scrape_collector_success{collector="cpu"} 1
`

	registry := prometheus.NewRegistry()
	registry.MustRegister(&NodeCollector{
		Collectors: map[string]Collector{"cpu": collector},
		logger:     logger,
	})
	if err := testutil.GatherAndCompare(
		registry,
		strings.NewReader(expected),
		"node_cpu_bug_info",
		"node_cpu_flag_info",
		"node_cpu_info",
		"node_scrape_collector_success",
	); err != nil {
		t.Fatal(err)
	}
}
