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

//go:build !nonetclass && linux

package collector

import "testing"

func TestVRFTable(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want uint32
		ok   bool
	}{
		{
			// IFLA_INFO_DATA of a vrf device using table 100, as
			// returned by RTM_GETLINK.
			name: "table 100",
			data: []byte{0x08, 0x00, 0x01, 0x00, 0x64, 0x00, 0x00, 0x00},
			want: 100,
			ok:   true,
		},
		{
			// Table IDs above 255 exercise the multi-byte decode.
			name: "table 1000",
			data: []byte{0x08, 0x00, 0x01, 0x00, 0xe8, 0x03, 0x00, 0x00},
			want: 1000,
			ok:   true,
		},
		{
			name: "no attributes",
			data: nil,
			ok:   false,
		},
		{
			// IFLA_VRF_TABLE absent, only an unrelated attribute.
			name: "missing table attribute",
			data: []byte{0x08, 0x00, 0x02, 0x00, 0x64, 0x00, 0x00, 0x00},
			ok:   false,
		},
		{
			name: "truncated attribute",
			data: []byte{0x08, 0x00, 0x01, 0x00, 0x64},
			ok:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := vrfTable(test.data)
			if ok != test.ok {
				t.Fatalf("vrfTable() ok = %v, want %v", ok, test.ok)
			}
			if ok && got != test.want {
				t.Errorf("vrfTable() = %d, want %d", got, test.want)
			}
		})
	}
}
