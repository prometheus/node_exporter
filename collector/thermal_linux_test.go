// Copyright 2019 The Prometheus Authors
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

	"github.com/prometheus/procfs/sysfs"
)

func TestThermalStats(t *testing.T) {
	fs, err := sysfs.NewFS("fixtures/sys")
	if err != nil {
		t.Fatalf("Error in getting fixture data: %v", err)
	}
	stats, err := fs.NewClassThermalZoneStats()
	if err != nil {
		t.Fatalf("Error in getting fixture data: %v", err)
	}
	if len(stats) != 1 {
		t.Errorf("wrong number of thermal stat: want 1, got %d", len(stats))
	}

	if want, got := uint64(56000), stats[0].Temp; want != got {
		t.Errorf("want thermal temp %d, got %d", want, got)
	}

}
