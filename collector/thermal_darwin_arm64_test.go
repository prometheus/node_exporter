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

//go:build !notherm && darwin && arm64 && cgo

package collector

import "testing"

func TestResolveSensorNames(t *testing.T) {
	for _, tc := range []struct {
		name    string
		sensors []thermalSensor
		want    []string
	}{
		{
			name:    "unique names are left alone",
			sensors: []thermalSensor{{name: "NAND CH0 temp", location: "1", registryID: "10"}},
			want:    []string{"NAND CH0 temp"},
		},
		{
			// Several "gas gauge battery" services are distinct sensors that
			// only share a product name.
			name: "a shared name is qualified with the location",
			sensors: []thermalSensor{
				{name: "gas gauge battery", location: "1413951555", registryID: "10"},
				{name: "gas gauge battery", location: "1413951574", registryID: "11"},
			},
			want: []string{"gas gauge battery_1413951555", "gas gauge battery_1413951574"},
		},
		{
			// The PMU sensors report the same product name and the same
			// location, so only the registry ID separates them.
			name: "a shared location falls back to the registry ID",
			sensors: []thermalSensor{
				{name: "PMU tdie2", location: "1414541922", registryID: "10"},
				{name: "PMU tdie2", location: "1414541922", registryID: "11"},
				{name: "PMU tdie2", location: "1414541922", registryID: "12"},
			},
			want: []string{"PMU tdie2_10", "PMU tdie2_11", "PMU tdie2_12"},
		},
		{
			name: "a missing location falls back to the registry ID",
			sensors: []thermalSensor{
				{name: "sensor", location: "", registryID: "10"},
				{name: "sensor", location: "", registryID: "11"},
			},
			want: []string{"sensor_10", "sensor_11"},
		},
		{
			// Only the services that cannot be told apart by location fall back
			// to the registry ID.
			name: "location and registry ID are mixed within one name",
			sensors: []thermalSensor{
				{name: "sensor", location: "1", registryID: "10"},
				{name: "sensor", location: "2", registryID: "11"},
				{name: "sensor", location: "2", registryID: "12"},
			},
			want: []string{"sensor_1", "sensor_11", "sensor_12"},
		},
		{
			name:    "no sensors",
			sensors: nil,
			want:    []string{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := resolveSensorNames(tc.sensors)
			if len(got) != len(tc.want) {
				t.Fatalf("got %d labels, want %d: %q", len(got), len(tc.want), got)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Errorf("label %d: got %q, want %q", i, got[i], tc.want[i])
				}
			}

			seen := make(map[string]struct{}, len(got))
			for _, l := range got {
				if _, dup := seen[l]; dup {
					t.Errorf("duplicate label %q", l)
				}
				seen[l] = struct{}{}
			}
		})
	}
}

// Every reading must be reported. Sensors sharing a product name are
// disambiguated rather than dropped.
func TestResolveSensorNamesKeepsEveryReading(t *testing.T) {
	sensors := []thermalSensor{
		{name: "a", location: "1", registryID: "10"},
		{name: "a", location: "1", registryID: "11"},
		{name: "a", location: "2", registryID: "12"},
		{name: "b", location: "1", registryID: "13"},
	}

	got := resolveSensorNames(sensors)
	if len(got) != len(sensors) {
		t.Fatalf("got %d labels for %d sensors", len(got), len(sensors))
	}

	seen := make(map[string]struct{}, len(got))
	for _, l := range got {
		if _, dup := seen[l]; dup {
			t.Errorf("duplicate label %q", l)
		}
		seen[l] = struct{}{}
	}
}
