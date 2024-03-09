// Copyright 2023 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file ewcept in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !nowatchdog
// +build !nowatchdog

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

type testWatchdogCollector struct {
	wc Collector
}

func (c testWatchdogCollector) Collect(ch chan<- prometheus.Metric) {
	c.wc.Update(ch)
}

func (c testWatchdogCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func TestWatchdogStats(t *testing.T) {
	testcase := `# HELP node_watchdog_access_cs0 Value of /sys/class/watchdog/<watchdog>/access_cs0
	# TYPE node_watchdog_access_cs0 gauge
	node_watchdog_access_cs0{name="watchdog0"} 0
	# HELP node_watchdog_bootstatus Value of /sys/class/watchdog/<watchdog>/bootstatus
	# TYPE node_watchdog_bootstatus gauge
	node_watchdog_bootstatus{name="watchdog0"} 1
	# HELP node_watchdog_fw_version Value of /sys/class/watchdog/<watchdog>/fw_version
	# TYPE node_watchdog_fw_version gauge
	node_watchdog_fw_version{name="watchdog0"} 2
	# HELP node_watchdog_info Info of /sys/class/watchdog/<watchdog>
	# TYPE node_watchdog_info gauge
	node_watchdog_info{identity="",name="watchdog1",options="",pretimeout_governor="",state="",status=""} 1
	node_watchdog_info{identity="Software Watchdog",name="watchdog0",options="0x8380",pretimeout_governor="noop",state="active",status="0x8000"} 1
	# HELP node_watchdog_nowayout Value of /sys/class/watchdog/<watchdog>/nowayout
	# TYPE node_watchdog_nowayout gauge
	node_watchdog_nowayout{name="watchdog0"} 0
	# HELP node_watchdog_pretimeout_seconds Value of /sys/class/watchdog/<watchdog>/pretimeout
	# TYPE node_watchdog_pretimeout_seconds gauge
	node_watchdog_pretimeout_seconds{name="watchdog0"} 120
	# HELP node_watchdog_timeleft_seconds Value of /sys/class/watchdog/<watchdog>/timeleft
	# TYPE node_watchdog_timeleft_seconds gauge
	node_watchdog_timeleft_seconds{name="watchdog0"} 300
	# HELP node_watchdog_timeout_seconds Value of /sys/class/watchdog/<watchdog>/timeout
	# TYPE node_watchdog_timeout_seconds gauge
	node_watchdog_timeout_seconds{name="watchdog0"} 60
	`
	*sysPath = "fixtures/sys"

	logger := log.NewLogfmtLogger(os.Stderr)
	c, err := NewWatchdogCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	reg := prometheus.NewRegistry()
	reg.MustRegister(&testWatchdogCollector{wc: c})

	sink := make(chan prometheus.Metric)
	go func() {
		err = c.Update(sink)
		if err != nil {
			panic(fmt.Errorf("failed to update collector: %s", err))
		}
		close(sink)
	}()

	err = testutil.GatherAndCompare(reg, strings.NewReader(testcase))
	if err != nil {
		t.Fatal(err)
	}
}
