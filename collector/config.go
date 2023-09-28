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

import "time"

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

type WifiConfig struct {
	Fixtures *string
}

type VmStatConfig struct {
	Fields *string
}

type TextFileConfig struct {
	Directory *string
}
type TapestatsConfig struct {
	IgnoredDevices *string
}

type SystemdConfig struct {
	UnitInclude            *string
	UnitIncludeSet         bool
	UnitExclude            *string
	UnitExcludeSet         bool
	OldUnitInclude         *string
	OldUnitExclude         *string
	Private                *bool
	EnableTaskMetrics      *bool
	EnableRestartsMetrics  *bool
	EnableStartTimeMetrics *bool
}

type SysctlConfig struct {
	Include     *[]string
	IncludeInfo *[]string
}

type SupervisordConfig struct {
	URL *string
}

type RunitConfig struct {
	ServiceDir *string
}

type StatConfig struct {
	Softirq *bool
}

type RaplConfig struct {
	ZoneLabel *bool
}

type QdiscConfig struct {
	Fixtures         *string
	DeviceInclude    *string
	OldDeviceInclude *string
	DeviceExclude    *string
	OldDeviceExclude *string
}

type PowerSupplyClassConfig struct {
	IgnoredPowerSupplies *string
}

type PerfConfig struct {
	CPUs           *string
	Tracepoint     *[]string
	NoHwProfiler   *bool
	HwProfiler     *[]string
	NoSwProfiler   *bool
	SwProfiler     *[]string
	NoCaProfiler   *bool
	CaProfilerFlag *[]string
}

type PathConfig struct {
	ProcPath     *string
	SysPath      *string
	RootfsPath   *string
	UdevDataPath *string
}

type NTPConfig struct {
	Server          *string
	ServerPort      *int
	ProtocolVersion *int
	ServerIsLocal   *bool
	IPTTL           *int
	MaxDistance     *time.Duration
	OffsetTolerance *time.Duration
}

type NetStatConfig struct {
	Fields *string
}

type NetDevConfig struct {
	DeviceInclude    *string
	OldDeviceInclude *string
	DeviceExclude    *string
	OldDeviceExclude *string
	AddressInfo      *bool
	DetailedMetrics  *bool
	Netlink          *bool
}

type NetClassConfig struct {
	IgnoredDevices *string
	InvalidSpeed   *bool
	Netlink        *bool
	RTNLWithStats  *bool
}

type ArpConfig struct {
	DeviceInclude *string
	DeviceExclude *string
	Netlink       *bool
}

type BcacheConfig struct {
	PriorityStats *bool
}

type CPUConfig struct {
	EnableCPUGuest *bool
	EnableCPUInfo  *bool
	FlagsInclude   *string
	BugsInclude    *string
}

type DiskstatsDeviceFilterConfig struct {
	DeviceExclude    *string
	DeviceExcludeSet bool
	OldDeviceExclude *string
	DeviceInclude    *string
}

type EthtoolConfig struct {
	DeviceInclude   *string
	DeviceExclude   *string
	IncludedMetrics *string
}

type HwMonConfig struct {
	ChipInclude *string
	ChipExclude *string
}

type FilesystemConfig struct {
	MountPointsExclude     *string
	MountPointsExcludeSet  bool
	OldMountPointsExcluded *string
	FSTypesExclude         *string
	FSTypesExcludeSet      bool
	OldFSTypesExcluded     *string
	MountTimeout           *time.Duration
	StatWorkerCount        *int
}

type IPVSConfig struct {
	Labels    *string
	LabelsSet bool
}
