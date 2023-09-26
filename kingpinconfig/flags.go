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

package kingpinconfig

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/node_exporter/collector"
)

func AddFlags(a *kingpin.Application) collector.NodeCollectorConfig {
	config := collector.NodeCollectorConfig{}

	config.Arp.DeviceInclude = a.Flag("collector.arp.device-include", "Regexp of arp devices to include (mutually exclusive to device-exclude).").String()
	config.Arp.DeviceExclude = a.Flag("collector.arp.device-exclude", "Regexp of arp devices to exclude (mutually exclusive to device-include).").String()

	config.Bcache.PriorityStats = a.Flag("collector.bcache.priorityStats", "Expose expensive priority stats.").Bool()

	config.CPU.EnableCPUGuest = a.Flag("collector.cpu.guest", "Enables metric node_cpu_guest_seconds_total").Default("true").Bool()
	config.CPU.EnableCPUInfo = a.Flag("collector.cpu.info", "Enables metric cpu_info").Bool()
	config.CPU.FlagsInclude = a.Flag("collector.cpu.info.flags-include", "Filter the `flags` field in cpuInfo with a value that must be a regular expression").String()
	config.CPU.BugsInclude = a.Flag("collector.cpu.info.bugs-include", "Filter the `bugs` field in cpuInfo with a value that must be a regular expression").String()

	config.DiskstatsDeviceFilter.DiskstatsDeviceExclude = a.Flag(
		"collector.diskstats.device-exclude",
		"Regexp of diskstats devices to exclude (mutually exclusive to device-include).",
	).PreAction(func(c *kingpin.ParseContext) error {
		diskstatsDeviceExcludeSet = true
		return nil
	}).String()
	config.DiskstatsDeviceFilter.OldDiskstatsDeviceExclude = a.Flag(
		"collector.diskstats.ignored-devices",
		"DEPRECATED: Use collector.diskstats.device-exclude",
	).Hidden().String()
	config.DiskstatsDeviceFilter.DiskstatsDeviceInclude = a.Flag("collector.diskstats.device-include", "Regexp of diskstats devices to include (mutually exclusive to device-exclude).").String()

	config.Ethtool.DeviceInclude = a.Flag("collector.ethtool.device-include", "Regexp of ethtool devices to include (mutually exclusive to device-exclude).").String()
	config.Ethtool.DeviceExclude = a.Flag("collector.ethtool.device-exclude", "Regexp of ethtool devices to exclude (mutually exclusive to device-include).").String()
	config.Ethtool.IncludedMetrics = a.Flag("collector.ethtool.metrics-include", "Regexp of ethtool stats to include.").Default(".*").String()

	config.Filesystem.MountPointsExclude = a.Flag(
		"collector.filesystem.mount-points-exclude",
		"Regexp of mount points to exclude for filesystem collector.",
	).PreAction(func(c *kingpin.ParseContext) error {
		mountPointsExcludeSet = true
		return nil
	}).String()
	config.Filesystem.OldMountPointsExcluded = a.Flag(
		"collector.filesystem.ignored-mount-points",
		"Regexp of mount points to ignore for filesystem collector.",
	).Hidden().String()

	config.Filesystem.FSTypesExclude = a.Flag(
		"collector.filesystem.fs-types-exclude",
		"Regexp of filesystem types to exclude for filesystem collector.",
	).PreAction(func(c *kingpin.ParseContext) error {
		fsTypesExcludeSet = true
		return nil
	}).String()
	config.Filesystem.OldFSTypesExcluded = a.Flag(
		"collector.filesystem.ignored-fs-types",
		"Regexp of filesystem types to ignore for filesystem collector.",
	).Hidden().String()
	config.Filesystem.MountTimeout = a.Flag("collector.filesystem.mount-timeout",
		"how long to wait for a mount to respond before marking it as stale").
		Hidden().Default("5s").Duration()
	config.Filesystem.StatWorkerCount = a.Flag("collector.filesystem.stat-workers",
		"how many stat calls to process simultaneously").
		Hidden().Default("4").Int()

	config.HwMon.ChipInclude = a.Flag("collector.hwmon.chip-include", "Regexp of hwmon chip to include (mutually exclusive to device-exclude).").String()
	config.HwMon.ChipExclude = a.Flag("collector.hwmon.chip-exclude", "Regexp of hwmon chip to exclude (mutually exclusive to device-include).").String()

	config.IPVS.Labels = a.Flag("collector.ipvs.backend-labels", "Comma separated list for IPVS backend stats labels.").String()

	config.NetClass.IgnoredDevices = a.Flag("collector.netclass.ignored-devices", "Regexp of net devices to ignore for netclass collector.").Default("^$").String()
	config.NetClass.InvalidSpeed = a.Flag("collector.netclass.ignore-invalid-speed", "Ignore devices where the speed is invalid. This will be the default behavior in 2.x.").Bool()
	config.NetClass.Netlink = a.Flag("collector.netclass.netlink", "Use netlink to gather stats instead of /proc/net/dev.").Default("false").Bool()
	config.NetClass.RTNLWithStats = a.Flag("collector.netclass_rtnl.with-stats", "Expose the statistics for each network device, replacing netdev collector.").Bool()

	config.NetDev.DeviceInclude = a.Flag("collector.netdev.device-include", "Regexp of net devices to include (mutually exclusive to device-exclude).").String()
	config.NetDev.OldDeviceInclude = a.Flag("collector.netdev.device-whitelist", "DEPRECATED: Use collector.netdev.device-include").Hidden().String()
	config.NetDev.DeviceExclude = a.Flag("collector.netdev.device-exclude", "Regexp of net devices to exclude (mutually exclusive to device-include).").String()
	config.NetDev.OldDeviceExclude = a.Flag("collector.netdev.device-blacklist", "DEPRECATED: Use collector.netdev.device-exclude").Hidden().String()
	config.NetDev.AddressInfo = a.Flag("collector.netdev.address-info", "Collect address-info for every device").Bool()
	config.NetDev.DetailedMetrics = a.Flag("collector.netdev.enable-detailed-metrics", "Use (incompatible) metric names that provide more detailed stats on Linux").Bool()
	config.NetDev.Netlink = a.Flag("collector.netdev.netlink", "Use netlink to gather stats instead of /proc/net/dev.").Default("true").Bool()

	config.NetStat.Fields = a.Flag("collector.netstat.fields", "Regexp of fields to return for netstat collector.").Default("^(.*_(InErrors|InErrs)|Ip_Forwarding|Ip(6|Ext)_(InOctets|OutOctets)|Icmp6?_(InMsgs|OutMsgs)|TcpExt_(Listen.*|Syncookies.*|TCPSynRetrans|TCPTimeouts)|Tcp_(ActiveOpens|InSegs|OutSegs|OutRsts|PassiveOpens|RetransSegs|CurrEstab)|Udp6?_(InDatagrams|OutDatagrams|NoPorts|RcvbufErrors|SndbufErrors))$").String()

	config.NTP.Server = a.Flag("collector.ntp.server", "NTP server to use for ntp collector").Default("127.0.0.1").String()
	config.NTP.ServerPort = a.Flag("collector.ntp.server-port", "UDP port number to connect to on NTP server").Default("123").Int()
	config.NTP.ProtocolVersion = a.Flag("collector.ntp.protocol-version", "NTP protocol version").Default("4").Int()
	config.NTP.ServerIsLocal = a.Flag("collector.ntp.server-is-local", "Certify that collector.ntp.server address is not a public ntp server").Default("false").Bool()
	config.NTP.IPTTL = a.Flag("collector.ntp.ip-ttl", "IP TTL to use while sending NTP query").Default("1").Int()
	// 3.46608s ~ 1.5s + PHI * (1 << maxPoll), where 1.5s is MAXDIST from ntp.org, it is 1.0 in RFC5905
	// max-distance option is used as-is without phi*(1<<poll)
	config.NTP.MaxDistance = a.Flag("collector.ntp.max-distance", "Max accumulated distance to the root").Default("3.46608s").Duration()
	config.NTP.OffsetTolerance = a.Flag("collector.ntp.local-offset-tolerance", "Offset between local clock and local ntpd time to tolerate").Default("1ms").Duration()

	config.Perf.CPUs = a.Flag("collector.perf.cpus", "List of CPUs from which perf metrics should be collected").Default("").String()
	config.Perf.Tracepoint = a.Flag("collector.perf.tracepoint", "perf tracepoint that should be collected").Strings()
	config.Perf.NoHwProfiler = a.Flag("collector.perf.disable-hardware-profilers", "disable perf hardware profilers").Default("false").Bool()
	config.Perf.HwProfiler = a.Flag("collector.perf.hardware-profilers", "perf hardware profilers that should be collected").Strings()
	config.Perf.NoSwProfiler = a.Flag("collector.perf.disable-software-profilers", "disable perf software profilers").Default("false").Bool()
	config.Perf.SwProfiler = a.Flag("collector.perf.software-profilers", "perf software profilers that should be collected").Strings()
	config.Perf.NoCaProfiler = a.Flag("collector.perf.disable-cache-profilers", "disable perf cache profilers").Default("false").Bool()
	config.Perf.CaProfilerFlag = a.Flag("collector.perf.cache-profilers", "perf cache profilers that should be collected").Strings()

	config.PowerSupplyClass.IgnoredPowerSupplies = a.Flag("collector.powersupply.ignored-supplies", "Regexp of power supplies to ignore for powersupplyclass collector.").Default("^$").String()

	config.Qdisc.Fixtures = a.Flag("collector.qdisc.fixtures", "test fixtures to use for qdisc collector end-to-end testing").Default("").String()
	config.Qdisc.DeviceInclude = a.Flag("collector.qdisk.device-include", "Regexp of qdisk devices to include (mutually exclusive to device-exclude).").String()
	config.Qdisc.DeviceExclude = a.Flag("collector.qdisk.device-exclude", "Regexp of qdisk devices to exclude (mutually exclusive to device-include).").String()

	config.Rapl.ZoneLabel = a.Flag("collector.rapl.enable-zone-label", "Enables service unit metric unit_start_time_seconds").Bool()

	config.Runit.ServiceDir = a.Flag("collector.runit.servicedir", "Path to runit service directory.").Default("/etc/service").String()

	config.Stat.Softirq = a.Flag("collector.stat.softirq", "Export softirq calls per vector").Default("false").Bool()

	config.Supervisord.URL = a.Flag("collector.supervisord.url", "XML RPC endpoint.").Default("http://localhost:9001/RPC2").Envar("SUPERVISORD_URL").String()

	config.Sysctl.Include = a.Flag("collector.sysctl.include", "Select sysctl metrics to include").Strings()
	config.Sysctl.IncludeInfo = a.Flag("collector.sysctl.include-info", "Select sysctl metrics to include as info metrics").Strings()

	config.Systemd.UnitInclude = a.Flag("collector.systemd.unit-include", "Regexp of systemd units to include. Units must both match include and not match exclude to be included.").Default(".+").PreAction(func(c *kingpin.ParseContext) error {
		systemdUnitIncludeSet = true
		return nil
	}).String()
	config.Systemd.UnitExclude = kingpin.Flag("collector.systemd.unit-exclude", "Regexp of systemd units to exclude. Units must both match include and not match exclude to be included.").Default(".+\\.(automount|device|mount|scope|slice)").PreAction(func(c *kingpin.ParseContext) error {
		systemdUnitExcludeSet = true
		return nil
	}).String()
	config.Systemd.OldUnitInclude = kingpin.Flag("collector.systemd.unit-whitelist", "DEPRECATED: Use --collector.systemd.unit-include").Hidden().String()
	config.Systemd.OldUnitExclude = kingpin.Flag("collector.systemd.unit-blacklist", "DEPRECATED: Use collector.systemd.unit-exclude").Hidden().String()
	config.Systemd.Private = kingpin.Flag("collector.systemd.private", "Establish a private, direct connection to systemd without dbus (Strongly discouraged since it requires root. For testing purposes only).").Hidden().Bool()
	config.Systemd.EnableTaskMetrics = kingpin.Flag("collector.systemd.enable-task-metrics", "Enables service unit tasks metrics unit_tasks_current and unit_tasks_max").Bool()
	config.Systemd.EnableRestartsMetrics = kingpin.Flag("collector.systemd.enable-restarts-metrics", "Enables service unit metric service_restart_total").Bool()
	config.Systemd.EnableStartTimeMetrics = kingpin.Flag("collector.systemd.enable-start-time-metrics", "Enables service unit metric unit_start_time_seconds").Bool()

	config.Tapestats.IgnoredDevices = kingpin.Flag("collector.tapestats.ignored-devices", "Regexp of devices to ignore for tapestats.").Default("^$").String()

	config.TextFile.Directory = a.Flag("collector.textfile.directory", "Directory to read text files with metrics from.").Default("").String()

	config.VmStat.Fields = a.Flag("collector.vmstat.fields", "Regexp of fields to return for vmstat collector.").Default("^(oom_kill|pgpg|pswp|pg.*fault).*").String()

	config.Wifi.Fixtures = a.Flag("collector.wifi.fixtures", "test fixtures to use for wifi collector end-to-end testing").Default("").String()
	return config
}
