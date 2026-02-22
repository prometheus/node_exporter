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
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type nfqueueCollector struct {
	fs             procfs.FS
	queueLength    *prometheus.Desc
	packetsDropped *prometheus.Desc
	info           *prometheus.Desc
	logger         *slog.Logger
}

func init() {
	registerCollector("nfqueue", defaultDisabled, NewNFQueueCollector)
}

// NewNFQueueCollector returns a new Collector exposing netfilter queue stats
// from /proc/net/netfilter/nfnetlink_queue.
func NewNFQueueCollector(logger *slog.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &nfqueueCollector{
		fs: fs,
		queueLength: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "nfqueue", "queue_length"),
			"Current number of packets waiting in the queue.",
			[]string{"queue"}, nil,
		),
		packetsDropped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "nfqueue", "packets_dropped_total"),
			"Total number of packets dropped.",
			[]string{"queue", "reason"}, nil,
		),
		// nfqueue_info cardinality:
		//
		// queue       | Queue ID (uint16): 0-65535 in theory, but in practice single digits to low tens.
		// peer_portid | PID of the userspace listener: one per queue, only changes on process restart.
		// copy_mode   | Packet copy mode: only 3 possible values (none/meta/packet) - set per queue.
		// copy_range  | Bytes copied to userspace: 0-65535, but typically a single fixed value per queue (e.g. 65535 or MTU size).
		//
		// The total number of time series = number of NFQUEUE queues. Even in an extreme case of 100 queues (unrealistically many),
		// that's just 100 series.
		info: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "nfqueue", "info"),
			"Non-numeric metadata about the queue (value is always 1).",
			[]string{"queue", "peer_portid", "copy_mode", "copy_range"}, nil,
		),
		logger: logger,
	}, nil
}

func (c *nfqueueCollector) Update(ch chan<- prometheus.Metric) error {
	queues, err := c.fs.NFNetLinkQueue()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.logger.Debug("nfqueue: file not found, NFQUEUE probably not in use")
			return ErrNoData
		}
		return fmt.Errorf("failed to retrieve nfqueue stats: %w", err)
	}

	for _, q := range queues {
		queueID := strconv.FormatUint(uint64(q.QueueID), 10)
		ch <- prometheus.MustNewConstMetric(
			c.queueLength, prometheus.GaugeValue,
			float64(q.QueueTotal), queueID,
		)
		ch <- prometheus.MustNewConstMetric(
			c.packetsDropped, prometheus.CounterValue,
			float64(q.QueueDropped), queueID, "queue_full",
		)
		ch <- prometheus.MustNewConstMetric(
			c.packetsDropped, prometheus.CounterValue,
			float64(q.QueueUserDropped), queueID, "user",
		)
		ch <- prometheus.MustNewConstMetric(
			c.info, prometheus.GaugeValue, 1,
			queueID,
			strconv.FormatUint(uint64(q.PeerPID), 10),
			nfqueueCopyModeString(q.CopyMode),
			strconv.FormatUint(uint64(q.CopyRange), 10),
		)
	}
	return nil
}

func nfqueueCopyModeString(mode uint) string {
	switch mode {
	case 0:
		return "none"
	case 1:
		return "meta"
	case 2:
		return "packet"
	default:
		return strconv.FormatUint(uint64(mode), 10)
	}
}
