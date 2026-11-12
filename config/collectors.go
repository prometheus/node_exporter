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
	"runtime"
	"time"
)

const (
	DefaultCPUEnableGuest = true
)

// CPUConfig contains the options for the cpu collector.
type CPUConfig struct {
	// EnableGuest exposes node_cpu_guest_seconds_total, like
	// --collector.cpu.guest.
	EnableGuest bool
	// EnableInfo exposes node_cpu_info, like --collector.cpu.info.
	EnableInfo bool
	// FlagsInclude and BugsInclude are regular expressions that filter the
	// flags and bugs fields of node_cpu_info, like
	// --collector.cpu.info.flags-include and --collector.cpu.info.bugs-include.
	// Setting either one also enables node_cpu_info.
	FlagsInclude string
	BugsInclude  string
}

// TextfileConfig contains the options for the textfile collector.
type TextfileConfig struct {
	// Directories to read text files with metrics from, like
	// --collector.textfile.directory. Glob patterns are supported.
	Directories []string
}

// DefaultNTPMaxDistance is 3.46608s ~ 1.5s + PHI * (1 << maxPoll), where 1.5s is MAXDIST
// from ntp.org, it is 1.0 in RFC5905. The max-distance option is used as-is without phi*(1<<poll).
const (
	DefaultNTPServer               = "127.0.0.1"
	DefaultNTPServerPort           = 123
	DefaultNTPProtocolVersion      = 4
	DefaultNTPIPTTL                = 1
	DefaultNTPMaxDistance          = 3466080 * time.Microsecond
	DefaultNTPLocalOffsetTolerance = time.Millisecond
)

// NTPConfig contains the options for the ntp collector.
type NTPConfig struct {
	// Server is --collector.ntp.server: NTP server to use for ntp collector.
	Server string
	// ServerPort is --collector.ntp.server-port: UDP port number to connect to on NTP server.
	ServerPort int
	// ProtocolVersion is --collector.ntp.protocol-version: NTP protocol version.
	ProtocolVersion int
	// ServerIsLocal is --collector.ntp.server-is-local: Certify that collector.ntp.server address is not a public ntp server.
	ServerIsLocal bool
	// IPTTL is --collector.ntp.ip-ttl: IP TTL to use while sending NTP query.
	IPTTL int
	// MaxDistance is --collector.ntp.max-distance: Max accumulated distance to the root.
	MaxDistance time.Duration
	// LocalOffsetTolerance is --collector.ntp.local-offset-tolerance: Offset between local clock and local ntpd time to tolerate.
	LocalOffsetTolerance time.Duration
}

// StatConfig contains the options for the stat collector.
type StatConfig struct {
	// Softirq is --collector.stat.softirq: Export softirq calls per vector.
	Softirq bool
}

const (
	DefaultVMStatFields = "^(oom_kill|pgpg|pswp|pg.*fault).*"
)

// VMStatConfig contains the options for the vmstat collector.
type VMStatConfig struct {
	// Fields is --collector.vmstat.fields: Regexp of fields to return for vmstat collector.
	Fields string
}

const (
	DefaultNetstatFields = "^(.*_(InErrors|InErrs)|Ip_Forwarding|Ip(6|Ext)_(InOctets|OutOctets)|Icmp6?_(InMsgs|OutMsgs)|TcpExt_(Listen.*|Syncookies.*|TCPSynRetrans|TCPTimeouts|TCPOFOQueue|TCPRcvQDrop)|Tcp_(ActiveOpens|InSegs|OutSegs|OutRsts|PassiveOpens|RetransSegs|CurrEstab)|Udp6?_(InDatagrams|OutDatagrams|NoPorts|RcvbufErrors|SndbufErrors))$"
)

// NetstatConfig contains the options for the netstat collector.
type NetstatConfig struct {
	// Fields is --collector.netstat.fields: Regexp of fields to return for netstat collector.
	Fields string
}

// BcacheConfig contains the options for the bcache collector.
type BcacheConfig struct {
	// PriorityStats is --collector.bcache.priorityStats: Expose expensive priority stats.
	PriorityStats bool
}

const (
	DefaultNVMeNamespaceUsage = true
)

