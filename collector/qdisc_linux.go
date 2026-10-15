// Copyright 2017 The Prometheus Authors
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

//go:build !noqdisc

package collector

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin/v2"
	"github.com/ema/qdisc"
	"github.com/prometheus/client_golang/prometheus"
)

type qdiscStatCollector struct {
	logger       *slog.Logger
	deviceFilter deviceFilter
	bytes        typedDesc
	packets      typedDesc
	drops        typedDesc
	requeues     typedDesc
	overlimits   typedDesc
	qlength      typedDesc
	backlog      typedDesc
}

var (
	collectorQdisc                 = kingpin.Flag("collector.qdisc.fixtures", "test fixtures to use for qdisc collector end-to-end testing").Default("").String()
	collectorQdiscIncludeChild     = kingpin.Flag("collector.qdisc.include-child", "Also expose non-root qdiscs, adding 'parent' and 'handle' labels to all qdisc metrics. Note that packets traversing a child qdisc are also counted by its parents, so generic counters must not be aggregated across a device without grouping by these labels.").Bool()
	collectorQdiscDeviceInclude    = kingpin.Flag("collector.qdisc.device-include", "Regexp of qdisc devices to include (mutually exclusive to device-exclude).").String()
	oldCollectorQdiskDeviceInclude = kingpin.Flag("collector.qdisk.device-include", "DEPRECATED: Use collector.qdisc.device-include").Hidden().String()
	collectorQdiscDeviceExclude    = kingpin.Flag("collector.qdisc.device-exclude", "Regexp of qdisc devices to exclude (mutually exclusive to device-include).").String()
	oldCollectorQdiskDeviceExclude = kingpin.Flag("collector.qdisk.device-exclude", "DEPRECATED: Use collector.qdisc.device-exclude").Hidden().String()
)

func init() {
	registerCollector("qdisc", defaultDisabled, NewQdiscStatCollector)
}

// NewQdiscStatCollector returns a new Collector exposing queuing discipline statistics.
func NewQdiscStatCollector(logger *slog.Logger) (Collector, error) {
	if *oldCollectorQdiskDeviceInclude != "" {
		if *collectorQdiscDeviceInclude == "" {
			logger.Warn("--collector.qdisk.device-include is DEPRECATED and will be removed in 2.0.0, use --collector.qdisc.device-include")
			*collectorQdiscDeviceInclude = *oldCollectorQdiskDeviceInclude
		} else {
			return nil, fmt.Errorf("--collector.qdisk.device-include and --collector.qdisc.device-include are mutually exclusive")
		}
	}

	if *oldCollectorQdiskDeviceExclude != "" {
		if *collectorQdiscDeviceExclude == "" {
			logger.Warn("--collector.qdisk.device-exclude is DEPRECATED and will be removed in 2.0.0, use --collector.qdisc.device-exclude")
			*collectorQdiscDeviceExclude = *oldCollectorQdiskDeviceExclude
		} else {
			return nil, fmt.Errorf("--collector.qdisk.device-exclude and --collector.qdisc.device-exclude are mutually exclusive")
		}
	}

	if *collectorQdiscDeviceExclude != "" && *collectorQdiscDeviceInclude != "" {
		return nil, fmt.Errorf("collector.qdisc.device-include and collector.qdisc.device-exclude are mutaly exclusive")
	}

	labels := []string{"device", "kind"}
	if *collectorQdiscIncludeChild {
		labels = append(labels, "parent", "handle")
	}

	return &qdiscStatCollector{
		bytes: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "qdisc", "bytes_total"),
			"Number of bytes sent.",
			labels, nil,
		), prometheus.CounterValue},
		packets: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "qdisc", "packets_total"),
			"Number of packets sent.",
			labels, nil,
		), prometheus.CounterValue},
		drops: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "qdisc", "drops_total"),
			"Number of packets dropped.",
			labels, nil,
		), prometheus.CounterValue},
		requeues: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "qdisc", "requeues_total"),
			"Number of packets dequeued, not transmitted, and requeued.",
			labels, nil,
		), prometheus.CounterValue},
		overlimits: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "qdisc", "overlimits_total"),
			"Number of overlimit packets.",
			labels, nil,
		), prometheus.CounterValue},
		qlength: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "qdisc", "current_queue_length"),
			"Number of packets currently in queue to be sent.",
			labels, nil,
		), prometheus.GaugeValue},
		backlog: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "qdisc", "backlog"),
			"Number of bytes currently in queue to be sent.",
			labels, nil,
		), prometheus.GaugeValue},
		logger:       logger,
		deviceFilter: newDeviceFilter(*collectorQdiscDeviceExclude, *collectorQdiscDeviceInclude),
	}, nil
}

// tcHandleString renders a tc handle in the usual major:minor hex notation.
func tcHandleString(handle uint32) string {
	return fmt.Sprintf("%x:%x", handle>>16, handle&0xffff)
}

// tcParentString renders the parent of a qdisc for use as a label value.
// The ema/qdisc library normalizes TC_H_ROOT to 0, which is rendered as
// "root" to match tc's own terminology.
func tcParentString(parent uint32) string {
	if parent == 0 {
		return "root"
	}
	return tcHandleString(parent)
}

func testQdiscGet(fixtures string) ([]qdisc.QdiscInfo, error) {
	var res []qdisc.QdiscInfo

	b, err := os.ReadFile(filepath.Join(fixtures, "results.json"))
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(b, &res)
	return res, err
}

func (c *qdiscStatCollector) Update(ch chan<- prometheus.Metric) error {
	var msgs []qdisc.QdiscInfo
	var err error

	fixtures := *collectorQdisc

	if fixtures == "" {
		msgs, err = qdisc.Get()
	} else {
		msgs, err = testQdiscGet(fixtures)
	}

	if err != nil {
		return err
	}

	for _, msg := range msgs {
		// By default, only report root qdisc information. Packets traversing
		// a child qdisc are also counted by its parents, so reporting only
		// root qdiscs keeps per-device aggregations free of double counting.
		if !*collectorQdiscIncludeChild && msg.Parent != 0 {
			continue
		}

		if c.deviceFilter.ignored(msg.IfaceName) {
			continue
		}

		labelValues := []string{msg.IfaceName, msg.Kind}
		if *collectorQdiscIncludeChild {
			labelValues = append(labelValues, tcParentString(msg.Parent), tcHandleString(msg.Handle))
		}

		ch <- c.bytes.mustNewConstMetric(float64(msg.Bytes), labelValues...)
		ch <- c.packets.mustNewConstMetric(float64(msg.Packets), labelValues...)
		ch <- c.drops.mustNewConstMetric(float64(msg.Drops), labelValues...)
		ch <- c.requeues.mustNewConstMetric(float64(msg.Requeues), labelValues...)
		ch <- c.overlimits.mustNewConstMetric(float64(msg.Overlimits), labelValues...)
		ch <- c.qlength.mustNewConstMetric(float64(msg.Qlen), labelValues...)
		ch <- c.backlog.mustNewConstMetric(float64(msg.Backlog), labelValues...)
	}

	return nil
}
