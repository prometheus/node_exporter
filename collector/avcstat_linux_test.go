// Copyright 2022 The Prometheus Authors
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
	"testing"
)

func TestAVCStat(t *testing.T) {
	avcStats, err := getAVCStats("fixtures/sys/fs/selinux/avc/cache_stats")
	if err != nil {
		t.Fatal(err)
	}

	if want, got := uint64(91590784), avcStats["lookups"]; want != got {
		t.Errorf("want avcstat lookups %v, got %v", want, got)
	}

	if want, got := uint64(91569452), avcStats["hits"]; want != got {
		t.Errorf("want avcstat hits %v, got %v", want, got)
	}

	if want, got := uint64(21332), avcStats["misses"]; want != got {
		t.Errorf("want avcstat misses %v, got %v", want, got)
	}

	if want, got := uint64(21332), avcStats["allocations"]; want != got {
		t.Errorf("want avcstat allocations %v, got %v", want, got)
	}

	if want, got := uint64(20400), avcStats["reclaims"]; want != got {
		t.Errorf("want avcstat reclaims %v, got %v", want, got)
	}

	if want, got := uint64(20826), avcStats["frees"]; want != got {
		t.Errorf("want avcstat frees %v, got %v", want, got)
	}
}

func TestAVCHashStat(t *testing.T) {
	avcHashStats, err := getAVCHashStats("fixtures/sys/fs/selinux/avc/hash_stats")
	if err != nil {
		t.Fatal(err)
	}

	if want, got := uint64(503), avcHashStats["entries"]; want != got {
		t.Errorf("want avc hash stat entries %v, got %v", want, got)
	}

	if want, got := uint64(512), avcHashStats["buckets_available"]; want != got {
		t.Errorf("want avc hash stat buckets available %v, got %v", want, got)
	}

	if want, got := uint64(257), avcHashStats["buckets_used"]; want != got {
		t.Errorf("want avc hash stat buckets used %v, got %v", want, got)
	}

	if want, got := uint64(8), avcHashStats["longest_chain"]; want != got {
		t.Errorf("want avc hash stat longest chain %v, got %v", want, got)
	}
}
