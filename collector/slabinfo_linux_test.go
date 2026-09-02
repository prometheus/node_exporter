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
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

const slabinfoHeader = `slabinfo - version: 2.1
# name            <active_objs> <num_objs> <objsize> <objperslab> <pagesperslab> : tunables <limit> <batchcount> <sharedfactor> : slabdata <active_slabs> <num_slabs> <sharedavail>
`

// setProcPath points the collector at dir for the duration of the test. procPath
// is a package-global flag value, so it is restored afterwards to keep these
// tests independent of execution order. These tests must not call t.Parallel()
// for the same reason.
func setProcPath(t *testing.T, dir string) {
	t.Helper()
	orig := *procPath
	t.Cleanup(func() { *procPath = orig })
	*procPath = dir
}

// writeSlabinfo writes body into a temporary procfs and points the collector at
// it for the duration of the test.
func writeSlabinfo(t *testing.T, body string) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "slabinfo"), []byte(slabinfoHeader+body), 0o644); err != nil {
		t.Fatal(err)
	}
	setProcPath(t, dir)
}

// slabinfoTestCollector adapts the collector to prometheus.Collector so the
// metrics can be gathered through a registry, which is what surfaces duplicate
// label sets.
//
// Describe is deliberately empty, which registers this as an unchecked
// collector. Because the index label is conditional, one metric family can now
// carry both the {slab} and the {slab,index} descriptor, and Register rejects
// two descriptors that share a fully-qualified name but declare different label
// names:
//
//	descriptors reported by collector have inconsistent label names or help
//	strings for the same fully-qualified name
//
// Gathering mixed label sets is fine; only the Describe-time check objects.
// Production never reaches it because NodeCollector.Describe reports just the
// two scrape descriptors. Do not "fix" this by describing through Collect
// (prometheus.DescribeByCollect or prometheus.CollectorFunc): that reintroduces
// the registration failure.
type slabinfoTestCollector struct {
	t *testing.T
	c Collector
}

func (slabinfoTestCollector) Describe(chan<- *prometheus.Desc) {}

func (tc slabinfoTestCollector) Collect(ch chan<- prometheus.Metric) {
	if err := tc.c.Update(ch); err != nil {
		tc.t.Errorf("Update failed: %v", err)
	}
}

// expectSlabinfo gathers the collector and compares the complete exposition, so
// extra series, missing series, wrong labels, wrong values and changed HELP or
// TYPE lines all fail.
func expectSlabinfo(t *testing.T, want string) {
	t.Helper()
	c, err := NewSlabinfoCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(slabinfoTestCollector{t: t, c: c})

	if err := testutil.GatherAndCompare(reg, strings.NewReader(want)); err != nil {
		t.Error(err)
	}
}

// Slab names that occur once must keep the exact identity they had before
// duplicate handling existed: bare {slab}, with no index label anywhere. This is
// the guarantee the whole change rests on, mirroring
// TestHwmonUniqueChipNamesAreUnchanged.
func TestSlabinfoUniqueSlabNamesAreUnchanged(t *testing.T) {
	writeSlabinfo(t, `tw_sock_TCP          704    864    256   32    2 : tunables    0    0    0 : slabdata     27     27      0
kmem_cache           320    320    256   32    2 : tunables    0    0    0 : slabdata     10     10      0
`)
	expectSlabinfo(t, `# HELP node_slabinfo_active_objects The number of objects that are currently active (i.e., in use).
# TYPE node_slabinfo_active_objects gauge
node_slabinfo_active_objects{slab="kmem_cache"} 320
node_slabinfo_active_objects{slab="tw_sock_TCP"} 704
# HELP node_slabinfo_object_size_bytes The size of objects in this slab, in bytes.
# TYPE node_slabinfo_object_size_bytes gauge
node_slabinfo_object_size_bytes{slab="kmem_cache"} 256
node_slabinfo_object_size_bytes{slab="tw_sock_TCP"} 256
# HELP node_slabinfo_objects The total number of allocated objects (i.e., objects that are both in use and not in use).
# TYPE node_slabinfo_objects gauge
node_slabinfo_objects{slab="kmem_cache"} 320
node_slabinfo_objects{slab="tw_sock_TCP"} 864
# HELP node_slabinfo_objects_per_slab The number of objects stored in each slab.
# TYPE node_slabinfo_objects_per_slab gauge
node_slabinfo_objects_per_slab{slab="kmem_cache"} 32
node_slabinfo_objects_per_slab{slab="tw_sock_TCP"} 32
# HELP node_slabinfo_pages_per_slab The number of pages allocated for each slab.
# TYPE node_slabinfo_pages_per_slab gauge
node_slabinfo_pages_per_slab{slab="kmem_cache"} 2
node_slabinfo_pages_per_slab{slab="tw_sock_TCP"} 2
`)
}

