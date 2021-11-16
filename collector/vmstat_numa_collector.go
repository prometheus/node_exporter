// Copyright 2021 The Prometheus Authors
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

//go:build linux && !nonumavmstat
// +build linux,!nonumavmstat

package collector

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-kit/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

const vmStatNumaSubsystem = "vmstat_numa"

type aggregatedMetricPair struct {
	baseMetric, label string
}

type vmstatNumaCollector struct {
	metricDescs         map[string]*prometheus.Desc
	aggregatedMetricMap map[string]aggregatedMetricPair
	logger              log.Logger
	fs                  sysfs.FS
}

func init() {
	registerCollector("vmstat_numa", defaultDisabled, NewVmstatNumaCollector)
}

// NewVmstatNumaCollector returns a new Collector exposing memory stats.
// Returns an error when cAdvisor can't read procfs.
func NewVmstatNumaCollector(logger log.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}
	return &vmstatNumaCollector{
		metricDescs:         createMetricDescriptions(),
		aggregatedMetricMap: createAggregatedMetricMap(),
		logger:              logger,
		fs:                  fs,
	}, nil
}

func (c *vmstatNumaCollector) Update(ch chan<- prometheus.Metric) error {
	metrics, err := c.fs.VMStatNUMA()
	if err != nil {
		return fmt.Errorf("couldn't get NUMA vmstat: %w", err)
	}
	for k, v := range metrics {
		metricStruct := reflect.ValueOf(v)
		typeOfMetricStruct := metricStruct.Type()
		for i := 0; i < metricStruct.NumField(); i++ {
			if aggregatedMetric, ok := c.aggregatedMetricMap[typeOfMetricStruct.Field(i).Name]; ok {
				desc := c.metricDescs[aggregatedMetric.baseMetric]
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(metricStruct.Field(i).Uint()),
					strconv.Itoa(k), aggregatedMetric.label)
			} else {
				desc := c.metricDescs[typeOfMetricStruct.Field(i).Name]
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(metricStruct.Field(i).Uint()),
					strconv.Itoa(k))

			}
		}

	}
	return nil
}

func createAggregatedMetricMap() map[string]aggregatedMetricPair {
	return map[string]aggregatedMetricPair{
		"NrFreePages":                aggregatedMetricPair{"NrPages", "free"},
		"NrZoneInactiveAnon":         aggregatedMetricPair{"NrPages", "zone_inactive_anon"},
		"NrZoneActiveAnon":           aggregatedMetricPair{"NrPages", "zone_active_anon"},
		"NrZoneInactiveFile":         aggregatedMetricPair{"NrPages", "zone_inactive_file"},
		"NrZoneActiveFile":           aggregatedMetricPair{"NrPages", "zone_active_file"},
		"NrZoneUnevictable":          aggregatedMetricPair{"NrPages", "zone_unevictable"},
		"NrZoneWritePending":         aggregatedMetricPair{"NrPages", "zone_write_pending"},
		"NrMlock":                    aggregatedMetricPair{"NrPages", "mlock"},
		"NrPageTablePages":           aggregatedMetricPair{"NrPages", "page_table_pages"},
		"NrBounce":                   aggregatedMetricPair{"NrPages", "bounce"},
		"NrZspages":                  aggregatedMetricPair{"NrPages", "zspages"},
		"NrFreeCma":                  aggregatedMetricPair{"NrPages", "free_cma"},
		"NrInactiveAnon":             aggregatedMetricPair{"NrPages", "inactive_anon"},
		"NrActiveAnon":               aggregatedMetricPair{"NrPages", "active_anon"},
		"NrInactiveFile":             aggregatedMetricPair{"NrPages", "inactive_file"},
		"NrActiveFile":               aggregatedMetricPair{"NrPages", "active_file"},
		"NrUnevictable":              aggregatedMetricPair{"NrPages", "unevictable"},
		"NrSlabReclaimable":          aggregatedMetricPair{"NrPages", "slab_reclaimable"},
		"NrSlabUnreclaimable":        aggregatedMetricPair{"NrPages", "slab_unreclaimable"},
		"NrIsolatedAnon":             aggregatedMetricPair{"NrPages", "isolated_anon"},
		"NrIsolatedFile":             aggregatedMetricPair{"NrPages", "isolated_file"},
		"NrAnonPages":                aggregatedMetricPair{"NrPages", "anon_pages"},
		"NrMapped":                   aggregatedMetricPair{"NrPages", "mapped"},
		"NrFilePages":                aggregatedMetricPair{"NrPages", "file_pages"},
		"NrDirty":                    aggregatedMetricPair{"NrPages", "dirty"},
		"NrWriteback":                aggregatedMetricPair{"NrPages", "writeback"},
		"NrWritebackTemp":            aggregatedMetricPair{"NrPages", "writeback_temp"},
		"NrShmem":                    aggregatedMetricPair{"NrPages", "shmem"},
		"NrShmemHugepages":           aggregatedMetricPair{"NrPages", "shmem_hugepages"},
		"NrShmemPmdmapped":           aggregatedMetricPair{"NrPages", "shmem_pmdmapped"},
		"NrFilePmdmapped":            aggregatedMetricPair{"NrPages", "file_pmdmapped"},
		"NrFileHugepages":            aggregatedMetricPair{"NrPages", "file_hugepages"},
		"NrAnonTransparentHugepages": aggregatedMetricPair{"NrPages", "anon_transparent_hugepages"},
		"NrKernelMiscReclaimable":    aggregatedMetricPair{"NrPages", "kernel_misc_reclaimable"},
		"NrFollPinAcquired":          aggregatedMetricPair{"NrPages", "foll_pin_acquired"},
		"NrFollPinReleased":          aggregatedMetricPair{"NrPages", "foll_pin_released"},
	}

}

