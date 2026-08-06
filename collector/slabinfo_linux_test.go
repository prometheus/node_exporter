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

//go:build !noslabinfo

package collector

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

const slabinfoHeader = `slabinfo - version: 2.1
# name            <active_objs> <num_objs> <objsize> <objperslab> <pagesperslab> : tunables <limit> <batchcount> <sharedfactor> : slabdata <active_slabs> <num_slabs> <sharedavail>
`

// writeSlabinfo writes body into a temporary procfs and points the collector at
// it for the duration of the test. procPath is a package-global flag value, so
// it is restored afterwards to keep these tests independent of execution order.
// These tests must not call t.Parallel() for the same reason.
func writeSlabinfo(t *testing.T, body string) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "slabinfo"), []byte(slabinfoHeader+body), 0o644); err != nil {
		t.Fatal(err)
	}
	orig := *procPath
	t.Cleanup(func() { *procPath = orig })
	*procPath = dir
}

// gatherSlabinfo collects through a registry, which is what surfaces duplicate
// label sets, and keys the result by "metric/slab/index".
func gatherSlabinfo(t *testing.T) map[string]float64 {
	t.Helper()
	c, err := NewSlabinfoCollector(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	if err != nil {
		t.Fatal(err)
	}
	reg := prometheus.NewPedanticRegistry()
	if err := reg.Register(prometheus.CollectorFunc(func(ch chan<- prometheus.Metric) {
		if err := c.Update(ch); err != nil {
			t.Errorf("Update failed: %v", err)
		}
	})); err != nil {
		t.Fatal(err)
	}
	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("Gather failed: %v", err)
	}
	got := map[string]float64{}
	for _, mf := range mfs {
		for _, m := range mf.Metric {
			var name, index string
			for _, l := range m.Label {
				switch l.GetName() {
				case "slab":
					name = l.GetValue()
				case "index":
					index = l.GetValue()
				}
			}
			got[mf.GetName()+"/"+name+"/"+index] = m.GetGauge().GetValue()
		}
	}
	return got
}

func checkSlabinfo(t *testing.T, got map[string]float64, want map[string]float64) {
	t.Helper()
	for key, w := range want {
		if g, ok := got[key]; !ok {
			t.Errorf("%s: missing", key)
		} else if g != w {
			t.Errorf("%s: got %v, want %v", key, g, w)
		}
	}
}

// A slab name may appear more than once, for example one cache per device
// instance. Each entry keeps its own series, distinguished by the index label.
func TestSlabinfoDuplicateSlabNames(t *testing.T) {
	writeSlabinfo(t, `kmem_cache           320    320    256   32    2 : tunables    0    0    0 : slabdata     10     10      0
mlx5_fs_ftes          44     44    736   44    8 : tunables    0    0    0 : slabdata      1      1      0
mlx5_fs_fgs           42     42    776   42    8 : tunables    0    0    0 : slabdata      1      1      0
mlx5_fs_ftes         396    396    736   44    8 : tunables    0    0    0 : slabdata      9      9      0
mlx5_fs_fgs           84     84    776   42    8 : tunables    0    0    0 : slabdata      2      2      0
`)
	checkSlabinfo(t, gatherSlabinfo(t), map[string]float64{
		"node_slabinfo_active_objects/mlx5_fs_ftes/0":    44,
		"node_slabinfo_active_objects/mlx5_fs_ftes/1":    396,
		"node_slabinfo_objects/mlx5_fs_ftes/1":           396,
		"node_slabinfo_active_objects/mlx5_fs_fgs/0":     42,
		"node_slabinfo_active_objects/mlx5_fs_fgs/1":     84,
		"node_slabinfo_object_size_bytes/mlx5_fs_ftes/1": 736,
		"node_slabinfo_objects_per_slab/mlx5_fs_ftes/1":  44,
		"node_slabinfo_pages_per_slab/mlx5_fs_ftes/1":    8,
		"node_slabinfo_active_objects/kmem_cache/0":      320,
	})
}

// Slabs with distinct names each get index 0.
func TestSlabinfoDistinctSlabNames(t *testing.T) {
	writeSlabinfo(t, `tw_sock_TCP          704    864    256   32    2 : tunables    0    0    0 : slabdata     27     27      0
kmem_cache           320    320    256   32    2 : tunables    0    0    0 : slabdata     10     10      0
`)
	checkSlabinfo(t, gatherSlabinfo(t), map[string]float64{
		"node_slabinfo_active_objects/tw_sock_TCP/0": 704,
		"node_slabinfo_objects/tw_sock_TCP/0":        864,
		"node_slabinfo_active_objects/kmem_cache/0":  320,
		"node_slabinfo_objects/kmem_cache/0":         320,
	})
}

// The kernel does not guarantee that same-named caches share geometry: the
// duplicate-name check in kmem_cache_sanity_check() only WARNs and is compiled
// out without CONFIG_DEBUG_VM, and object_size is per-cache. Indexing keeps each
// cache's own geometry, so no representative value has to be chosen.
func TestSlabinfoDuplicateSlabNamesDifferentGeometry(t *testing.T) {
	writeSlabinfo(t, `weird_cache          10     10    128   32    2 : tunables    0    0    0 : slabdata      1      1      0
weird_cache          20     20    256   16    4 : tunables    0    0    0 : slabdata      2      2      0
`)
	checkSlabinfo(t, gatherSlabinfo(t), map[string]float64{
		"node_slabinfo_active_objects/weird_cache/0":    10,
		"node_slabinfo_object_size_bytes/weird_cache/0": 128,
		"node_slabinfo_objects_per_slab/weird_cache/0":  32,
		"node_slabinfo_pages_per_slab/weird_cache/0":    2,
		"node_slabinfo_active_objects/weird_cache/1":    20,
		"node_slabinfo_object_size_bytes/weird_cache/1": 256,
		"node_slabinfo_objects_per_slab/weird_cache/1":  16,
		"node_slabinfo_pages_per_slab/weird_cache/1":    4,
	})
}

// Excluded slabs must not consume an index, so the ordinals of the remaining
// entries stay contiguous.
func TestSlabinfoIndexIgnoresFilteredSlabs(t *testing.T) {
	orig := *slabNameExclude
	t.Cleanup(func() { *slabNameExclude = orig })
	*slabNameExclude = "^drop_me$"

	writeSlabinfo(t, `dup_cache            10     10    128   32    2 : tunables    0    0    0 : slabdata      1      1      0
drop_me              99     99    999   99    9 : tunables    0    0    0 : slabdata      9      9      0
dup_cache            20     20    128   32    2 : tunables    0    0    0 : slabdata      2      2      0
`)
	got := gatherSlabinfo(t)
	checkSlabinfo(t, got, map[string]float64{
		"node_slabinfo_active_objects/dup_cache/0": 10,
		"node_slabinfo_active_objects/dup_cache/1": 20,
	})
	if _, ok := got["node_slabinfo_active_objects/drop_me/0"]; ok {
		t.Error("excluded slab was collected")
	}
}
