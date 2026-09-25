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

package numatopology

import "testing"

func TestCountCPUList(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"0", 1},
		{"0-3", 4},
		{"0-3,8-11", 8},
		{"0,2,4", 3},
		{"0-1,4-7", 6},
	}
	for _, tt := range tests {
		got := CountCPUList(tt.input)
		if got != tt.want {
			t.Errorf("CountCPUList(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