// NVMeConfig contains the options for the nvme collector.
type NVMeConfig struct {
	// NamespaceUsage is --collector.nvme.namespace-usage: Expose namespace used bytes (reads nuse, which makes the kernel query the controller).
	NamespaceUsage bool
}

// RAPLConfig contains the options for the rapl collector.
type RAPLConfig struct {
	// EnableZoneLabel is --collector.rapl.enable-zone-label: adds a zone label to the RAPL metrics.
	EnableZoneLabel bool
}

const (
	DefaultARPNetlink = true
)

// ARPConfig contains the options for the arp collector.
type ARPConfig struct {
	// DeviceInclude is --collector.arp.device-include: Regexp of arp devices to include (mutually exclusive to device-exclude).
	DeviceInclude string
	// DeviceExclude is --collector.arp.device-exclude: Regexp of arp devices to exclude (mutually exclusive to device-include).
	DeviceExclude string
	// Netlink is --collector.arp.netlink: Use netlink to gather stats instead of /proc/net/arp.
	Netlink bool
}

const (
	DefaultEthtoolMetricsInclude = ".*"
)

// EthtoolConfig contains the options for the ethtool collector.
type EthtoolConfig struct {
	// DeviceInclude is --collector.ethtool.device-include: Regexp of ethtool devices to include (mutually exclusive to device-exclude).
	DeviceInclude string
	// DeviceExclude is --collector.ethtool.device-exclude: Regexp of ethtool devices to exclude (mutually exclusive to device-include).
	DeviceExclude string
	// MetricsInclude is --collector.ethtool.metrics-include: Regexp of ethtool stats to include.
	MetricsInclude string
}

// InfinibandConfig contains the options for the infiniband collector.
type InfinibandConfig struct {
	// DeviceInclude is --collector.infiniband.device-include: Regexp of infiniband devices to include (mutually exclusive to device-exclude).
	DeviceInclude string
	// DeviceExclude is --collector.infiniband.device-exclude: Regexp of infiniband devices to exclude (mutually exclusive to device-include).
	DeviceExclude string
}

const (
	DefaultNetClassIgnoredDevices = "^$"
)

// NetClassConfig contains the options for the netclass collector.
type NetClassConfig struct {
	// IgnoredDevices is --collector.netclass.ignored-devices: Regexp of net devices to ignore for netclass collector.
	IgnoredDevices string
	// IgnoreInvalidSpeed is --collector.netclass.ignore-invalid-speed: Ignore devices where the speed is invalid.
	IgnoreInvalidSpeed bool
	// Netlink is --collector.netclass.netlink: Use netlink to gather stats instead of /proc/net/dev.
	Netlink bool
	// RTNLWithStats is --collector.netclass_rtnl.with-stats: Expose the statistics for each network device, replacing netdev collector.
	RTNLWithStats bool
}

const (
	DefaultNetDevNetlink = true
)

// NetDevConfig contains the options for the netdev collector.
type NetDevConfig struct {
	// DeviceInclude is --collector.netdev.device-include: Regexp of net devices to include (mutually exclusive to device-exclude).
	DeviceInclude string
	// DeviceExclude is --collector.netdev.device-exclude: Regexp of net devices to exclude (mutually exclusive to device-include).
	DeviceExclude string
	// AddressInfo is --collector.netdev.address-info: Collect address-info for every device.
	AddressInfo bool
	// DetailedMetrics is --collector.netdev.enable-detailed-metrics: Use (incompatible) metric names that provide more detailed stats on Linux.
	DetailedMetrics bool
	// Netlink is --collector.netdev.netlink: Use netlink to gather stats instead of /proc/net/dev.
	Netlink bool
	// LabelIfAlias is --collector.netdev.label-ifalias: Add ifAlias label.
	LabelIfAlias bool
}

// QdiscConfig contains the options for the qdisc collector.
type QdiscConfig struct {
	// Fixtures is --collector.qdisc.fixtures: test fixtures to use for qdisc collector end-to-end testing.
	Fixtures string
	// DeviceInclude is --collector.qdisc.device-include: Regexp of qdisc devices to include (mutually exclusive to device-exclude).
	DeviceInclude string
	// DeviceExclude is --collector.qdisc.device-exclude: Regexp of qdisc devices to exclude (mutually exclusive to device-include).
	DeviceExclude string
}

