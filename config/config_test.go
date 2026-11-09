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

package config

import (
	"strings"
	"testing"
)

func TestValidateRejectsInvalidCPUInfoFilters(t *testing.T) {
	for _, test := range []struct {
		name      string
		cpu       CPUConfig
		wantField string
	}{
		{name: "flags include", cpu: CPUConfig{FlagsInclude: "("}, wantField: "CPU.FlagsInclude"},
		{name: "bugs include", cpu: CPUConfig{BugsInclude: "("}, wantField: "CPU.BugsInclude"},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := NewConfigWithDefaults()
			cfg.CPU = test.cpu
			if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), test.wantField) {
				t.Fatalf("want error naming %s, got %v", test.wantField, err)
			}
		})
	}
}
