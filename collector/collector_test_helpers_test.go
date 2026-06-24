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

package collector

import (
	"strings"
	"testing"

	dto "github.com/prometheus/client_model/go"
)

type labelMap map[string]string

func assertGaugeValue(t *testing.T, metrics map[string][]*dto.Metric, metricSubstring string, labels labelMap, expected float64) {
	t.Helper()
	for desc, ms := range metrics {
		if !strings.Contains(desc, metricSubstring) {
			continue
		}
		for _, m := range ms {
			if matchLabels(m.GetLabel(), labels) {
				got := m.GetGauge().GetValue()
				if got != expected {
					t.Errorf("%s%v: got %v, want %v", metricSubstring, labels, got, expected)
				}
				return
			}
		}
	}
	t.Errorf("metric %s%v not found", metricSubstring, labels)
}

func matchLabels(pairs []*dto.LabelPair, want labelMap) bool {
	if want == nil {
		return len(pairs) == 0
	}
	found := 0
	for _, lp := range pairs {
		if v, ok := want[lp.GetName()]; ok && v == lp.GetValue() {
			found++
		}
	}
	return found == len(want)
}
