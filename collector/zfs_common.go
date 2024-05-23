// Copyright 2024 The Prometheus Authors
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

var (
	bsdsysctl = map[string]string{
		"kstat.zfs.misc.abdstats.linear_cnt":             "abdstats_linear_count_total",
		"kstat.zfs.misc.abdstats.linear_data_size":       "abdstats_linear_data_bytes",
		"kstat.zfs.misc.abdstats.scatter_chunk_waste":    "abdstats_scatter_chunk_waste_bytes",
		"kstat.zfs.misc.abdstats.scatter_cnt":            "abdstats_scatter_count_total",
		"kstat.zfs.misc.abdstats.scatter_data_size":      "abdstats_scatter_data_bytes",
		"kstat.zfs.misc.abdstats.struct_size":            "abdstats_struct_bytes",
		"kstat.zfs.misc.arcstats.anon_size":              "arcstats_anon_bytes",
		"kstat.zfs.misc.arcstats.c":                      "arcstats_c_bytes",
		"kstat.zfs.misc.arcstats.c_max":                  "arcstats_c_max_bytes",
		"kstat.zfs.misc.arcstats.c_min":                  "arcstats_c_min_bytes",
		"kstat.zfs.misc.arcstats.data_size":              "arcstats_data_bytes",
		"kstat.zfs.misc.arcstats.demand_data_hits":       "arcstats_demand_data_hits_total",
		"kstat.zfs.misc.arcstats.demand_data_misses":     "arcstats_demand_data_misses_total",
		"kstat.zfs.misc.arcstats.demand_metadata_hits":   "arcstats_demand_metadata_hits_total",
		"kstat.zfs.misc.arcstats.demand_metadata_misses": "arcstats_demand_metadata_misses_total",
		"kstat.zfs.misc.arcstats.hdr_size":               "arcstats_hdr_bytes",
		"kstat.zfs.misc.arcstats.hits":                   "arcstats_hits_total",
		"kstat.zfs.misc.arcstats.misses":                 "arcstats_misses_total",
		"kstat.zfs.misc.arcstats.mfu_ghost_hits":         "arcstats_mfu_ghost_hits_total",
		"kstat.zfs.misc.arcstats.mfu_ghost_size":         "arcstats_mfu_ghost_size",
		"kstat.zfs.misc.arcstats.mfu_size":               "arcstats_mfu_bytes",
		"kstat.zfs.misc.arcstats.mru_ghost_hits":         "arcstats_mru_ghost_hits_total",
		"kstat.zfs.misc.arcstats.mru_ghost_size":         "arcstats_mru_ghost_bytes",
		"kstat.zfs.misc.arcstats.mru_size":               "arcstats_mru_bytes",
		"kstat.zfs.misc.arcstats.other_size":             "arcstats_other_bytes",
		"kstat.zfs.misc.arcstats.p":                      "arcstats_p_bytes",
		"kstat.zfs.misc.arcstats.meta":                   "arcstats_meta_bytes",
		"kstat.zfs.misc.arcstats.pd":                     "arcstats_pd_bytes",
		"kstat.zfs.misc.arcstats.pm":                     "arcstats_pm_bytes",
		"kstat.zfs.misc.arcstats.size":                   "arcstats_size_bytes",
		"kstat.zfs.misc.zfetchstats.hits":                "zfetchstats_hits_total",
		"kstat.zfs.misc.zfetchstats.misses":              "zfetchstats_misses_total",
	}
)