const (
	DefaultIPVSBackendLabels = "local_address,local_port,remote_address,remote_port,proto,local_mark"
)

// IPVSConfig contains the options for the ipvs collector.
type IPVSConfig struct {
	// BackendLabels is --collector.ipvs.backend-labels: Comma separated list for IPVS backend stats labels.
	BackendLabels string
}

// WifiConfig contains the options for the wifi collector.
type WifiConfig struct {
	// Fixtures is --collector.wifi.fixtures: test fixtures to use for wifi collector metrics.
	Fixtures string
}

// DefaultDiskstatsDeviceExclude depends on the operating system.
var DefaultDiskstatsDeviceExclude = func() string {
	if runtime.GOOS == "linux" {
		return "^(z?ram|loop|fd|(h|s|v|xv)d[a-z]|nvme\\d+n\\d+p)\\d+$"
	}
	return ""
}()

// DiskstatsConfig contains the options for the diskstats collector.
type DiskstatsConfig struct {
	// DeviceExclude is --collector.diskstats.device-exclude: Regexp of diskstats devices to exclude (mutually exclusive to device-include).
	// It is ignored when it has the default value and DeviceInclude is set.
	DeviceExclude string
	// DeviceInclude is --collector.diskstats.device-include: Regexp of diskstats devices to include (mutually exclusive to device-exclude).
	DeviceInclude string
}

var (
	// DefaultFilesystemMountPointsExclude depends on the operating system.
	DefaultFilesystemMountPointsExclude = defaultFilesystemExcludes[0]
	// DefaultFilesystemFSTypesExclude depends on the operating system.
	DefaultFilesystemFSTypesExclude = defaultFilesystemExcludes[1]
)

// defaultFilesystemExcludes holds the mount point and filesystem type
// excludes for the operating system.
var defaultFilesystemExcludes = func() [2]string {
	switch runtime.GOOS {
	case "linux":
		return [2]string{
			"^/(dev|proc|run/credentials/.+|sys|var/lib/docker/.+|var/lib/containers/storage/.+)($|/)",
			"^(autofs|binfmt_misc|bpf|cgroup2?|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|iso9660|mqueue|nsfs|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|selinuxfs|squashfs|erofs|sysfs|tracefs)$",
		}
	case "aix":
		return [2]string{"^/(dev|aha)($|/)", "^procfs$"}
	case "netbsd":
		return [2]string{"^/(dev)($|/)", "^(kernfs|procfs|ptyfs|fdesc)$"}
	case "darwin", "dragonfly", "freebsd", "openbsd":
		return [2]string{"^/(dev)($|/)", "^devfs$"}
	}
	return [2]string{}
}()

const (
	DefaultFilesystemMountTimeout    = 5 * time.Second
	DefaultFilesystemStatWorkerCount = 4
)

// FilesystemConfig contains the options for the filesystem collector.
type FilesystemConfig struct {
	// MountPointsExclude is --collector.filesystem.mount-points-exclude: Regexp of mount points to exclude for filesystem collector. (mutually exclusive to mount-points-include)
	// It is ignored when it has the default value and MountPointsInclude is set.
	MountPointsExclude string
	// MountPointsInclude is --collector.filesystem.mount-points-include: Regexp of mount points to include for filesystem collector. (mutually exclusive to mount-points-exclude)
	MountPointsInclude string
	// FSTypesExclude is --collector.filesystem.fs-types-exclude: Regexp of filesystem types to exclude for filesystem collector. (mutually exclusive to fs-types-include)
	// It is ignored when it has the default value and FSTypesInclude is set.
	FSTypesExclude string
	// FSTypesInclude is --collector.filesystem.fs-types-include: Regexp of filesystem types to include for filesystem collector. (mutually exclusive to fs-types-exclude)
	FSTypesInclude string
	// MountTimeout is --collector.filesystem.mount-timeout: how long to wait for a mount to respond before marking it as stale.
	MountTimeout time.Duration
	// StatWorkerCount is --collector.filesystem.stat-workers: how many stat calls to process simultaneously.
	StatWorkerCount int
}

