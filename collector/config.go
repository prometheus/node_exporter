// Copyright 2023 The Prometheus Authors
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

type NodeCollectorConfig struct {
	Arp                   ArpConfig
	Bcache                BcacheConfig
	CPU                   CPUConfig
	DiskstatsDeviceFilter DiskstatsDeviceFilterConfig
	Ethtool               EthtoolConfig
	Filesystem            FilesystemConfig
	HwMon                 HwMonConfig
	IPVS                  IPVSConfig
	NetClass              NetClassConfig
	NetDev                NetDevConfig
	NetStat               NetStatConfig
	NTP                   NTPConfig
	Path                  PathConfig
	Perf                  PerfConfig
	PowerSupplyClass      PowerSupplyClassConfig
	Qdisc                 QdiscConfig
	Rapl                  RaplConfig
	Runit                 RunitConfig
	Stat                  StatConfig
	Supervisord           SupervisordConfig
	Sysctl                SysctlConfig
	Systemd               SystemdConfig
	Tapestats             TapestatsConfig
	TextFile              TextFileConfig
	VmStat                VmStatConfig
	Wifi                  WifiConfig
}