// A slab name may appear more than once, for example one cache per device
// instance. Each entry then keeps its own series, distinguished by an ordinal in
// /proc/slabinfo order.
func TestSlabinfoDuplicateSlabNames(t *testing.T) {
	writeSlabinfo(t, `mlx5_fs_ftes          44     44    736   44    8 : tunables    0    0    0 : slabdata      1      1      0
mlx5_fs_ftes         396    396    736   44    8 : tunables    0    0    0 : slabdata      9      9      0
mlx5_fs_ftes         792    792    736   44    8 : tunables    0    0    0 : slabdata     18     18      0
`)
	expectSlabinfo(t, `# HELP node_slabinfo_active_objects The number of objects that are currently active (i.e., in use).
# TYPE node_slabinfo_active_objects gauge
node_slabinfo_active_objects{index="0",slab="mlx5_fs_ftes"} 44
node_slabinfo_active_objects{index="1",slab="mlx5_fs_ftes"} 396
node_slabinfo_active_objects{index="2",slab="mlx5_fs_ftes"} 792
# HELP node_slabinfo_object_size_bytes The size of objects in this slab, in bytes.
# TYPE node_slabinfo_object_size_bytes gauge
node_slabinfo_object_size_bytes{index="0",slab="mlx5_fs_ftes"} 736
node_slabinfo_object_size_bytes{index="1",slab="mlx5_fs_ftes"} 736
node_slabinfo_object_size_bytes{index="2",slab="mlx5_fs_ftes"} 736
# HELP node_slabinfo_objects The total number of allocated objects (i.e., objects that are both in use and not in use).
# TYPE node_slabinfo_objects gauge
node_slabinfo_objects{index="0",slab="mlx5_fs_ftes"} 44
node_slabinfo_objects{index="1",slab="mlx5_fs_ftes"} 396
node_slabinfo_objects{index="2",slab="mlx5_fs_ftes"} 792
# HELP node_slabinfo_objects_per_slab The number of objects stored in each slab.
# TYPE node_slabinfo_objects_per_slab gauge
node_slabinfo_objects_per_slab{index="0",slab="mlx5_fs_ftes"} 44
node_slabinfo_objects_per_slab{index="1",slab="mlx5_fs_ftes"} 44
node_slabinfo_objects_per_slab{index="2",slab="mlx5_fs_ftes"} 44
# HELP node_slabinfo_pages_per_slab The number of pages allocated for each slab.
# TYPE node_slabinfo_pages_per_slab gauge
node_slabinfo_pages_per_slab{index="0",slab="mlx5_fs_ftes"} 8
node_slabinfo_pages_per_slab{index="1",slab="mlx5_fs_ftes"} 8
node_slabinfo_pages_per_slab{index="2",slab="mlx5_fs_ftes"} 8
`)
}

