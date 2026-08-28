// Copyright The Prometheus Authors
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
	"fmt"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type testCPUFreqCollector struct {
	mc Collector
}

func (c testCPUFreqCollector) Collect(ch chan<- prometheus.Metric) {
	c.mc.Update(ch)
}

func (c testCPUFreqCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func newTestCPUFreqCollector(logger *slog.Logger) (prometheus.Collector, error) {
	mc, err := NewCPUFreqCollector(logger)
	if err != nil {
		return testCPUFreqCollector{}, err
	}
	return &testCPUFreqCollector{mc}, nil
}

const cpuFreqFixtureMetrics = `# HELP node_%[1]s_scaling_frequency_hertz Current scaled CPU thread frequency in hertz.
# TYPE node_%[1]s_scaling_frequency_hertz gauge
node_%[1]s_scaling_frequency_hertz{cpu="0"} 1.699981e+09
node_%[1]s_scaling_frequency_hertz{cpu="1"} 1.699981e+09
node_%[1]s_scaling_frequency_hertz{cpu="2"} 8e+06
node_%[1]s_scaling_frequency_hertz{cpu="3"} 8e+06
# HELP node_%[1]s_scaling_frequency_max_hertz Maximum scaled CPU thread frequency in hertz.
# TYPE node_%[1]s_scaling_frequency_max_hertz gauge
node_%[1]s_scaling_frequency_max_hertz{cpu="0"} 3.7e+09
node_%[1]s_scaling_frequency_max_hertz{cpu="1"} 3.7e+09
node_%[1]s_scaling_frequency_max_hertz{cpu="2"} 4.2e+09
node_%[1]s_scaling_frequency_max_hertz{cpu="3"} 4.2e+09
# HELP node_%[1]s_scaling_frequency_min_hertz Minimum scaled CPU thread frequency in hertz.
# TYPE node_%[1]s_scaling_frequency_min_hertz gauge
node_%[1]s_scaling_frequency_min_hertz{cpu="0"} 8e+08
node_%[1]s_scaling_frequency_min_hertz{cpu="1"} 8e+08
node_%[1]s_scaling_frequency_min_hertz{cpu="2"} 1e+06
node_%[1]s_scaling_frequency_min_hertz{cpu="3"} 1e+06
# HELP node_%[1]s_scaling_governor Current enabled CPU frequency governor.
# TYPE node_%[1]s_scaling_governor gauge
node_%[1]s_scaling_governor{cpu="0",governor="performance"} 0
node_%[1]s_scaling_governor{cpu="0",governor="powersave"} 1
node_%[1]s_scaling_governor{cpu="1",governor="performance"} 0
node_%[1]s_scaling_governor{cpu="1",governor="powersave"} 1
node_%[1]s_scaling_governor{cpu="2",governor="performance"} 0
node_%[1]s_scaling_governor{cpu="2",governor="powersave"} 1
node_%[1]s_scaling_governor{cpu="3",governor="performance"} 0
node_%[1]s_scaling_governor{cpu="3",governor="powersave"} 1
`

func TestCPUFreqMetrics(t *testing.T) {
	*sysPath = "fixtures/sys"

	for _, tc := range []struct {
		name          string
		cpufreqPrefix bool
		subsystem     string
	}{
		{name: "cpu prefix", cpufreqPrefix: false, subsystem: "cpu"},
		{name: "cpufreq prefix", cpufreqPrefix: true, subsystem: "cpufreq"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			*useCPUFreqPrefix = tc.cpufreqPrefix
			defer func() { *useCPUFreqPrefix = false }()

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			collector, err := newTestCPUFreqCollector(logger)
			if err != nil {
				t.Fatal(err)
			}

			expected := fmt.Sprintf(cpuFreqFixtureMetrics, tc.subsystem)
			if err := testutil.CollectAndCompare(collector, strings.NewReader(expected)); err != nil {
				t.Fatal(err)
			}
		})
	}
}
