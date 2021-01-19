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
	"regexp"
	"strings"

	"github.com/go-kit/kit/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

const zoneinfoSubsystem = "zoneinfo"

type zoneinfoCollector struct {
	metricDescs map[string]*prometheus.Desc
	logger      log.Logger
	fs          procfs.FS
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
		metricDescs: createMetricDescriptions(),
		logger:      logger,
		fs:          fs,
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
			desc := c.metricDescs[typeOfMetricStruct.Field(i).Name]
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue,
				float64(reflect.Indirect(metricStruct.Field(i)).Int()),
				node, zone)
		}
		for i, value := range metric.Protection {
			metricName := fmt.Sprintf("protection_%d", i)
			desc, ok := c.metricDescs[metricName]
			if !ok {
				desc = prometheus.NewDesc(
					prometheus.BuildFQName(namespace, zoneinfoSubsystem, metricName),
					fmt.Sprintf("Protection array %d. field", i),
					[]string{"node", "zone"}, nil)
				c.metricDescs[metricName] = desc
			}
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue,
				float64(*value), node, zone)
		}

	}
	return nil
}
func  createMetricDescriptions() map[string]*prometheus.Desc  {
	return map[string]*prometheus.Desc{
		"NrFreePages" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_free_pages" ),
			"Total number of free pages in the zone",
			[]string{"node", "zone"}, nil),
		"Min" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "min" ),
			"Zone watermark pages_min",
			[]string{"node", "zone"}, nil),
		"Low" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "low" ),
			"Zone watermark pages_low",
			[]string{"node", "zone"}, nil),
		"High" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "high" ),
			"Zone watermark pages_high",
			[]string{"node", "zone"}, nil),
		"Scanned" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "scanned" ),
			"Pages scanned since last reclaim",
			[]string{"node", "zone"}, nil),
		"Spanned" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "spanned" ),
			"Total pages spanned by the zone, including holes",
			[]string{"node", "zone"}, nil),
		"Present" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "present" ),
			"Physical pages existing within the zone",
			[]string{"node", "zone"}, nil),
		"Managed" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "managed" ),
			"Present pages managed by the buddy system",
			[]string{"node", "zone"}, nil),
		"NrActiveAnon" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_active_anon" ),
			"Number of anonymous pages recently more used",
			[]string{"node", "zone"}, nil),
		"NrInactiveAnon" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_inactive_anon" ),
			"Number of anonymous pages recently less used",
			[]string{"node", "zone"}, nil),
		"NrIsolatedAnon" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_isolated_anon" ),
			"Temporary isolated pages from anon lru",
			[]string{"node", "zone"}, nil),
		"NrAnonPages" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_anon_pages" ),
			"Number of anonymous pages currently used by the system",
			[]string{"node", "zone"}, nil),
		"NrAnonTransparentHugepages" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_anon_transparent_hugepages" ),
			"Number of anonymous transparent huge pages currently used by the system",
			[]string{"node", "zone"}, nil),
		"NrActiveFile" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_active_anon" ),
			"Number of active pages with file-backing ",
			[]string{"node", "zone"}, nil),
		"NrInactiveFile" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_inactive_anon" ),
			"Number of inactive pages with file-backing ",
			[]string{"node", "zone"}, nil),
		"NrIsolatedFile" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_isolated_file" ),
			"Temporary isolated pages from file lru",
			[]string{"node", "zone"}, nil),
		"NrFilePages" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_file_pages" ),
			"Number of file pages",
			[]string{"node", "zone"}, nil),
		"NrSlabReclaimable" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_slab_reclaimable" ),
			"Number of reclaimable slab pages",
			[]string{"node", "zone"}, nil),
		"NrSlabUnreclaimable" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_slab_unreclaimable" ),
			"Number of unreclaimable slab pages",
			[]string{"node", "zone"}, nil),
		"NrMlockStack" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_mlock_stack" ),
			"mlock()ed pages found and moved off LRU",
			[]string{"node", "zone"}, nil),
		"NrKernelStack" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_kernel_stack" ),
			"Amount of memory allocated to kernel stacks",
			[]string{"node", "zone"}, nil),
		"NrMapped" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_mapped" ),
			"Mapped paged",
			[]string{"node", "zone"}, nil),
		"NrDirty" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_dirty" ),
			"Number of dirty pages",
			[]string{"node", "zone"}, nil),
		"NrWriteback" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_writeback" ),
			"Number of writeback pages",
			[]string{"node", "zone"}, nil),
		"NrUnevictable" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_unevictable" ),
			"Number of unevictable pages",
			[]string{"node", "zone"}, nil),
		"NrShmem" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_shmem" ),
			"Shmem pages (included tmpfs/GEM pages)",
			[]string{"node", "zone"}, nil),
		"NrDirtied" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_dirtied" ),
			"Page dirtyings since bootup",
			[]string{"node", "zone"}, nil),
		"NrWritten" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "nr_written" ),
			"Page writings since bootu",
			[]string{"node", "zone"}, nil),
		"NumaHit" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "numa_hit" ),
			"Allocated in intended node",
			[]string{"node", "zone"}, nil),
		"NumaMiss" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "numa_hit" ),
			"Allocated in non intended node",
			[]string{"node", "zone"}, nil),
		"NumaForeign" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "numa_foreign" ),
			"Was intended here, hit elsewhere",
			[]string{"node", "zone"}, nil),
		"NumaInterleave" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "numa_interleave" ),
			"Interleaver preferred this zone",
			[]string{"node", "zone"}, nil),
		"NumaLocal" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "numa_local" ),
			"Allocation from local node",
			[]string{"node", "zone"}, nil),
		"NumaOther" : prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zoneinfoSubsystem, "numa_other" ),
			"Allocation from other node",
			[]string{"node", "zone"}, nil),

	}
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
