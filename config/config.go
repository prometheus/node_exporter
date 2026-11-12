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
	"fmt"
	"regexp"

	"github.com/prometheus/procfs"
)

// Config contains every option of a node_exporter runtime. Each field has a
// command-line flag of the same meaning, named in its comment.
type Config struct {
	// CollectorDisableDefaults turns off every collector that Collectors
	// does not name, like --collector.disable-defaults.
	CollectorDisableDefaults bool
	// Collectors enables or disables collectors by name, like
	// --collector.<name> and --no-collector.<name>. It takes precedence
	// over CollectorDisableDefaults and the collector defaults.
	Collectors map[string]bool

	Path PathConfig

	CPU         CPUConfig
	Textfile    TextfileConfig
	NTP         NTPConfig
	Stat        StatConfig
	VMStat      VMStatConfig
	Netstat     NetstatConfig
	Bcache      BcacheConfig
	NVMe        NVMeConfig
	RAPL        RAPLConfig
	ARP         ARPConfig
	Ethtool     EthtoolConfig
	Infiniband  InfinibandConfig
	NetClass    NetClassConfig
	NetDev      NetDevConfig
	Qdisc       QdiscConfig
	IPVS        IPVSConfig
	Wifi        WifiConfig
	Diskstats   DiskstatsConfig
	Filesystem  FilesystemConfig
	Systemd     SystemdConfig
	Perf        PerfConfig
	Sysctl      SysctlConfig
	Hwmon       HwmonConfig
	CPUFreq     CPUFreqConfig
	Interrupts  InterruptsConfig
	PCIDevice   PCIDeviceConfig
	PowerSupply PowerSupplyConfig
	Resctrl     ResctrlConfig
	Runit       RunitConfig
	Supervisord SupervisordConfig
	Slabinfo    SlabinfoConfig
	Tapestats   TapestatsConfig
}

const (
	DefaultPathProcFS   = procfs.DefaultMountPoint
	DefaultPathSysFS    = "/sys"
	DefaultPathRootFS   = "/"
	DefaultPathUdevData = "/run/udev/data"
)

// PathConfig contains the mount points the collectors read from.
type PathConfig struct {
	// ProcFS is the procfs mount point, like --path.procfs.
	ProcFS string
	// SysFS is the sysfs mount point, like --path.sysfs.
	SysFS string
	// RootFS is the rootfs mount point, like --path.rootfs.
	RootFS string
	// UdevData is the udev data path, like --path.udev.data.
	UdevData string
}

