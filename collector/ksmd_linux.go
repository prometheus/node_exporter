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

//go:build !noksmd

package collector

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
)

const ksmdCollectorSubsystem = "ksmd"

var (
	ksmdFiles = []string{
		"full_scans",
		"general_profit",
		"merge_across_nodes",
		"pages_shared",
		"pages_sharing",
		"pages_to_scan",
		"pages_unshared",
		"pages_volatile",
		"run",
		"sleep_millisecs",
	}

	// Help text from https://docs.kernel.org/admin-guide/mm/ksm.html
	ksmdFullScansDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ksmdCollectorSubsystem, "full_scans_total"),
		"How many times all mergeable areas have been scanned.",
		nil, nil,
	)
	ksmdGeneralProfitDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ksmdCollectorSubsystem, "general_profit_bytes"),
		"How effective is KSM. ksm_saved_pages * sizeof(page) - (all_rmap_items) * sizeof(rmap_item)",
		nil, nil,
	)
	ksmdMergeAcrossNodesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ksmdCollectorSubsystem, "merge_across_nodes"),
		"Specifies if pages from different NUMA nodes can be merged.",
		nil, nil,
	)
	ksmdPagesSharedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ksmdCollectorSubsystem, "pages_shared"),
		"How many shared pages are being used.",
		nil, nil,
	)
	ksmdPagesSharingDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ksmdCollectorSubsystem, "pages_sharing"),
		"How many more sites are sharing them i.e. how much saved.",
		nil, nil,
	)
	ksmdPagesToScanDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ksmdCollectorSubsystem, "pages_to_scan"),
		"How many pages to scan before ksmd goes to sleep.",
		nil, nil,
	)
	ksmdPagesUnsharedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ksmdCollectorSubsystem, "pages_unshared"),
		"How many pages unique but repeatedly checked for merging.",
		nil, nil,
	)
	ksmdPagesVolatileDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ksmdCollectorSubsystem, "pages_volatile"),
		"How many pages changing too fast to be placed in a tree.",
		nil, nil,
	)
	ksmdRunDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ksmdCollectorSubsystem, "run"),
		"The status of ksmd. stopped = 0, running = 1",
		nil, nil,
	)
	ksmdSleepSecondsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ksmdCollectorSubsystem, "sleep_seconds"),
		"How many seconds ksmd should sleep before next scan.",
		nil, nil,
	)
)

type ksmdCollector struct {
	logger *slog.Logger
}

func init() {
	registerCollector(ksmdCollectorSubsystem, defaultDisabled, NewKsmdCollector)
}

// NewKsmdCollector returns a new Collector exposing kernel/system statistics.
func NewKsmdCollector(logger *slog.Logger) (Collector, error) {
	return &ksmdCollector{logger: logger}, nil
}

// Update implements Collector and exposes kernel and system statistics.
func (c *ksmdCollector) Update(ch chan<- prometheus.Metric) error {
	for _, n := range ksmdFiles {
		val, err := readUintFromFile(sysFilePath(filepath.Join("kernel/mm/ksm", n)))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				c.logger.Debug("ksmd file not found, skipping", "file", n)
				continue
			}
			return err
		}

		v := float64(val)
		switch n {
		case "full_scans":
			ch <- prometheus.MustNewConstMetric(ksmdFullScansDesc, prometheus.CounterValue, v)
		case "general_profit":
			ch <- prometheus.MustNewConstMetric(ksmdGeneralProfitDesc, prometheus.GaugeValue, v)
		case "merge_across_nodes":
			ch <- prometheus.MustNewConstMetric(ksmdMergeAcrossNodesDesc, prometheus.GaugeValue, v)
		case "pages_shared":
			ch <- prometheus.MustNewConstMetric(ksmdPagesSharedDesc, prometheus.GaugeValue, v)
		case "pages_sharing":
			ch <- prometheus.MustNewConstMetric(ksmdPagesSharingDesc, prometheus.GaugeValue, v)
		case "pages_to_scan":
			ch <- prometheus.MustNewConstMetric(ksmdPagesToScanDesc, prometheus.GaugeValue, v)
		case "pages_unshared":
			ch <- prometheus.MustNewConstMetric(ksmdPagesUnsharedDesc, prometheus.GaugeValue, v)
		case "pages_volatile":
			ch <- prometheus.MustNewConstMetric(ksmdPagesVolatileDesc, prometheus.GaugeValue, v)
		case "run":
			ch <- prometheus.MustNewConstMetric(ksmdRunDesc, prometheus.GaugeValue, v)
		case "sleep_millisecs":
			ch <- prometheus.MustNewConstMetric(ksmdSleepSecondsDesc, prometheus.GaugeValue, v/1000)
		}
	}

	return nil
}
