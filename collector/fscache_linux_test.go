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

//go:build !nofscache
// +build !nofscache

package collector

import (
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestFscacheStats(t *testing.T) {
	testcases := []struct {
		name     string
		procPath string
		expected string
		err      error
	}{
		{
			name:     "stats format",
			procPath: "fixtures/proc",
			expected: `
# HELP node_fscache_acquire_attempts_total Number of acquire operations attempted (Acquire: n=attempts).
# TYPE node_fscache_acquire_attempts_total counter
node_fscache_acquire_attempts_total 31998
# HELP node_fscache_acquire_success_total Number of acquire operations successful (Acquire: ok=success).
# TYPE node_fscache_acquire_success_total counter
node_fscache_acquire_success_total 31986
# HELP node_fscache_allocations_success_total Number of successful allocation operations (Allocs: ok=success).
# TYPE node_fscache_allocations_success_total counter
node_fscache_allocations_success_total 0
# HELP node_fscache_allocations_total Number of allocation operations attempted (Allocs: n=attempts).
# TYPE node_fscache_allocations_total counter
node_fscache_allocations_total 0
# HELP node_fscache_attribute_changes_success_total Number of successful attribute change operations (AttrChg: ok=success).
# TYPE node_fscache_attribute_changes_success_total counter
node_fscache_attribute_changes_success_total 0
# HELP node_fscache_attribute_changes_total Number of attribute change operations attempted (AttrChg: n=attempts).
# TYPE node_fscache_attribute_changes_total counter
node_fscache_attribute_changes_total 0
# HELP node_fscache_invalidations_total Number of invalidation operations (Invals: n=tot).
# TYPE node_fscache_invalidations_total counter
node_fscache_invalidations_total 409
# HELP node_fscache_lookups_negative_total Number of negative lookup operations (Lookups: neg=negative).
# TYPE node_fscache_lookups_negative_total counter
node_fscache_lookups_negative_total 0
# HELP node_fscache_lookups_positive_total Number of positive lookup operations (Lookups: pos=positive).
# TYPE node_fscache_lookups_positive_total counter
node_fscache_lookups_positive_total 0
# HELP node_fscache_lookups_total Number of lookup operations (Lookups: n=tot).
# TYPE node_fscache_lookups_total counter
node_fscache_lookups_total 0
# HELP node_fscache_objects_allocated_total Number of index cookies allocated (Cookies: idx=allocated/available/unused).
# TYPE node_fscache_objects_allocated_total counter
node_fscache_objects_allocated_total 16
# HELP node_fscache_objects_available_total Number of index cookies available (Cookies: idx=allocated/available/unused).
# TYPE node_fscache_objects_available_total counter
node_fscache_objects_available_total 0
# HELP node_fscache_objects_culled_total Number of objects culled (CacheEv: cul=culled).
# TYPE node_fscache_objects_culled_total counter
node_fscache_objects_culled_total 0
# HELP node_fscache_objects_retired_total Number of objects retired (CacheEv: rtr=retired).
# TYPE node_fscache_objects_retired_total counter
node_fscache_objects_retired_total 0
# HELP node_fscache_relinquishes_total Number of relinquish operations (Relinqs: n=tot).
# TYPE node_fscache_relinquishes_total counter
node_fscache_relinquishes_total 31939
# HELP node_fscache_retrievals_nobuffer_total Number of retrieval (read) operations failed due to no buffer (Retrvls: nbf=nobuff).
# TYPE node_fscache_retrievals_nobuffer_total counter
node_fscache_retrievals_nobuffer_total 2551742
# HELP node_fscache_retrievals_success_total Number of successful retrieval (read) operations (Retrvls: ok=success).
# TYPE node_fscache_retrievals_success_total counter
node_fscache_retrievals_success_total 0
# HELP node_fscache_retrievals_total Number of retrieval (read) operations attempted (Retrvls: n=attempts).
# TYPE node_fscache_retrievals_total counter
node_fscache_retrievals_total 2551742
# HELP node_fscache_stores_success_total Number of successful store (write) operations (Stores: ok=success).
# TYPE node_fscache_stores_success_total counter
node_fscache_stores_success_total 0
# HELP node_fscache_stores_total Number of store (write) operations attempted (Stores: n=attempts).
# TYPE node_fscache_stores_total counter
node_fscache_stores_total 0
# HELP node_fscache_updates_total Number of update operations (Updates: n=tot).
# TYPE node_fscache_updates_total counter
node_fscache_updates_total 0
`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", tc.procPath}); err != nil {
				t.Fatal(err)
			}

			// Create collector
			collector, err := NewFscacheCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
			if err != nil {
				t.Fatal(err)
			}

			// Register collector
			registry := prometheus.NewRegistry()
			registry.MustRegister(collector)

			// Compare metrics
			err = testutil.GatherAndCompare(registry, strings.NewReader(tc.expected))
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