// NewConfigWithDefaults returns a Config with the node_exporter defaults.
func NewConfigWithDefaults() Config {
	return Config{
		Path: PathConfig{
			ProcFS:   DefaultPathProcFS,
			SysFS:    DefaultPathSysFS,
			RootFS:   DefaultPathRootFS,
			UdevData: DefaultPathUdevData,
		},
		CPU: CPUConfig{
			EnableGuest: DefaultCPUEnableGuest,
		},
		NTP: NTPConfig{
			Server:               DefaultNTPServer,
			ServerPort:           DefaultNTPServerPort,
			ProtocolVersion:      DefaultNTPProtocolVersion,
			IPTTL:                DefaultNTPIPTTL,
			MaxDistance:          DefaultNTPMaxDistance,
			LocalOffsetTolerance: DefaultNTPLocalOffsetTolerance,
		},
		VMStat: VMStatConfig{
			Fields: DefaultVMStatFields,
		},
		Netstat: NetstatConfig{
			Fields: DefaultNetstatFields,
		},
		NVMe: NVMeConfig{
			NamespaceUsage: DefaultNVMeNamespaceUsage,
		},
		ARP: ARPConfig{
			Netlink: DefaultARPNetlink,
		},
		Ethtool: EthtoolConfig{
			MetricsInclude: DefaultEthtoolMetricsInclude,
		},
		NetClass: NetClassConfig{
			IgnoredDevices: DefaultNetClassIgnoredDevices,
		},
		NetDev: NetDevConfig{
			Netlink: DefaultNetDevNetlink,
		},
		IPVS: IPVSConfig{
			BackendLabels: DefaultIPVSBackendLabels,
		},
		Diskstats: DiskstatsConfig{
			DeviceExclude: DefaultDiskstatsDeviceExclude,
		},
		Filesystem: FilesystemConfig{
			MountPointsExclude: DefaultFilesystemMountPointsExclude,
			FSTypesExclude:     DefaultFilesystemFSTypesExclude,
			MountTimeout:       DefaultFilesystemMountTimeout,
			StatWorkerCount:    DefaultFilesystemStatWorkerCount,
		},
		Systemd: SystemdConfig{
			UnitInclude: DefaultSystemdUnitInclude,
			UnitExclude: DefaultSystemdUnitExclude,
		},
		Interrupts: InterruptsConfig{
			IncludeZeros: DefaultInterruptsIncludeZeros,
		},
		PowerSupply: PowerSupplyConfig{
			IgnoredSupplies: DefaultPowerSupplyIgnoredSupplies,
		},
		Resctrl: ResctrlConfig{
			Path: DefaultResctrlPath,
		},
		Runit: RunitConfig{
			ServiceDir: DefaultRunitServiceDir,
		},
		Supervisord: SupervisordConfig{
			URL: DefaultSupervisordURL,
		},
		Slabinfo: SlabinfoConfig{
			SlabsInclude: DefaultSlabinfoSlabsInclude,
		},
		Tapestats: TapestatsConfig{
			IgnoredDevices: DefaultTapestatsIgnoredDevices,
		},
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	for _, filter := range []func() (string, string, error){
		c.Diskstats.DeviceFilter,
		c.NetDev.DeviceFilter,
		c.Filesystem.MountPointsFilter,
		c.Filesystem.FSTypesFilter,
	} {
		if _, _, err := filter(); err != nil {
			return err
		}
	}
	for _, r := range []struct{ name, expr string }{
		{"CPU.FlagsInclude", c.CPU.FlagsInclude},
		{"CPU.BugsInclude", c.CPU.BugsInclude},
		{"VMStat.Fields", c.VMStat.Fields},
		{"Netstat.Fields", c.Netstat.Fields},
		{"ARP.DeviceInclude", c.ARP.DeviceInclude},
		{"ARP.DeviceExclude", c.ARP.DeviceExclude},
		{"Ethtool.DeviceInclude", c.Ethtool.DeviceInclude},
		{"Ethtool.DeviceExclude", c.Ethtool.DeviceExclude},
		{"Ethtool.MetricsInclude", c.Ethtool.MetricsInclude},
		{"Infiniband.DeviceInclude", c.Infiniband.DeviceInclude},
		{"Infiniband.DeviceExclude", c.Infiniband.DeviceExclude},
		{"NetClass.IgnoredDevices", c.NetClass.IgnoredDevices},
		{"NetDev.DeviceInclude", c.NetDev.DeviceInclude},
		{"NetDev.DeviceExclude", c.NetDev.DeviceExclude},
		{"Qdisc.DeviceInclude", c.Qdisc.DeviceInclude},
		{"Qdisc.DeviceExclude", c.Qdisc.DeviceExclude},
		{"Diskstats.DeviceExclude", c.Diskstats.DeviceExclude},
		{"Diskstats.DeviceInclude", c.Diskstats.DeviceInclude},
		{"Filesystem.MountPointsExclude", c.Filesystem.MountPointsExclude},
		{"Filesystem.MountPointsInclude", c.Filesystem.MountPointsInclude},
		{"Filesystem.FSTypesExclude", c.Filesystem.FSTypesExclude},
		{"Filesystem.FSTypesInclude", c.Filesystem.FSTypesInclude},
		{"Systemd.UnitInclude", c.Systemd.UnitInclude},
		{"Systemd.UnitExclude", c.Systemd.UnitExclude},
		{"Hwmon.ChipInclude", c.Hwmon.ChipInclude},
		{"Hwmon.ChipExclude", c.Hwmon.ChipExclude},
		{"Hwmon.SensorInclude", c.Hwmon.SensorInclude},
		{"Hwmon.SensorExclude", c.Hwmon.SensorExclude},
		{"Interrupts.NameInclude", c.Interrupts.NameInclude},
		{"Interrupts.NameExclude", c.Interrupts.NameExclude},
		{"PowerSupply.IgnoredSupplies", c.PowerSupply.IgnoredSupplies},
		{"Slabinfo.SlabsInclude", c.Slabinfo.SlabsInclude},
		{"Slabinfo.SlabsExclude", c.Slabinfo.SlabsExclude},
		{"Tapestats.IgnoredDevices", c.Tapestats.IgnoredDevices},
	} {
		if _, err := regexp.Compile(r.expr); err != nil {
			return fmt.Errorf("invalid %s: %w", r.name, err)
		}
	}
	return nil
}
