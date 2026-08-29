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

//go:build !noresctrl

package collector

import (
	"errors"
	"log/slog"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/resctrlfs"
)

const resctrlSubsystem = "resctrl"

// The resctrl filesystem exposes the cache and memory bandwidth monitoring of
// the CPU (Intel RDT, AMD PQoS, Arm MPAM), documented in
// https://docs.kernel.org/filesystems/resctrl.html and
// https://docs.kernel.org/arch/arm64/mpam.html.
//
// It is not mounted by default, hence this collector is disabled by default.
// Mount it with:
//
//	mount -t resctrl resctrl /sys/fs/resctrl
var resctrlPath = kingpin.Flag("collector.resctrl.path", "resctrl filesystem mountpoint.").Default(resctrlfs.DefaultMountPoint).String()

type resctrlCollector struct{}

func init() {
	registerCollector("resctrl", defaultDisabled, NewResctrlCollector)
}

var (
	resctrlMemoryBandwidth = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, resctrlSubsystem, "memory_bandwidth_bytes_total"),
		"Bytes moved between the last level cache and memory. scope=total includes remote sockets, scope=local is memory attached to this domain.",
		[]string{"domain", "scope"}, nil,
	)
	resctrlLLCOccupancy = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, resctrlSubsystem, "llc_occupancy_bytes"),
		"Last level cache bytes occupied in this domain.",
		[]string{"domain"}, nil,
	)
)

// NewResctrlCollector returns a new Collector exposing resctrl monitoring counters.
func NewResctrlCollector(_ *slog.Logger) (Collector, error) {
	return &resctrlCollector{}, nil
}

func (c *resctrlCollector) Update(ch chan<- prometheus.Metric) error {
	// The filesystem is opened on every scrape because it can be mounted and
	// unmounted while node_exporter runs.
	fs, err := resctrlfs.NewFS(*resctrlPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoData
		}
		return err
	}

	domains, err := fs.MonData()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoData
		}
		return err
	}

	found := false
	for _, domain := range domains {
		// The metrics describe the last level cache, so only L3 fits them.
		if domain.Resource != "L3" {
			continue
		}
		found = true

		// The bandwidth counters wrap. No correction is applied: Prometheus
		// reads a wrap as a counter reset, which is the closest available
		// approximation.
		if v := domain.MBMTotalBytes; v != nil {
			ch <- prometheus.MustNewConstMetric(resctrlMemoryBandwidth, prometheus.CounterValue, float64(*v), domain.ID, "total")
		}
		if v := domain.MBMLocalBytes; v != nil {
			ch <- prometheus.MustNewConstMetric(resctrlMemoryBandwidth, prometheus.CounterValue, float64(*v), domain.ID, "local")
		}
		if v := domain.LLCOccupancy; v != nil {
			ch <- prometheus.MustNewConstMetric(resctrlLLCOccupancy, prometheus.GaugeValue, float64(*v), domain.ID)
		}
	}

	if !found {
		return ErrNoData
	}
	return nil
}
