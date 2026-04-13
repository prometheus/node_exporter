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

//go:build !nocifs

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

type testCIFSCollector struct {
	cc Collector
}

func (c testCIFSCollector) Collect(ch chan<- prometheus.Metric) {
	c.cc.Update(ch)
}

func (c testCIFSCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func TestCIFSStats(t *testing.T) {
	testcase := `# HELP node_cifs_read_bytes_total Total bytes read from this share.
	# TYPE node_cifs_read_bytes_total counter
	node_cifs_read_bytes_total{share="\\\\server1\\share1"} 123456
	node_cifs_read_bytes_total{share="\\\\server2\\share2"} 789012
	# HELP node_cifs_session_reconnects_total Total number of CIFS session reconnects.
	# TYPE node_cifs_session_reconnects_total counter
	node_cifs_session_reconnects_total 3
	# HELP node_cifs_sessions Number of active CIFS sessions.
	# TYPE node_cifs_sessions gauge
	node_cifs_sessions 2
	# HELP node_cifs_share_reconnects_total Total number of CIFS share reconnects.
	# TYPE node_cifs_share_reconnects_total counter
	node_cifs_share_reconnects_total 5
	# HELP node_cifs_shares Number of unique mount targets.
	# TYPE node_cifs_shares gauge
	node_cifs_shares 3
	# HELP node_cifs_smbs_total Total number of SMBs sent for this share.
	# TYPE node_cifs_smbs_total counter
	node_cifs_smbs_total{share="\\\\server1\\share1"} 1234
	node_cifs_smbs_total{share="\\\\server2\\share2"} 5678
	# HELP node_cifs_vfs_operations_total Total number of VFS operations.
	# TYPE node_cifs_vfs_operations_total counter
	node_cifs_vfs_operations_total 67
	# HELP node_cifs_write_bytes_total Total bytes written to this share.
	# TYPE node_cifs_write_bytes_total counter
	node_cifs_write_bytes_total{share="\\\\server1\\share1"} 654321
	node_cifs_write_bytes_total{share="\\\\server2\\share2"} 210987
	`
	*procPath = "fixtures/proc"

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := NewCIFSCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	reg := prometheus.NewRegistry()
	reg.MustRegister(&testCIFSCollector{cc: c})

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
