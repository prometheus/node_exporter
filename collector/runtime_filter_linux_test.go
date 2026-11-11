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

//go:build !nodiskstats && !nofilesystem

package collector

import (
	"strings"
	"testing"

	"github.com/prometheus/node_exporter/config"
)

func TestNewRuntimeIncludeAndExcludeFilters(t *testing.T) {
	for _, test := range []struct {
		name      string
		collector string
		modify    func(*config.Config)
		wantErr   string
	}{
		{
			name:      "diskstats include with the default exclude",
			collector: "diskstats",
			modify:    func(c *config.Config) { c.Diskstats.DeviceInclude = "^sd" },
		},
		{
			name:      "diskstats include with another exclude",
			collector: "diskstats",
			modify:    func(c *config.Config) { c.Diskstats.DeviceInclude, c.Diskstats.DeviceExclude = "^sd", "^sdb" },
			wantErr:   "Diskstats.DeviceExclude (--collector.diskstats.device-exclude) and Diskstats.DeviceInclude (--collector.diskstats.device-include) are mutually exclusive",
		},
		{
			name:      "filesystem includes with the default excludes",
			collector: "filesystem",
			modify: func(c *config.Config) {
				c.Filesystem.MountPointsInclude, c.Filesystem.FSTypesInclude = "^/data", "^ext4$"
			},
		},
		{
			name:      "filesystem mount points include with another exclude",
			collector: "filesystem",
			modify: func(c *config.Config) {
				c.Filesystem.MountPointsInclude, c.Filesystem.MountPointsExclude = "^/data", "^/data/tmp"
			},
			wantErr: "Filesystem.MountPointsExclude (--collector.filesystem.mount-points-exclude) and Filesystem.MountPointsInclude (--collector.filesystem.mount-points-include) are mutually exclusive",
		},
		{
			name:      "filesystem fs types include with another exclude",
			collector: "filesystem",
			modify:    func(c *config.Config) { c.Filesystem.FSTypesInclude, c.Filesystem.FSTypesExclude = "^ext4$", "^tmpfs$" },
			wantErr:   "Filesystem.FSTypesExclude (--collector.filesystem.fs-types-exclude) and Filesystem.FSTypesInclude (--collector.filesystem.fs-types-include) are mutually exclusive",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := config.NewConfigWithDefaults()
			cfg.CollectorDisableDefaults = true
			cfg.Collectors = map[string]bool{test.collector: true}
			test.modify(&cfg)

			_, err := NewRuntime(cfg, testLogger())
			if test.wantErr == "" {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("want error containing %q, got %v", test.wantErr, err)
			}
		})
	}
}
