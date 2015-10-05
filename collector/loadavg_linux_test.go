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

import "testing"

func TestLoad(t *testing.T) {
	want := []float64{0.21, 0.37, 0.39}
	loads, err := parseLoad("0.21 0.37 0.39 1/719 19737")
	if err != nil {
		t.Fatal(err)
	}

	for i, load := range loads {
		if want[i] != load {
			t.Fatalf("want load %f, got %f", want[i], load)
		}
	}
}