// Unique and duplicated names coexist inside the same metric family, one carrying
// the index label and one not. This is the mixed-label-set case: a registry
// accepts it at gather time even though the family holds two label shapes.
func TestSlabinfoMixedUniqueAndDuplicateSlabNames(t *testing.T) {
	writeSlabinfo(t, `kmem_cache           320    320    256   32    2 : tunables    0    0    0 : slabdata     10     10      0
mlx5_fs_fgs           42     42    776   42    8 : tunables    0    0    0 : slabdata      1      1      0
mlx5_fs_fgs           84     84    776   42    8 : tunables    0    0    0 : slabdata      2      2      0
tw_sock_TCP          704    864    256   32    2 : tunables    0    0    0 : slabdata     27     27      0
`)
	expectSlabinfo(t, `# HELP node_slabinfo_active_objects The number of objects that are currently active (i.e., in use).
# TYPE node_slabinfo_active_objects gauge
node_slabinfo_active_objects{slab="kmem_cache"} 320
node_slabinfo_active_objects{index="0",slab="mlx5_fs_fgs"} 42
node_slabinfo_active_objects{index="1",slab="mlx5_fs_fgs"} 84
node_slabinfo_active_objects{slab="tw_sock_TCP"} 704
# HELP node_slabinfo_object_size_bytes The size of objects in this slab, in bytes.
# TYPE node_slabinfo_object_size_bytes gauge
node_slabinfo_object_size_bytes{slab="kmem_cache"} 256
node_slabinfo_object_size_bytes{index="0",slab="mlx5_fs_fgs"} 776
node_slabinfo_object_size_bytes{index="1",slab="mlx5_fs_fgs"} 776
node_slabinfo_object_size_bytes{slab="tw_sock_TCP"} 256
# HELP node_slabinfo_objects The total number of allocated objects (i.e., objects that are both in use and not in use).
# TYPE node_slabinfo_objects gauge
node_slabinfo_objects{slab="kmem_cache"} 320
node_slabinfo_objects{index="0",slab="mlx5_fs_fgs"} 42
node_slabinfo_objects{index="1",slab="mlx5_fs_fgs"} 84
node_slabinfo_objects{slab="tw_sock_TCP"} 864
# HELP node_slabinfo_objects_per_slab The number of objects stored in each slab.
# TYPE node_slabinfo_objects_per_slab gauge
node_slabinfo_objects_per_slab{slab="kmem_cache"} 32
node_slabinfo_objects_per_slab{index="0",slab="mlx5_fs_fgs"} 42
node_slabinfo_objects_per_slab{index="1",slab="mlx5_fs_fgs"} 42
node_slabinfo_objects_per_slab{slab="tw_sock_TCP"} 32
# HELP node_slabinfo_pages_per_slab The number of pages allocated for each slab.
# TYPE node_slabinfo_pages_per_slab gauge
node_slabinfo_pages_per_slab{slab="kmem_cache"} 2
node_slabinfo_pages_per_slab{index="0",slab="mlx5_fs_fgs"} 8
node_slabinfo_pages_per_slab{index="1",slab="mlx5_fs_fgs"} 8
node_slabinfo_pages_per_slab{slab="tw_sock_TCP"} 2
`)
}

// The kernel does not guarantee that same-named caches share geometry: the
// duplicate-name check in kmem_cache_sanity_check() only WARNs and is compiled
// out without CONFIG_DEBUG_VM, and object_size is per-cache. Indexing keeps each
// cache's own geometry, so no representative value has to be chosen.
func TestSlabinfoDuplicateSlabNamesDifferentGeometry(t *testing.T) {
	writeSlabinfo(t, `weird_cache          10     10    128   32    2 : tunables    0    0    0 : slabdata      1      1      0
weird_cache          20     20    256   16    4 : tunables    0    0    0 : slabdata      2      2      0
`)
	expectSlabinfo(t, `# HELP node_slabinfo_active_objects The number of objects that are currently active (i.e., in use).
# TYPE node_slabinfo_active_objects gauge
node_slabinfo_active_objects{index="0",slab="weird_cache"} 10
node_slabinfo_active_objects{index="1",slab="weird_cache"} 20
# HELP node_slabinfo_object_size_bytes The size of objects in this slab, in bytes.
# TYPE node_slabinfo_object_size_bytes gauge
node_slabinfo_object_size_bytes{index="0",slab="weird_cache"} 128
node_slabinfo_object_size_bytes{index="1",slab="weird_cache"} 256
# HELP node_slabinfo_objects The total number of allocated objects (i.e., objects that are both in use and not in use).
# TYPE node_slabinfo_objects gauge
node_slabinfo_objects{index="0",slab="weird_cache"} 10
node_slabinfo_objects{index="1",slab="weird_cache"} 20
# HELP node_slabinfo_objects_per_slab The number of objects stored in each slab.
# TYPE node_slabinfo_objects_per_slab gauge
node_slabinfo_objects_per_slab{index="0",slab="weird_cache"} 32
node_slabinfo_objects_per_slab{index="1",slab="weird_cache"} 16
# HELP node_slabinfo_pages_per_slab The number of pages allocated for each slab.
# TYPE node_slabinfo_pages_per_slab gauge
node_slabinfo_pages_per_slab{index="0",slab="weird_cache"} 2
node_slabinfo_pages_per_slab{index="1",slab="weird_cache"} 4
`)
}

