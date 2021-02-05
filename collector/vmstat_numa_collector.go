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

// +build linux
// +build !nonumavmstat


package collector

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-kit/kit/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

const vmStatNumaSubsystem = "vmstat_numa"

type vmstatNumaCollector struct {
	metricDescs map[string]*prometheus.Desc
	logger      log.Logger
	fs          sysfs.FS
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
		metricDescs: createMetricDescriptions(),
		logger:      logger,
		fs:          fs,
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
			desc := c.metricDescs[typeOfMetricStruct.Field(i).Name]
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(metricStruct.Field(i).Uint()),
				strconv.Itoa(k))
		}

	}
	return nil
}

func createMetricDescriptions() map[string]*prometheus.Desc {
	return map[string]*prometheus.Desc{
		"NrFreePages": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_free_pages"),
			"Total number of free pages",
			[]string{"node"}, nil),
		"NrZoneInactiveAnon": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_zone_inactive_anon"),
			"Number of anonymous pages recently less used in NUMA node",
			[]string{"node"}, nil),
		"NrZoneActiveAnon": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_zone_active_anon"),
			"Number of anonymous pages recently more used in NUMA node",
			[]string{"node"}, nil),
		"NrZoneInactiveFile": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_zone_inactive_file"),
			"Number of inactive pages with file-backing in NUMA node",
			[]string{"node"}, nil),
		"NrZoneActiveFile": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_zone_active_file"),
			"Number of active pages with file-backing in NUMA node",
			[]string{"node"}, nil),
		"NrZoneUnevictable": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_zone_unevictable"),
			"Number of unevictable pages in NUMA node",
			[]string{"node"}, nil),
		"NrZoneWritePending": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_zone_write_pending"),
			"Count of dirty, writeback and unstable pages in NUMA node",
			[]string{"node"}, nil),
		"NrMlock": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_mlock"),
			"mlock()ed pages found and moved off LRU",
			[]string{"node"}, nil),
		"NrPageTablePages": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_page_table_pages"),
			"Number of pages allocated to page tables",
			[]string{"node"}, nil),
		"NrKernelStack": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_kernel_stack"),
			"Amount of memory allocated to kernel stacks",
			[]string{"node"}, nil),
		"NrBounce": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_bounce"),
			"Number of bounce pages",
			[]string{"node"}, nil),
		"NrZspages": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_zspages"),
			"Number of pages allocated in zsmalloc",
			[]string{"node"}, nil),
		"NrFreeCma": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_free_cma"),
			"Number of free cma pages",
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
		"NrInactiveAnon": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_inactive_anon"),
			"Number of anonymous pages recently less used",
			[]string{"node"}, nil),
		"NrActiveAnon": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_active_anon"),
			"Number of anonymous pages recently more used",
			[]string{"node"}, nil),
		"NrInactiveFile": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_inactive_file"),
			"Number of inactive pages with file-backing",
			[]string{"node"}, nil),
		"NrActiveFile": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_active_file"),
			"Number of active pages with file-backing",
			[]string{"node"}, nil),
		"NrUnevictable": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_unevictable"),
			"Number of unevictable pages",
			[]string{"node"}, nil),
		"NrSlabReclaimable": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_slab_reclaimable"),
			"Number of reclaimable slab pages",
			[]string{"node"}, nil),
		"NrSlabUnreclaimable": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_slab_unreclaimable"),
			"Number of unreclaimable slab pages",
			[]string{"node"}, nil),
		"NrIsolatedAnon": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_isolated_anon"),
			"Temporary isolated pages from anon lru",
			[]string{"node"}, nil),
		"NrIsolatedFile": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_isolated_file"),
			"Temporary isolated pages from file lru",
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
		"NrAnonPages": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_anon_pages"),
			"Number of anonymous pages currently used by the system",
			[]string{"node"}, nil),
		"NrMapped": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_mapped"),
			"Number of mapped pages",
			[]string{"node"}, nil),
		"NrFilePages": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_file_pages"),
			"Number of file pages",
			[]string{"node"}, nil),
		"NrDirty": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_dirty"),
			"Number of dirty pages",
			[]string{"node"}, nil),
		"NrWriteback": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_writeback"),
			"Number of writeback pages",
			[]string{"node"}, nil),
		"NrWritebackTemp": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_writeback_temp"),
			"Number of writeback pages using temporary buffers",
			[]string{"node"}, nil),
		"NrShmem": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_shmem"),
			"Number of shmem pages (included tmpfs/GEM pages)",
			[]string{"node"}, nil),
		"NrShmemHugepages": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_shmem_hugepages"),
			"Number of shmem hugepages",
			[]string{"node"}, nil),
		"NrShmemPmdmapped": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_shmem_pmdmapped"),
			"Number of pmdmapped shmem pages",
			[]string{"node"}, nil),
		"NrFilePmdmapped": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_file_pmdmapped"),
			"Number of pmdmapped pages with file-backing",
			[]string{"node"}, nil),
		"NrFileHugepages": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_file_hugepages"),
			"Number of hugepages with file-backing",
			[]string{"node"}, nil),
		"NrAnonTransparentHugepages": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_anon_transparent_hugepages"),
			"Number of anonymous transparent huge pages currently used by the system",
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
		"NrKernelMiscReclaimable": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, vmStatNumaSubsystem, "nr_kernel_misc_reclaimable"),
			"Number of reclaimable non-slab kernel pages",
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
