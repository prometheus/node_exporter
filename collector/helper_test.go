// Copyright 2020 The Prometheus Authors
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

func TestBytesToString(t *testing.T) {
	tests := []struct {
		name     string
		b        []byte
		expected string
	}{
		{
			"Single null byte",
			[]byte{0},
			"",
		},
		{
			"Empty byte array",
			[]byte{},
			"",
		},
		{
			"Not null terminated",
			[]byte{65, 66, 67},
			"ABC",
		},
		{
			"Null randomly in array",
			[]byte{65, 66, 67, 0, 65, 0, 65},
			"ABC",
		},
		{
			"Array starts with null and contains other valid bytes",
			[]byte{0, 65, 66, 67, 0},
			"",
		},
	}

	for _, tt := range tests {
		name := tt.name
		b := tt.b
		result := bytesToString(b)
		expected := tt.expected

		if result != expected {
			t.Errorf("bytesToString(%#v): Name: %s, expected %#v, got %#v)", b, name, expected, result)
		}
	}
}

func TestSanitizeMetricName(t *testing.T) {
	testcases := map[string]string{
		"":                             "",
		"rx_errors":                    "rx_errors",
		"Queue[0] AllocFails":          "Queue_0_AllocFails",
		"Tx LPI entry count":           "Tx_LPI_entry_count",
		"port.VF_admin_queue_requests": "port_VF_admin_queue_requests",
		"[3]: tx_bytes":                "_3_tx_bytes",
		"     err":                     "_err",
	}

	for metricName, expected := range testcases {
		got := SanitizeMetricName(metricName)
		if expected != got {
			t.Errorf("Expected '%s' but got '%s'", expected, got)
		}
	}
}