// The filter matches on slab name, so every entry sharing a name is kept or
// dropped together. An excluded name contributes no series at all, and excluding
// it leaves the collision handling of the names that remain untouched: a
// duplicated name is still indexed and a unique name is still bare.
func TestSlabinfoExcludedSlabsAreNotCollected(t *testing.T) {
	origExclude := *slabNameExclude
	t.Cleanup(func() { *slabNameExclude = origExclude })
	*slabNameExclude = "^drop_me$"

	writeSlabinfo(t, `keep_me              10     10    128   32    2 : tunables    0    0    0 : slabdata      1      1      0
drop_me              99     99    999   99    9 : tunables    0    0    0 : slabdata      9      9      0
drop_me              98     98    999   99    9 : tunables    0    0    0 : slabdata      8      8      0
dup_cache            20     20    128   32    2 : tunables    0    0    0 : slabdata      2      2      0
dup_cache            30     30    128   32    2 : tunables    0    0    0 : slabdata      3      3      0
`)
	expectSlabinfo(t, `# HELP node_slabinfo_active_objects The number of objects that are currently active (i.e., in use).
# TYPE node_slabinfo_active_objects gauge
node_slabinfo_active_objects{index="0",slab="dup_cache"} 20
node_slabinfo_active_objects{index="1",slab="dup_cache"} 30
node_slabinfo_active_objects{slab="keep_me"} 10
# HELP node_slabinfo_object_size_bytes The size of objects in this slab, in bytes.
# TYPE node_slabinfo_object_size_bytes gauge
node_slabinfo_object_size_bytes{index="0",slab="dup_cache"} 128
node_slabinfo_object_size_bytes{index="1",slab="dup_cache"} 128
node_slabinfo_object_size_bytes{slab="keep_me"} 128
# HELP node_slabinfo_objects The total number of allocated objects (i.e., objects that are both in use and not in use).
# TYPE node_slabinfo_objects gauge
node_slabinfo_objects{index="0",slab="dup_cache"} 20
node_slabinfo_objects{index="1",slab="dup_cache"} 30
node_slabinfo_objects{slab="keep_me"} 10
# HELP node_slabinfo_objects_per_slab The number of objects stored in each slab.
# TYPE node_slabinfo_objects_per_slab gauge
node_slabinfo_objects_per_slab{index="0",slab="dup_cache"} 32
node_slabinfo_objects_per_slab{index="1",slab="dup_cache"} 32
node_slabinfo_objects_per_slab{slab="keep_me"} 32
# HELP node_slabinfo_pages_per_slab The number of pages allocated for each slab.
# TYPE node_slabinfo_pages_per_slab gauge
node_slabinfo_pages_per_slab{index="0",slab="dup_cache"} 2
node_slabinfo_pages_per_slab{index="1",slab="dup_cache"} 2
node_slabinfo_pages_per_slab{slab="keep_me"} 2
`)
}

// An unreadable /proc/slabinfo must fail the collector rather than report a
// partial scrape.
func TestSlabinfoNoSlabinfoFile(t *testing.T) {
	setProcPath(t, t.TempDir())

	c, err := NewSlabinfoCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan prometheus.Metric, 1)
	if err := c.Update(ch); err == nil {
		t.Fatal("expected an error when /proc/slabinfo is missing, got nil")
	}
	close(ch)
	if n := len(ch); n != 0 {
		t.Errorf("expected no metrics, got %d", n)
	}
}