func createMetricDescriptions() map[string]*prometheus.Desc {
	return map[string]*prometheus.Desc{
		"NrPages": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_pages"),
			"Number of pages",
			[]string{"node", "type"}, nil),
		"NrKernelStack": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_kernel_stack"),
			"Amount of memory allocated to kernel stacks",
			[]string{"node"}, nil),
		"NumaHit": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "numa_hit"),
			"Allocated in intended node",
			[]string{"node"}, nil),
		"NumaMiss": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "numa_miss"),
			"Allocated in non intended node",
			[]string{"node"}, nil),
		"NumaForeign": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "numa_foreign"),
			"Was intended here, hit elsewhere",
			[]string{"node"}, nil),
		"NumaInterleave": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "numa_interleave"),
			"Interleaver preferred this zone",
			[]string{"node"}, nil),
		"NumaLocal": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "numa_local"),
			"Allocation from local node",
			[]string{"node"}, nil),
		"NumaOther": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "numa_other"),
			"Allocation from other node",
			[]string{"node"}, nil),
		"WorkingsetNodes": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "workingset_nodes"),
			"Number of nodes in working set",
			[]string{"node"}, nil),
		"WorkingsetRefault": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "workingset_refault"),
			"Number of refaults of previously evicted pages",
			[]string{"node"}, nil),
		"WorkingsetActivate": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "workingset_activate"),
			"Number of refaulted pages that were immediately activated",
			[]string{"node"}, nil),
		"WorkingsetRestore": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "workingset_restore"),
			"Number of restored pages which have been detected as an active",
			[]string{"node"}, nil),
		"WorkingsetNodereclaim": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "workingset_nodereclaim"),
			"Number of times a shadow node has been reclaimed",
			[]string{"node"}, nil),
		"NrVmscanWrite": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_vmscan_write"),
			"Number of writebacks of dirty pages during scan of LRU",
			[]string{"node"}, nil),
		"NrVmscanImmediateReclaim": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_vmscan_immediate_reclaim"),
			"Prioritise for reclaim when writeback ends",
			[]string{"node"}, nil),
		"NrDirtied": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_dirtied"),
			"Page dirtyings since bootup",
			[]string{"node"}, nil),
		"NrWritten": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_written"),
			"Page writings since bootup",
			[]string{"node"}, nil),
		"NrFollPinAcquired": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_foll_pin_acquired"),
			"Number of pages allocated via: pin_user_page(), gup flag: FOLL_PIN",
			[]string{"node"}, nil),
		"NrFollPinReleased": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_foll_pin_released"),
			"Number of pages returned via unpin_user_page()",
			[]string{"node"}, nil),
	}
}