const (
	DefaultSystemdUnitInclude = ".+"
	DefaultSystemdUnitExclude = ".+\\.(automount|device|mount|scope|slice)"
)

// SystemdConfig contains the options for the systemd collector.
type SystemdConfig struct {
	// UnitInclude is --collector.systemd.unit-include: Regexp of systemd units to include. Units must both match include and not match exclude to be included.
	UnitInclude string
	// UnitExclude is --collector.systemd.unit-exclude: Regexp of systemd units to exclude. Units must both match include and not match exclude to be included.
	UnitExclude string
	// Private is --collector.systemd.private: Establish a private, direct connection to systemd without dbus.
	Private bool
	// EnableTaskMetrics is --collector.systemd.enable-task-metrics: Enables service unit tasks metrics unit_tasks_current and unit_tasks_max.
	EnableTaskMetrics bool
	// EnableRestartsMetrics is --collector.systemd.enable-restarts-metrics: Enables service unit metric service_restart_total.
	EnableRestartsMetrics bool
	// EnableStartTimeMetrics is --collector.systemd.enable-start-time-metrics: Enables service unit metric unit_start_time_seconds.
	EnableStartTimeMetrics bool
}

// PerfConfig contains the options for the perf collector.
type PerfConfig struct {
	// CPUs is --collector.perf.cpus: List of CPUs from which perf metrics should be collected.
	CPUs string
	// Tracepoints is --collector.perf.tracepoint: perf tracepoint that should be collected.
	Tracepoints []string
	// DisableHardwareProfilers is --collector.perf.disable-hardware-profilers: disable perf hardware profilers.
	DisableHardwareProfilers bool
	// HardwareProfilers is --collector.perf.hardware-profilers: perf hardware profilers that should be collected.
	HardwareProfilers []string
	// DisableSoftwareProfilers is --collector.perf.disable-software-profilers: disable perf software profilers.
	DisableSoftwareProfilers bool
	// SoftwareProfilers is --collector.perf.software-profilers: perf software profilers that should be collected.
	SoftwareProfilers []string
	// DisableCacheProfilers is --collector.perf.disable-cache-profilers: disable perf cache profilers.
	DisableCacheProfilers bool
	// CacheProfilers is --collector.perf.cache-profilers: perf cache profilers that should be collected.
	CacheProfilers []string
}

// SysctlConfig contains the options for the sysctl collector.
type SysctlConfig struct {
	// Include is --collector.sysctl.include: Select sysctl metrics to include.
	Include []string
	// IncludeInfo is --collector.sysctl.include-info: Select sysctl metrics to include as info metrics.
	IncludeInfo []string
}

// HwmonConfig contains the options for the hwmon collector.
type HwmonConfig struct {
	// ChipInclude is --collector.hwmon.chip-include: Regexp of hwmon chip to include (mutually exclusive to device-exclude).
	ChipInclude string
	// ChipExclude is --collector.hwmon.chip-exclude: Regexp of hwmon chip to exclude (mutually exclusive to device-include).
	ChipExclude string
	// SensorInclude is --collector.hwmon.sensor-include: Regexp of hwmon sensor to include (mutually exclusive to sensor-exclude).
	SensorInclude string
	// SensorExclude is --collector.hwmon.sensor-exclude: Regexp of hwmon sensor to exclude (mutually exclusive to sensor-include).
	SensorExclude string
}

// CPUFreqConfig contains the options for the cpufreq collector.
type CPUFreqConfig struct {
	// EnablePrefix is --collector.cpufreq.enable-cpufreq-prefix: Expose cpufreq metrics with the node_cpufreq_ prefix instead of node_cpu_.
	EnablePrefix bool
}

const (
	DefaultInterruptsIncludeZeros = true
)

// InterruptsConfig contains the options for the interrupts collector.
type InterruptsConfig struct {
	// NameInclude is --collector.interrupts.name-include: Regexp of interrupts name to include (mutually exclusive to --collector.interrupts.name-exclude).
	NameInclude string
	// NameExclude is --collector.interrupts.name-exclude: Regexp of interrupts name to exclude (mutually exclusive to --collector.interrupts.name-include).
	NameExclude string
	// IncludeZeros is --collector.interrupts.include-zeros: Include interrupts that have a zero value.
	IncludeZeros bool
}

