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
	"testing"
)

func TestInvalidRaplName(t *testing.T) {
	testCases := []struct {
		desc string
		name string
	}{
		{
			desc: "invalid metrics name with numeric",
			name: "package9",
		},
		{
			desc: "invalid metrics name with dash",
			name: "package-",
		},
		{
			desc: "invalid metrics name",
			name: "package-0-die-0",
		},
	}

	// these are irrelevent in this test so accepting default value
	var raplZonePath, index string
	var newMicrojoules uint64

	for i := range testCases {
		t.Run(testCases[i].desc, func(t *testing.T) {
			// getRzMetric panics if any error occurs
			getRzMetric(testCases[i].name, raplZonePath, index, newMicrojoules)
		})
	}
}
