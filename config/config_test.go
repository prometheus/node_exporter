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

func TestValidateRejectsInvalidRegexps(t *testing.T) {
	for _, test := range []struct {
		name      string
		modify    func(*Config)
		wantField string
	}{
		{name: "cpu flags include", modify: func(c *Config) { c.CPU.FlagsInclude = "(" }, wantField: "CPU.FlagsInclude"},
		{name: "cpu bugs include", modify: func(c *Config) { c.CPU.BugsInclude = "(" }, wantField: "CPU.BugsInclude"},
		{name: "vmstat fields", modify: func(c *Config) { c.VMStat.Fields = "(" }, wantField: "VMStat.Fields"},
		{name: "netstat fields", modify: func(c *Config) { c.Netstat.Fields = "(" }, wantField: "Netstat.Fields"},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := NewConfigWithDefaults()
			test.modify(&cfg)
			if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), test.wantField) {
				t.Fatalf("want error naming %s, got %v", test.wantField, err)
			}
		})
	}
}
