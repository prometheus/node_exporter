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

func TestParseMeminfo(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{
			name:  "standard sysfs meminfo format",
			input: "Node 0 MemTotal:       16384000 kB\nNode 0 MemFree:         4096 kB\n",
			want:  16384000 * 1024,
		},
		{
			name:  "inline MemTotal without node prefix",
			input: "MemTotal: 8192000 kB\n",
			want:  8192000 * 1024,
		},
		{
			name:    "missing MemTotal",
			input:   "Node 0 MemFree: 4096 kB\n",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMeminfo(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseMeminfo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseMeminfo() = %d, want %d", got, tt.want)
			}
		})
	}
}