// PCIDeviceConfig contains the options for the pcidevice collector.
type PCIDeviceConfig struct {
	// IDsFile is --collector.pcidevice.idsfile: Path to pci.ids file to use for PCI device identification.
	IDsFile string
	// Names is --collector.pcidevice.names: Enable PCI device name resolution (requires pci.ids file).
	Names bool
}

const (
	DefaultPowerSupplyIgnoredSupplies = "^$"
)

// PowerSupplyConfig contains the options for the powersupplyclass collector.
type PowerSupplyConfig struct {
	// IgnoredSupplies is --collector.powersupply.ignored-supplies: Regexp of power supplies to ignore for powersupplyclass collector.
	IgnoredSupplies string
}

const (
	// DefaultResctrlPath matches resctrlfs.DefaultMountPoint, which only
	// builds on Linux.
	DefaultResctrlPath = "/sys/fs/resctrl"
)

// ResctrlConfig contains the options for the resctrl collector.
type ResctrlConfig struct {
	// Path is --collector.resctrl.path: resctrl filesystem mountpoint.
	Path string
}

const (
	DefaultRunitServiceDir = "/etc/service"
)

// RunitConfig contains the options for the runit collector.
type RunitConfig struct {
	// ServiceDir is --collector.runit.servicedir: Path to runit service directory.
	ServiceDir string
}

const (
	DefaultSupervisordURL = "http://localhost:9001/RPC2"
)

// SupervisordConfig contains the options for the supervisord collector.
type SupervisordConfig struct {
	// URL is --collector.supervisord.url: XML RPC endpoint.
	URL string
}

const (
	DefaultSlabinfoSlabsInclude = ".*"
)

// SlabinfoConfig contains the options for the slabinfo collector.
type SlabinfoConfig struct {
	// SlabsInclude is --collector.slabinfo.slabs-include: Regexp of slabs to include in slabinfo collector.
	SlabsInclude string
	// SlabsExclude is --collector.slabinfo.slabs-exclude: Regexp of slabs to exclude in slabinfo collector.
	SlabsExclude string
}

const (
	DefaultTapestatsIgnoredDevices = "^$"
)

// TapestatsConfig contains the options for the tapestats collector.
type TapestatsConfig struct {
	// IgnoredDevices is --collector.tapestats.ignored-devices: Regexp of devices to ignore for tapestats.
	IgnoredDevices string
}

// includeExclude returns the exclude and include patterns a collector uses.
// An exclude that still has its default value gives way to an include, like
// the command line does when only the include flag is set.
func includeExclude(name, flag, exclude, defaultExclude, include string) (string, string, error) {
	if include != "" && exclude == defaultExclude {
		exclude = ""
	}
	if exclude != "" && include != "" {
		return "", "", fmt.Errorf("%[1]sExclude (--%[2]s-exclude) and %[1]sInclude (--%[2]s-include) are mutually exclusive", name, flag)
	}
	return exclude, include, nil
}

// DeviceFilter returns the device exclude and include patterns.
func (c DiskstatsConfig) DeviceFilter() (exclude, include string, err error) {
	return includeExclude("Diskstats.Device", "collector.diskstats.device", c.DeviceExclude, DefaultDiskstatsDeviceExclude, c.DeviceInclude)
}

// DeviceFilter returns the device exclude and include patterns.
func (c NetDevConfig) DeviceFilter() (exclude, include string, err error) {
	return includeExclude("NetDev.Device", "collector.netdev.device", c.DeviceExclude, "", c.DeviceInclude)
}

// MountPointsFilter returns the mount point exclude and include patterns.
func (c FilesystemConfig) MountPointsFilter() (exclude, include string, err error) {
	return includeExclude("Filesystem.MountPoints", "collector.filesystem.mount-points", c.MountPointsExclude, DefaultFilesystemMountPointsExclude, c.MountPointsInclude)
}

// FSTypesFilter returns the filesystem type exclude and include patterns.
func (c FilesystemConfig) FSTypesFilter() (exclude, include string, err error) {
	return includeExclude("Filesystem.FSTypes", "collector.filesystem.fs-types", c.FSTypesExclude, DefaultFilesystemFSTypesExclude, c.FSTypesInclude)
}
