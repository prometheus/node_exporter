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

//go:build !nomdadm
// +build !nomdadm

package collector

import (
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type testMdadmCollector struct {
	mc Collector
}

func (c testMdadmCollector) Collect(ch chan<- prometheus.Metric) {
	c.mc.Update(ch)
}

func (c testMdadmCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func NewTestMdadmCollector(logger *slog.Logger) (prometheus.Collector, error) {
	mc, err := NewMdadmCollector(logger)
	if err != nil {
		return testMdadmCollector{}, err
	}
	return &testMdadmCollector{mc}, nil
}

func TestMdadmStats(t *testing.T) {
	*sysPath = "fixtures/sys"
	*procPath = "fixtures/proc"
	testcase := `# HELP node_md_blocks Total number of blocks on device.
        # TYPE node_md_blocks gauge
        node_md_blocks{device="md0"} 248896
        node_md_blocks{device="md00"} 4.186624e+06
        node_md_blocks{device="md10"} 3.14159265e+08
        node_md_blocks{device="md101"} 322560
        node_md_blocks{device="md11"} 4.190208e+06
        node_md_blocks{device="md12"} 3.886394368e+09
        node_md_blocks{device="md120"} 2.095104e+06
        node_md_blocks{device="md126"} 1.855870976e+09
        node_md_blocks{device="md127"} 3.12319552e+08
        node_md_blocks{device="md201"} 1.993728e+06
        node_md_blocks{device="md219"} 7932
        node_md_blocks{device="md3"} 5.853468288e+09
        node_md_blocks{device="md4"} 4.883648e+06
        node_md_blocks{device="md6"} 1.95310144e+08
        node_md_blocks{device="md7"} 7.813735424e+09
        node_md_blocks{device="md8"} 1.95310144e+08
        node_md_blocks{device="md9"} 523968
        # HELP node_md_blocks_synced Number of blocks synced on device.
        # TYPE node_md_blocks_synced gauge
        node_md_blocks_synced{device="md0"} 248896
        node_md_blocks_synced{device="md00"} 4.186624e+06
        node_md_blocks_synced{device="md10"} 3.14159265e+08
        node_md_blocks_synced{device="md101"} 322560
        node_md_blocks_synced{device="md11"} 0
        node_md_blocks_synced{device="md12"} 3.886394368e+09
        node_md_blocks_synced{device="md120"} 2.095104e+06
        node_md_blocks_synced{device="md126"} 1.855870976e+09
        node_md_blocks_synced{device="md127"} 3.12319552e+08
        node_md_blocks_synced{device="md201"} 114176
        node_md_blocks_synced{device="md219"} 7932
        node_md_blocks_synced{device="md3"} 5.853468288e+09
        node_md_blocks_synced{device="md4"} 4.883648e+06
        node_md_blocks_synced{device="md6"} 1.6775552e+07
        node_md_blocks_synced{device="md7"} 7.813735424e+09
        node_md_blocks_synced{device="md8"} 1.6775552e+07
        node_md_blocks_synced{device="md9"} 0
        # HELP node_md_degraded Number of degraded disks on device.
        # TYPE node_md_degraded gauge
        node_md_degraded{device="md0"} 0
        node_md_degraded{device="md1"} 0
        node_md_degraded{device="md10"} 0
        node_md_degraded{device="md4"} 0
        node_md_degraded{device="md5"} 1
        node_md_degraded{device="md6"} 1
        # HELP node_md_disks Number of active/failed/spare disks of device.
        # TYPE node_md_disks gauge
        node_md_disks{device="md0",state="active"} 2
        node_md_disks{device="md0",state="failed"} 0
        node_md_disks{device="md0",state="spare"} 0
        node_md_disks{device="md00",state="active"} 1
        node_md_disks{device="md00",state="failed"} 0
        node_md_disks{device="md00",state="spare"} 0
        node_md_disks{device="md10",state="active"} 2
        node_md_disks{device="md10",state="failed"} 0
        node_md_disks{device="md10",state="spare"} 0
        node_md_disks{device="md101",state="active"} 3
        node_md_disks{device="md101",state="failed"} 0
        node_md_disks{device="md101",state="spare"} 0
        node_md_disks{device="md11",state="active"} 2
        node_md_disks{device="md11",state="failed"} 1
        node_md_disks{device="md11",state="spare"} 2
        node_md_disks{device="md12",state="active"} 2
        node_md_disks{device="md12",state="failed"} 0
        node_md_disks{device="md12",state="spare"} 0
        node_md_disks{device="md120",state="active"} 2
        node_md_disks{device="md120",state="failed"} 0
        node_md_disks{device="md120",state="spare"} 0
        node_md_disks{device="md126",state="active"} 2
        node_md_disks{device="md126",state="failed"} 0
        node_md_disks{device="md126",state="spare"} 0
        node_md_disks{device="md127",state="active"} 2
        node_md_disks{device="md127",state="failed"} 0
        node_md_disks{device="md127",state="spare"} 0
        node_md_disks{device="md201",state="active"} 2
        node_md_disks{device="md201",state="failed"} 0
        node_md_disks{device="md201",state="spare"} 0
        node_md_disks{device="md219",state="active"} 0
        node_md_disks{device="md219",state="failed"} 0
        node_md_disks{device="md219",state="spare"} 3
        node_md_disks{device="md3",state="active"} 8
        node_md_disks{device="md3",state="failed"} 0
        node_md_disks{device="md3",state="spare"} 2
        node_md_disks{device="md4",state="active"} 0
        node_md_disks{device="md4",state="failed"} 1
        node_md_disks{device="md4",state="spare"} 1
        node_md_disks{device="md6",state="active"} 1
        node_md_disks{device="md6",state="failed"} 1
        node_md_disks{device="md6",state="spare"} 1
        node_md_disks{device="md7",state="active"} 3
        node_md_disks{device="md7",state="failed"} 1
        node_md_disks{device="md7",state="spare"} 0
        node_md_disks{device="md8",state="active"} 2
        node_md_disks{device="md8",state="failed"} 0
        node_md_disks{device="md8",state="spare"} 2
        node_md_disks{device="md9",state="active"} 4
        node_md_disks{device="md9",state="failed"} 2
        node_md_disks{device="md9",state="spare"} 1
        # HELP node_md_disks_required Total number of disks of device.
        # TYPE node_md_disks_required gauge
        node_md_disks_required{device="md0"} 2
        node_md_disks_required{device="md00"} 1
        node_md_disks_required{device="md10"} 2
        node_md_disks_required{device="md101"} 3
        node_md_disks_required{device="md11"} 2
        node_md_disks_required{device="md12"} 2
        node_md_disks_required{device="md120"} 2
        node_md_disks_required{device="md126"} 2
        node_md_disks_required{device="md127"} 2
        node_md_disks_required{device="md201"} 2
        node_md_disks_required{device="md219"} 0
        node_md_disks_required{device="md3"} 8
        node_md_disks_required{device="md4"} 0
        node_md_disks_required{device="md6"} 2
        node_md_disks_required{device="md7"} 4
        node_md_disks_required{device="md8"} 2
        node_md_disks_required{device="md9"} 4
        # HELP node_md_raid_disks Number of raid disks on device.
        # TYPE node_md_raid_disks gauge
        node_md_raid_disks{device="md0"} 2
        node_md_raid_disks{device="md1"} 2
        node_md_raid_disks{device="md10"} 4
        node_md_raid_disks{device="md4"} 3
        node_md_raid_disks{device="md5"} 3
        node_md_raid_disks{device="md6"} 4
        # HELP node_md_state Indicates the state of md-device.
        # TYPE node_md_state gauge
        node_md_state{device="md0",state="active"} 1
        node_md_state{device="md0",state="check"} 0
        node_md_state{device="md0",state="inactive"} 0
        node_md_state{device="md0",state="recovering"} 0
        node_md_state{device="md0",state="resync"} 0
        node_md_state{device="md00",state="active"} 1
        node_md_state{device="md00",state="check"} 0
        node_md_state{device="md00",state="inactive"} 0
        node_md_state{device="md00",state="recovering"} 0
        node_md_state{device="md00",state="resync"} 0
        node_md_state{device="md10",state="active"} 1
        node_md_state{device="md10",state="check"} 0
        node_md_state{device="md10",state="inactive"} 0
        node_md_state{device="md10",state="recovering"} 0
        node_md_state{device="md10",state="resync"} 0
        node_md_state{device="md101",state="active"} 1
        node_md_state{device="md101",state="check"} 0
        node_md_state{device="md101",state="inactive"} 0
        node_md_state{device="md101",state="recovering"} 0
        node_md_state{device="md101",state="resync"} 0
        node_md_state{device="md11",state="active"} 0
        node_md_state{device="md11",state="check"} 0
        node_md_state{device="md11",state="inactive"} 0
        node_md_state{device="md11",state="recovering"} 0
        node_md_state{device="md11",state="resync"} 1
        node_md_state{device="md12",state="active"} 1
        node_md_state{device="md12",state="check"} 0
        node_md_state{device="md12",state="inactive"} 0
        node_md_state{device="md12",state="recovering"} 0
        node_md_state{device="md12",state="resync"} 0
        node_md_state{device="md120",state="active"} 1
        node_md_state{device="md120",state="check"} 0
        node_md_state{device="md120",state="inactive"} 0
        node_md_state{device="md120",state="recovering"} 0
        node_md_state{device="md120",state="resync"} 0
        node_md_state{device="md126",state="active"} 1
        node_md_state{device="md126",state="check"} 0
        node_md_state{device="md126",state="inactive"} 0
        node_md_state{device="md126",state="recovering"} 0
        node_md_state{device="md126",state="resync"} 0
        node_md_state{device="md127",state="active"} 1
        node_md_state{device="md127",state="check"} 0
        node_md_state{device="md127",state="inactive"} 0
        node_md_state{device="md127",state="recovering"} 0
        node_md_state{device="md127",state="resync"} 0
        node_md_state{device="md201",state="active"} 0
        node_md_state{device="md201",state="check"} 1
        node_md_state{device="md201",state="inactive"} 0
        node_md_state{device="md201",state="recovering"} 0
        node_md_state{device="md201",state="resync"} 0
        node_md_state{device="md219",state="active"} 0
        node_md_state{device="md219",state="check"} 0
        node_md_state{device="md219",state="inactive"} 1
        node_md_state{device="md219",state="recovering"} 0
        node_md_state{device="md219",state="resync"} 0
        node_md_state{device="md3",state="active"} 1
        node_md_state{device="md3",state="check"} 0
        node_md_state{device="md3",state="inactive"} 0
        node_md_state{device="md3",state="recovering"} 0
        node_md_state{device="md3",state="resync"} 0
        node_md_state{device="md4",state="active"} 0
        node_md_state{device="md4",state="check"} 0
        node_md_state{device="md4",state="inactive"} 1
        node_md_state{device="md4",state="recovering"} 0
        node_md_state{device="md4",state="resync"} 0
        node_md_state{device="md6",state="active"} 0
        node_md_state{device="md6",state="check"} 0
        node_md_state{device="md6",state="inactive"} 0
        node_md_state{device="md6",state="recovering"} 1
        node_md_state{device="md6",state="resync"} 0
        node_md_state{device="md7",state="active"} 1
        node_md_state{device="md7",state="check"} 0
        node_md_state{device="md7",state="inactive"} 0
        node_md_state{device="md7",state="recovering"} 0
        node_md_state{device="md7",state="resync"} 0
        node_md_state{device="md8",state="active"} 0
        node_md_state{device="md8",state="check"} 0
        node_md_state{device="md8",state="inactive"} 0
        node_md_state{device="md8",state="recovering"} 0
        node_md_state{device="md8",state="resync"} 1
        node_md_state{device="md9",state="active"} 0
        node_md_state{device="md9",state="check"} 0
        node_md_state{device="md9",state="inactive"} 0
        node_md_state{device="md9",state="recovering"} 0
        node_md_state{device="md9",state="resync"} 1
`
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelError,
		AddSource: true,
	}))
	collector, err := NewMdadmCollector(logger)
	if err != nil {
		panic(err)
	}
	c, err := NewTestMdadmCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	sink := make(chan prometheus.Metric)
	go func() {
		err := collector.Update(sink)
		if err != nil {
			panic(err)
		}
		close(sink)
	}()

	err = testutil.GatherAndCompare(reg, strings.NewReader(testcase))
	if err != nil {
		t.Fatal(err)
	}
}
