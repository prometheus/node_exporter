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

// +build !nomeminfo

package collector

import (
	"reflect"
	"testing"
)

// var (
// 	reParens = regexp.MustCompile(`\((.*)\)`)
// )

func TestGetMemInfo(t *testing.T) {
	info := make(map[string]float64)
	swapInfo := map[string]float64{
		"total":       1024454656,
		"used":        1060864,
		"free":        1023393792,
		"usedPercent": 0.10355402201422569,
		"sin":         36864,
		"sout":        81920,
		"pgin":        27170140160,
		"pgout":       210821484544,
		"pgfault":     10854351642624,
	}

	virtualInfo := map[string]float64{
		"total":          16705363968,
		"available":      5313527808,
		"used":           10251444224,
		"usedPercent":    61.366183003478284,
		"free":           1512931328,
		"active":         11915513856,
		"inactive":       2571345920,
		"wired":          0,
		"laundry":        0,
		"buffers":        459743232,
		"cached":         4481245184,
		"writeback":      106496,
		"dirty":          892928,
		"writebacktmp":   0,
		"shared":         909332480,
		"slab":           454131712,
		"sreclaimable":   318140416,
		"sunreclaim":     135991296,
		"pagetables":     113631232,
		"swapcached":     8192,
		"commitlimit":    9377136640,
		"committedas":    23042527232,
		"hightotal":      0,
		"highfree":       0,
		"lowtotal":       0,
		"lowfree":        0,
		"swaptotal":      1024454656,
		"swapfree":       1023393792,
		"mapped":         1528872960,
		"vmalloctotal":   35184372087808,
		"vmallocused":    0,
		"vmallocchunk":   0,
		"hugepagestotal": 0,
		"hugepagesfree":  0,
		"hugepagesize":   2097152,
	}

	mergeInfo(swapInfo, info, "swap")
	mergeInfo(virtualInfo, info, "virtual")

	expectedInfo := map[string]float64{
		"swap_free_bytes":              1.023393792e+09,
		"swap_pgfault_bytes":           1.0854351642624e+13,
		"swap_pgin_bytes":              2.717014016e+10,
		"swap_pgout_bytes":             2.10821484544e+11,
		"swap_sin_bytes":               36864,
		"swap_sout_bytes":              81920,
		"swap_total_bytes":             1.024454656e+09,
		"swap_used_bytes":              1.060864e+06,
		"swap_used_percent":            0.10355402201422569,
		"virtual_active_bytes":         1.1915513856e+10,
		"virtual_available_bytes":      5.313527808e+09,
		"virtual_buffers_bytes":        4.59743232e+08,
		"virtual_cached_bytes":         4.481245184e+09,
		"virtual_commitlimit_bytes":    9.37713664e+09,
		"virtual_committedas_bytes":    2.3042527232e+10,
		"virtual_dirty_bytes":          892928,
		"virtual_free_bytes":           1.512931328e+09,
		"virtual_highfree_bytes":       0,
		"virtual_hightotal_bytes":      0,
		"virtual_hugepagesfree_bytes":  0,
		"virtual_hugepagesize_bytes":   2.097152e+06,
		"virtual_hugepagestotal_bytes": 0,
		"virtual_inactive_bytes":       2.57134592e+09,
		"virtual_laundry_bytes":        0,
		"virtual_lowfree_bytes":        0,
		"virtual_lowtotal_bytes":       0,
		"virtual_mapped_bytes":         1.52887296e+09,
		"virtual_pagetables_bytes":     1.13631232e+08,
		"virtual_shared_bytes":         9.0933248e+08,
		"virtual_slab_bytes":           4.54131712e+08,
		"virtual_sreclaimable_bytes":   3.18140416e+08,
		"virtual_sunreclaim_bytes":     1.35991296e+08,
		"virtual_swapcached_bytes":     8192,
		"virtual_swapfree_bytes":       1.023393792e+09,
		"virtual_swaptotal_bytes":      1.024454656e+09,
		"virtual_total_bytes":          1.6705363968e+10,
		"virtual_used_bytes":           1.0251444224e+10,
		"virtual_used_percent":         61.366183003478284,
		"virtual_vmallocchunk_bytes":   0,
		"virtual_vmalloctotal_bytes":   3.5184372087808e+13,
		"virtual_vmallocused_bytes":    0,
		"virtual_wired_bytes":          0,
		"virtual_writeback_bytes":      106496,
		"virtual_writebacktmp_bytes":   0,
	}
	if !reflect.DeepEqual(info, expectedInfo) {
		t.Error("want: ", expectedInfo, "but got: ", info)
	}
}
