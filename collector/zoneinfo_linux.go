// Copyright 2020 The Prometheus Authors
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

package collector

import (
	"fmt"
	"reflect"

	"github.com/go-kit/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

const zoneinfoSubsystem = "zoneinfo"

type zoneinfoCollector struct {
	gaugeMetricDescs   map[string]*prometheus.Desc
	counterMetricDescs map[string]*prometheus.Desc
	logger             log.Logger
	fs                 procfs.FS
}

func init() {
	registerCollector("zoneinfo", defaultDisabled, NewZoneinfoCollector)
}

// NewZoneinfoCollector returns a new Collector exposing zone stats.
func NewZoneinfoCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}
	return &zoneinfoCollector{
		gaugeMetricDescs:   createGaugeMetricDescriptions(),
		counterMetricDescs: createCounterMetricDescriptions(),
		logger:             logger,
		fs:                 fs,
	}, nil
}

func (c *zoneinfoCollector) Update(ch chan<- prometheus.Metric) error {
	metrics, err := c.fs.Zoneinfo()
	if err != nil {
		return fmt.Errorf("couldn't get zoneinfo: %w", err)
	}
	for _, metric := range metrics {
		node := metric.Node
		zone := metric.Zone
		metricStruct := reflect.ValueOf(metric)
		typeOfMetricStruct := metricStruct.Type()
		for i := 0; i < metricStruct.NumField(); i++ {
			value := reflect.Indirect(metricStruct.Field(i))
			if value.Kind() != reflect.Int64 {
				continue
			}
			metricName := typeOfMetricStruct.Field(i).Name
			desc, ok := c.gaugeMetricDescs[metricName]
			metricType := prometheus.GaugeValue
			if !ok {
				desc = c.counterMetricDescs[metricName]
				metricType = prometheus.CounterValue
			}
			ch <- prometheus.MustNewConstMetric(desc, metricType,
				float64(reflect.Indirect(metricStruct.Field(i)).Int()),
				node, zone)
		}
		for i, value := range metric.Protection {
			metricName := fmt.Sprintf("protection_%d", i)
			desc, ok := c.gaugeMetricDescs[metricName]
			if !ok {
				desc = prometheus.NewDesc(
					prometheus.BuildFQName(namespace, zoneinfoSubsystem, metricName),
					fmt.Sprintf("Protection array %d. field", i),
					[]string{"node", "zone"}, nil)
				c.gaugeMetricDescs[metricName] = desc
			}
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue,
				float64(*value), node, zone)
		}

	}
	return nil
}
func createGaugeMetricDescriptions() map[string]*prometheus.Desc {
	return map[string]*prometheus.Desc{
		"NrFreePages": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_free_pages"),
			"Total number of free pages in the zone",
			[]string{"node", "zone"}, nil),
		"Min": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "min_pages"),
			"Zone watermark pages_min",
			[]string{"node", "zone"}, nil),
		"Low": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "low_pages"),
			"Zone watermark pages_low",
			[]string{"node", "zone"}, nil),
		"High": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "high_pages"),
			"Zone watermark pages_high",
			[]string{"node", "zone"}, nil),
		"Scanned": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "scanned_pages"),
			"Pages scanned since last reclaim",
			[]string{"node", "zone"}, nil),
		"Spanned": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "spanned_pages"),
			"Total pages spanned by the zone, including holes",
			[]string{"node", "zone"}, nil),
		"Present": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "present_pages"),
			"Physical pages existing within the zone",
			[]string{"node", "zone"}, nil),
		"Managed": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "managed_pages"),
			"Present pages managed by the buddy system",
			[]string{"node", "zone"}, nil),
		"NrActiveAnon": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_active_anon_pages"),
			"Number of anonymous pages recently more used",
			[]string{"node", "zone"}, nil),
		"NrInactiveAnon": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_inactive_anon_pages"),
			"Number of anonymous pages recently less used",
			[]string{"node", "zone"}, nil),
		"NrIsolatedAnon": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_isolated_anon_pages"),
			"Temporary isolated pages from anon lru",
			[]string{"node", "zone"}, nil),
		"NrAnonPages": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_anon_pages"),
			"Number of anonymous pages currently used by the system",
			[]string{"node", "zone"}, nil),
		"NrAnonTransparentHugepages": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_anon_transparent_hugepages"),
			"Number of anonymous transparent huge pages currently used by the system",
			[]string{"node", "zone"}, nil),
		"NrActiveFile": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_active_file_pages"),
			"Number of active pages with file-backing",
			[]string{"node", "zone"}, nil),
		"NrInactiveFile": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_inactive_file_pages"),
			"Number of inactive pages with file-backing",
			[]string{"node", "zone"}, nil),
		"NrIsolatedFile": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_isolated_file_pages"),
			"Temporary isolated pages from file lru",
			[]string{"node", "zone"}, nil),
		"NrFilePages": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_file_pages"),
			"Number of file pages",
			[]string{"node", "zone"}, nil),
		"NrSlabReclaimable": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_slab_reclaimable_pages"),
			"Number of reclaimable slab pages",
			[]string{"node", "zone"}, nil),
		"NrSlabUnreclaimable": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_slab_unreclaimable_pages"),
			"Number of unreclaimable slab pages",
			[]string{"node", "zone"}, nil),
		"NrMlockStack": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_mlock_stack_pages"),
			"mlock()ed pages found and moved off LRU",
			[]string{"node", "zone"}, nil),
		"NrKernelStack": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_kernel_stacks"),
			"Number of kernel stacks",
			[]string{"node", "zone"}, nil),
		"NrMapped": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_mapped_pages"),
			"Number of mapped pages",
			[]string{"node", "zone"}, nil),
		"NrDirty": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_dirty_pages"),
			"Number of dirty pages",
			[]string{"node", "zone"}, nil),
		"NrWriteback": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_writeback_pages"),
			"Number of writeback pages",
			[]string{"node", "zone"}, nil),
		"NrUnevictable": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_unevictable_pages"),
			"Number of unevictable pages",
			[]string{"node", "zone"}, nil),
		"NrShmem": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_shmem_pages"),
			"Number of shmem pages (included tmpfs/GEM pages)",
			[]string{"node", "zone"}, nil),
	}

}
func createCounterMetricDescriptions() map[string]*prometheus.Desc {
	return map[string]*prometheus.Desc{
		"NrDirtied": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_dirtied_total"),
			"Page dirtyings since bootup",
			[]string{"node", "zone"}, nil),
		"NrWritten": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_written_total"),
			"Page writings since bootup",
			[]string{"node", "zone"}, nil),
		"NumaHit": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "numa_hit_total"),
			"Allocated in intended node",
			[]string{"node", "zone"}, nil),
		"NumaMiss": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "numa_miss_total"),
			"Allocated in non intended node",
			[]string{"node", "zone"}, nil),
		"NumaForeign": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "numa_foreign_total"),
			"Was intended here, hit elsewhere",
			[]string{"node", "zone"}, nil),
		"NumaInterleave": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "numa_interleave_total"),
			"Interleaver preferred this zone",
			[]string{"node", "zone"}, nil),
		"NumaLocal": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "numa_local_total"),
			"Allocation from local node",
			[]string{"node", "zone"}, nil),
		"NumaOther": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "numa_other_total"),
			"Allocation from other node",
			[]string{"node", "zone"}, nil),
	}
}
