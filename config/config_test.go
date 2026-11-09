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
		{name: "arp deviceinclude", modify: func(c *Config) { c.ARP.DeviceInclude = "(" }, wantField: "ARP.DeviceInclude"},
		{name: "arp deviceexclude", modify: func(c *Config) { c.ARP.DeviceExclude = "(" }, wantField: "ARP.DeviceExclude"},
		{name: "ethtool deviceinclude", modify: func(c *Config) { c.Ethtool.DeviceInclude = "(" }, wantField: "Ethtool.DeviceInclude"},
		{name: "ethtool deviceexclude", modify: func(c *Config) { c.Ethtool.DeviceExclude = "(" }, wantField: "Ethtool.DeviceExclude"},
		{name: "ethtool metricsinclude", modify: func(c *Config) { c.Ethtool.MetricsInclude = "(" }, wantField: "Ethtool.MetricsInclude"},
		{name: "infiniband deviceinclude", modify: func(c *Config) { c.Infiniband.DeviceInclude = "(" }, wantField: "Infiniband.DeviceInclude"},
		{name: "infiniband deviceexclude", modify: func(c *Config) { c.Infiniband.DeviceExclude = "(" }, wantField: "Infiniband.DeviceExclude"},
		{name: "netclass ignoreddevices", modify: func(c *Config) { c.NetClass.IgnoredDevices = "(" }, wantField: "NetClass.IgnoredDevices"},
		{name: "netdev deviceinclude", modify: func(c *Config) { c.NetDev.DeviceInclude = "(" }, wantField: "NetDev.DeviceInclude"},
		{name: "netdev deviceexclude", modify: func(c *Config) { c.NetDev.DeviceExclude = "(" }, wantField: "NetDev.DeviceExclude"},
		{name: "qdisc deviceinclude", modify: func(c *Config) { c.Qdisc.DeviceInclude = "(" }, wantField: "Qdisc.DeviceInclude"},
		{name: "qdisc deviceexclude", modify: func(c *Config) { c.Qdisc.DeviceExclude = "(" }, wantField: "Qdisc.DeviceExclude"},
		{name: "diskstats deviceexclude", modify: func(c *Config) { c.Diskstats.DeviceExclude = "(" }, wantField: "Diskstats.DeviceExclude"},
		{name: "diskstats deviceinclude", modify: func(c *Config) { c.Diskstats.DeviceInclude = "(" }, wantField: "Diskstats.DeviceInclude"},
		{name: "filesystem mountpointsexclude", modify: func(c *Config) { c.Filesystem.MountPointsExclude = "(" }, wantField: "Filesystem.MountPointsExclude"},
		{name: "filesystem mountpointsinclude", modify: func(c *Config) { c.Filesystem.MountPointsInclude = "(" }, wantField: "Filesystem.MountPointsInclude"},
		{name: "filesystem fstypesexclude", modify: func(c *Config) { c.Filesystem.FSTypesExclude = "(" }, wantField: "Filesystem.FSTypesExclude"},
		{name: "filesystem fstypesinclude", modify: func(c *Config) { c.Filesystem.FSTypesInclude = "(" }, wantField: "Filesystem.FSTypesInclude"},
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
