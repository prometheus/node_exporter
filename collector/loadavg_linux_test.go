// Copyright 2015 The Prometheus Authors
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
	"syscall"
	"testing"

	"github.com/bouk/monkey"
)

func TestLoad(t *testing.T) {
	var scale float64 = 65536
	monkey.Patch(syscall.Sysinfo, func(in *syscall.Sysinfo_t) error {
		in.Loads = [3]uint64{
			uint64(0.20 * scale),
			uint64(0.40 * scale),
			uint64(0.60 * scale),
		}
		return nil
	})
	var (
		want       = []float64{0.20, 0.40, 0.60}
		loads, err = getLoad()
	)
	if err != nil {
		t.Fatal(err)
	}
	for i, load := range loads {
		if want[i] != load {
			t.Fatalf("want load %f, got %f", want[i], load)
		}
	}
}
