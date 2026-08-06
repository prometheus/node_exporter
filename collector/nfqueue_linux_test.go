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

//go:build !nonfqueue

package collector

import (
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type testNFQueueCollector struct {
	nc Collector
}

func (c testNFQueueCollector) Collect(ch chan<- prometheus.Metric) {
	c.nc.Update(ch)
}

func (c testNFQueueCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func TestNFQueueStats(t *testing.T) {
	testcase := `# HELP node_nfqueue_packets_dropped_total Total number of packets dropped.
	# TYPE node_nfqueue_packets_dropped_total counter
	node_nfqueue_packets_dropped_total{queue="0",reason="queue_full"} 0
	node_nfqueue_packets_dropped_total{queue="0",reason="user"} 0
	node_nfqueue_packets_dropped_total{queue="1",reason="queue_full"} 150
	node_nfqueue_packets_dropped_total{queue="1",reason="user"} 10
	node_nfqueue_packets_dropped_total{queue="2",reason="queue_full"} 20
	node_nfqueue_packets_dropped_total{queue="2",reason="user"} 5
	# HELP node_nfqueue_info Non-numeric metadata about the queue (value is always 1).
	# TYPE node_nfqueue_info gauge
	node_nfqueue_info{copy_mode="packet",copy_range="65531",peer_portid="31621",queue="0"} 1
	node_nfqueue_info{copy_mode="meta",copy_range="1024",peer_portid="31622",queue="1"} 1
	node_nfqueue_info{copy_mode="none",copy_range="512",peer_portid="31623",queue="2"} 1
	# HELP node_nfqueue_queue_length Current number of packets waiting in the queue.
	# TYPE node_nfqueue_queue_length gauge
	node_nfqueue_queue_length{queue="0"} 0
	node_nfqueue_queue_length{queue="1"} 100
	node_nfqueue_queue_length{queue="2"} 25
	`
	*procPath = "fixtures/proc"

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := NewNFQueueCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	reg := prometheus.NewRegistry()
	reg.MustRegister(&testNFQueueCollector{nc: c})

	err = testutil.GatherAndCompare(reg, strings.NewReader(testcase))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNFQueueStatsErrNoData(t *testing.T) {
	*procPath = t.TempDir() // valid dir, but no nfnetlink_queue file inside

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := NewNFQueueCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan prometheus.Metric)
	err = c.Update(ch)
	if !errors.Is(err, ErrNoData) {
		t.Fatalf("expected ErrNoData, got: %v", err)
	}
}
