## master / unreleased

### **Breaking changes**

### Changes

* [CHANGE]
* [FEATURE]
* [ENHANCEMENT]
* [BUGFIX] Fix incorrect sysctl call in BSD meminfo collector, resulting in broken swap metrics on FreeBSD #1345
* [BUGFIX] Fix rollover bug in mountstats collector #1364

## 0.18.0 / 2019-05-09

### **Breaking changes**

* Renamed `interface` label to `device` in netclass collector for consistency with
  other network metrics #1224
* The cpufreq metrics now separate the `cpufreq` and `scaling` data based on what the driver provides. #1248
* The labels for the network_up metric have changed, see issue #1236
* Bonding collector now uses `mii_status` instead of `operstatus` #1124
* Several systemd metrics have been turned off by default to improve performance #1254
  These include unit_tasks_current, unit_tasks_max, service_restart_total, and unit_start_time_seconds
* The systemd collector blacklist now includes automount, device, mount, and slice units by default. #1255

### Changes

* [CHANGE] Bonding state uses mii_status #1124
* [CHANGE] Add a limit to the number of in-flight requests #1166
* [CHANGE] Renamed `interface` label to `device` in netclass collector #1224
* [CHANGE] Add separate cpufreq and scaling metrics #1248
* [CHANGE] Several systemd metrics have been turned off by default to improve performance #1254
* [CHANGE] Expand systemd collector blacklist #1255
* [CHANGE] Split cpufreq metrics into a separate collector #1253
* [FEATURE] Add a flag to disable exporter metrics #1148
* [FEATURE] Add kstat-based Solaris metrics for boottime, cpu and zfs collectors #1197
* [FEATURE] Add uname collector for FreeBSD #1239
* [FEATURE] Add diskstats collector for OpenBSD #1250
* [FEATURE] Add pressure collector exposing pressure stall information for Linux #1174
* [FEATURE] Add perf exporter for Linux #1274
* [ENHANCEMENT] Add Infiniband counters #1120
* [ENHANCEMENT] Add TCPSynRetrans to netstat default filter #1143
* [ENHANCEMENT] Move network_up labels into new metric network_info #1236
* [ENHANCEMENT] Use 64-bit counters for Darwin netstat
* [BUGFIX] Add fallback for missing /proc/1/mounts #1172
* [BUGFIX] Fix node_textfile_mtime_seconds to work properly on symlinks #1326

## 0.17.0 / 2018-11-30

Build note: Linux builds can now be built without CGO.

### **Breaking changes**

supvervisord collector reports `start_time_seconds` rather than `uptime` #952

The wifi collector is disabled by default due to suspected caching issues and goroutine leaks.
* https://github.com/prometheus/node_exporter/issues/870
* https://github.com/prometheus/node_exporter/issues/1008

Darwin meminfo metrics have been renamed to match Prometheus conventions. #1060

### Changes

* [CHANGE] Filter out non-installed units when collecting all systemd units #1011
* [CHANGE] `service_restart_total` and `socket_refused_connections_total` will not be reported if you're running an older version of systemd
* [CHANGE] Use /proc/mounts instead of statfs(2) for ro state #1002
* [CHANGE] collector/timex: remove cgo dependency #1079
* [CHANGE] filesystem: Ignore Docker netns mounts #1047
* [CHANGE] Ignore additional virtual filesystems #1104
* [FEATURE] Add netclass collector #851
* [FEATURE] Add processes collector #950
* [FEATURE] Collect start time for systemd units #952
* [FEATURE] Add socket unit stats to systemd collector #968
* [FEATURE] Collect NRestarts property for systemd service units #992
* [FEATURE] Collect NRefused property for systemd socket units (available as of systemd v239) #995
* [FEATURE] Allow removal of rootfs prefix for run in docker #1058
* [ENHANCEMENT] Support for octal characters in mountpoints #954
* [ENHANCEMENT] Update wifi stats to support multiple stations #980
* [ENHANCEMENT] Add transmit/receive bytes total for wifi stations #1150
* [ENHANCEMENT] Handle stuck NFS mounts #997
* [ENHANCEMENT] infiniband: Handle iWARP RDMA modules N/A #974
* [ENHANCEMENT] Update diskstats for linux kernel 4.19 #1109
* [ENHANCEMENT] Collect TasksCurrent, TasksMax per systemd unit #1098

* [BUGFIX] Fix FreeBSD CPU temp #965
* [BUGFIX] Fix goroutine leak in supervisord collector #978
* [BUGFIX] Fix mdadm collector issues #985
* [BUGFIX] Fix ntp collector thread safety #1014
* [BUGFIX] Systemd units will not be ignored if you're running older versions of systemd #1039
* [BUGFIX] Handle vanishing PIDs #1043
* [BUGFIX] Correctly cast Darwin memory info #1060
* [BUGFIX] Filter systemd units in Go for compatibility with older versions #1083
* [BUGFIX] Update cpu collector for OpenBSD 6.4 #1094
* [BUGFIX] Fix typo on HELP of `read_time_seconds_total` #1057
* [BUGFIX] collector/diskstats: don't fail if there are extra stats #1125
* [BUGFIX] collector/hwmon\_linux: handle temperature sensor file #1123
* [BUGFIX] collector/filesystem: add bounds check #1133
* [BUGFIX] Fix dragonfly's CPU counting frequency #1140
* [BUGFIX] Add fallback for missing /proc/1/mounts #1172

## 0.16.0 / 2018-05-15

**Breaking changes**

This release contains major breaking changes to metric names.  Many metrics have new names, labels, and label values in order to conform to current naming conventions.
* Linux node_cpu metrics now break out `guest` values into separate metrics.  See Issue #737
* Many counter metrics have been renamed to include `_total`.
* Many metrics have been renamed/modified to include base units, for example `node_cpu` is now `node_cpu_seconds_total`.

In order to help with the transition we have an [upgrade guide](docs/V0_16_UPGRADE_GUIDE.md).

Other breaking changes:
* The megacli collector has been removed, is now replaced by the storcli.py textfile helper.
* The gmond collector has been removed.
* The textfile collector will now treat timestamps as errors.

* [CHANGE] Split out guest cpu metrics on Linux. #744
* [CHANGE] Exclude Linux proc from filesystem type regexp #774
* [CHANGE] Ignore more virtual filesystems #775
* [CHANGE] Remove obsolete megacli collector. #798
* [CHANGE] Ignore /var/lib/docker by default. #814
* [CHANGE] Cleanup NFS metrics #834
* [CHANGE] Only report core throttles per core, not per cpu #836
* [CHANGE] Treat custom textfile metric timestamps as errors #769
* [CHANGE] Use lowercase cpu label name in interrupts #849
* [CHANGE] Enable bonding collector by default. #872
* [CHANGE] Greatly reduce the metrics vmstat returns by default. #874
* [CHANGE] Greatly trim what netstat collector exposes by default #876
* [CHANGE] Drop `exec_` prefix and move `node_boot_time_seconds` from `exec` to new `boottime` collector and enable for Darwin/Dragonfly/FreeBSD/NetBSD/OpenBSD. #839, #901
* [CHANGE] Remove depreated gmond collector #852
* [CHANGE] align Darwin disk stat names with Linux #930
* [FEATURE] Add `collect[]` parameter #699
* [FEATURE] Add text collector conversion for ipmitool output. #746
* [FEATURE] Add openbsd meminfo #724
* [FEATURE] Add systemd summary metrics #765
* [FEATURE] Add OpenBSD CPU collector #805
* [FEATURE] Add NFS Server metrics collector. #803
* [FEATURE] add sample directory size exporter #789
* [ENHANCEMENT] added Wear_Leveling_Count attribute to smartmon.sh script #707
* [ENHANCEMENT] Simplify Utsname string conversion #716
* [ENHANCEMENT] apt.sh: handle multiple origins in apt-get output #757
* [ENHANCEMENT] Export systemd timers last trigger seconds. #807
* [ENHANCEMENT] updates for zfsonlinux 0.7.5 #779
* [BUGFIX] Fix smartmon.sh textfile script #700
* [BUGFIX] netdev: Change valueType to CounterValue #749
* [BUGFIX] textfile: fix duplicate metrics error #738
* [BUGFIX] Fix panic by updating github.com/ema/qdisc dependency #778
* [BUGFIX] Use uint64 in the ZFS collector #714
* [BUGFIX] multiply page size after float64 coercion to avoid signed integer overflow #780
* [BUGFIX] smartmon: Escape double quotes in device model family #772
* [BUGFIX] Fix log level regression in #533 #815
* [BUGFIX] Correct the ClocksPerSec scaling factor on Darwin #846
* [BUGFIX] Count core throttles per core and per package #871
* [BUGFIX] Fix netdev collector for linux #890 #910
* [BUGFIX] Fix memory corruption when number of filesystems > 16 on FreeBSD #900
* [BUGFIX] Fix parsing of interface aliases in netdev linux #904

## 0.15.2 / 2017-12-06

* [BUGFIX]  cpu: Support processor-less (memory-only) NUMA nodes #734

## 0.15.1 / 2017-11-07

* [BUGFIX] xfs: expose correct fields, fix metric names #708
* [BUGFIX] Correct buffer_bytes > INT_MAX on BSD/amd64. #712
* [BUGFIX] netstat: return nothing when /proc/net/snmp6 not found #718
* [BUGFIX] Fix off by one in Linux interrupts collector #721
* [BUGFIX] Add and use sysReadFile in hwmon collector #728

## 0.15.0 / 2017-10-06

**Breaking changes**

This release contains major breaking changes to flag handling.
* The flag library has been changed, all flags now require double-dashs. (`-foo` becomes `--foo`).
* The collector selection flag has been replaced by individual boolean flags.
* The `-collector.procfs` and `-collector.sysfs` flags have been renamed to `--path.procfs` and `--path.sysfs` respectively.

The `ntp` collector has been replaced with a new NTP-based check that is designed to expose the state of a localhost NTP server rather than provide the offset of the node to a remote NTP server.  By default the `ntp` collector is now locked to localhost.  This is to avoid accidental spamming of public internet NTP pools.

Windows support is now removed, the [wmi_exporter](https://github.com/martinlindhe/wmi_exporter) is recommended as a replacement.

* [CHANGE] `node_cpu` metrics moved from `stats` to `cpu` collector on linux (enabled by default). #548
* [CHANGE] Blacklist systemd scope units #534
* [CHANGE] Remove netbsd/arm #551
* [CHANGE] Remove Windows support #549
* [CHANGE] Enable IPVS collector by default #623
* [CHANGE] Switch to kingpin flags #639
* [CHANGE] Replace --collectors.enabled with per-collector flags #640
* [FEATURE] Add ARP collector for Linux #540
* [FEATURE] Add XFS colector for Linux #568, #575
* [FEATURE] Add qdisc collector for Linux #580
* [FEATURE] Add cpufreq stats for Linux #548
* [FEATURE] Add diskstats for Darwin #593
* [FEATURE] Add bcache collector for Linux #597
* [FEATURE] Add parsing /proc/net/snmp6 file for Linux #615
* [FEATURE] Add timex collector for Linux #664
* [ENHANCEMENT] Include overal health status in smartmon.sh example script #546
* [ENHANCEMENT] Include `guest_nice` in CPU collector #554
* [ENHANCEMENT] Add exec_boot_time for freebsd, dragonfly #550
* [ENHANCEMENT] Get full resolution for node_time #555
* [ENHANCEMENT] infiniband: Multiply port data XMIT/RCV metrics by 4 #579
* [ENHANCEMENT] cpu: Metric 'package_throttles_total' is per package. #657
* [BUGFIX] Fix stale device error metrics #533
* [BUGFIX] edac: Fix typo in node_edac_csrow_uncorrectable_errors_total #564
* [BUGFIX] Use int64 throughout the ZFS collector #653
* [BUGFIX] Silently ignore nonexisting bonding_masters file #569
* [BUGFIX] Change raid0 status line regexp for mdadm collector (bug #618) #619
* [BUGFIX] Ignore wifi collector permission errors #646
* [BUGFIX] Always try to return smartmon_device_info metric #663

## 0.14.0 / 2017-03-21

NOTE: We are deprecating several collectors in this release.
    * `gmond` - Out of scope.
    * `megacli` - Requires forking, to be moved to textfile collection.
    * `ntp` - Out of scope.

Breaking changes:
    * Collector errors are now a separate metric, `node_scrape_collector_success`, not a label on `node_exporter_scrape_duration_seconds` (#516)

* [CHANGE] Report collector success/failure as a bool metric, not a label. #516
* [FEATURE] Add loadavg collector for Solaris #311
* [FEATURE] Add StorCli text collector example script #320
* [FEATURE] Add collector for Linux EDAC #324
* [FEATURE] Add text file utility for SMART metrics #354
* [FEATURE] Add a collector for NFS client statistics. #360
* [FEATURE] Add mountstats collector for detailed NFS statistics #367
* [FEATURE] Add a collector for DRBD #365
* [FEATURE] Add cpu collector for darwin #391
* [FEATURE] Add netdev collector for darwin #393
* [FEATURE] Collect CPU temperatures on FreeBSD #397
* [FEATURE] Add ZFS collector #410
* [FEATURE] Add initial wifi collector #413
* [FEATURE] Add NFS event metrics to mountstats collector #415
* [FEATURE] Add an example rules file #422
* [FEATURE] infiniband: Add new collector for InfiniBand statistics #450
* [FEATURE] buddyinfo: Add support for /proc/buddyinfo for linux free memory fragmentation. #454
* [IMPROVEMENT] hwmon: Provide annotation metric to link chip sysfs paths to human-readable chip types #359
* [IMPROVEMENT] Add node_filesystem_device_errors_total metric #374
* [IMPROVEMENT] Add runit service dir flag #375
* [IMPROVEMENT] Improve Docker documentation #376
* [IMPROVEMENT] Ignore autofs filesystems on linux #384
* [IMPROVEMENT] Replace some FreeBSD collectors with pure Go versions #385
* [IMPROVEMENT] Use filename as label, move 'label' to own metric #411 (hwmon)
* [BUGFIX] mips64 build fix #361
* [BUGFIX] Update vendoring #372 (to fix #242)
* [BUGFIX] Convert remaining collectors to use ConstMetrics #389
* [BUGFIX] Check for errors in netdev scanner #398
* [BUGFIX] Don't leak or race in FreeBSD devstat collector #396
* [BUGFIX] Allow graceful failure in hwmon collector #427
* [BUGFIX] Fix the reporting of active+total disk metrics for inactive raids. #522

## 0.13.0 / 2016-11-26

NOTE: We have disabled builds of linux/ppc64 and linux/ppc64le due to build bugs.

* [FEATURE] Add flag to ignore certain filesystem types (Copy of #217) #241
* [FEATURE] Add NTP stratum to NTP collector. #247
* [FEATURE] Add ignored-units flag for systemd collector #286
* [FEATURE] Compile netdev on dragonfly #314
* [FEATURE] Compile meminfo for dfly #315
* [FEATURE] Add hwmon /sensors support #278
* [FEATURE] Add Linux NUMA "numastat" metrics #249
* [FEATURE] export DragonFlyBSD CPU time #310
* [FEATURE] Dragonfly devstat #323
* [IMPROVEMENT] Use the offset calculation that includes round trip time in the ntp collector #250
* [IMPROVEMENT] Enable `*bsd` collector on darwin #265
* [IMPROVEMENT] Use meminfo_freebsd on darwin as well #266
* [IMPROVEMENT] sockstat: add support for RHE4 #267
* [IMPROVEMENT] Compile fs stats for dfly #302
* [BUGFIX] Add support for raid0 devices in mdadm_linux collector. #253
* [BUGFIX] Close file handler in textfile #263
* [BUGFIX] Ignore partitions on NVME devices by default #268
* [BUGFIX] Fix mdstat tabs parsing #275
* [BUGFIX] Fix mdadm collector for resync=PENDING. #309
* [BUGFIX] mdstat: Fix parsing of RAID0 lines that contain additional attributes. #341
* [BUGFIX] Fix additional mdadm parsing cases #346

## 0.12.0 / 2016-05-05

* [CHANGE] Remove lastlogin collector.
* [CHANGE] Remove -debug.memprofile-file flag.
* [CHANGE] Sync BSD filesystem collector labels with Linux.
* [CHANGE] Remove HTTP Basic Auth support.
* [FEATURE] Add -version flag.
* [FEATURE] Add Linux logind collector.
* [FEATURE] Add Linux ksmd collector.
* [FEATURE] Add Linux memory NUMA collector.
* [FEATURE] Add Linux entropy collector.
* [FEATURE] Add Linux vmstat collector.
* [FEATURE] Add Linux conntrack collector.
* [FEATURE] Add systemd collector.
* [FEATURE] Add OpenBSD support for filesystem, interrupt and netdev collectors.
* [FEATURE] Add supervisord collector.
* [FEATURE] Add Linux /proc/mdstat collector.
* [FEATURE] Add Linux uname collector.
* [FEATURE] Add Linux /proc/sys/fs/file-nr collector.
* [FEATURE] Add Linux /proc/net/sockstat collector.
* [IMPROVEMENT] Provide statically linked Linux binaries.
* [IMPROVEMENT] Remove root requirement for FreeBSD CPU metrics.
* [IMPROVEMENT] Add node_exporter build info metric.
* [IMPROVEMENT] Add disk bytes read/written metrics on Linux.
* [IMPROVEMENT] Add filesystem read-only metric.
* [IMPROVEMENT] Use common Prometheus log formatting.
* [IMPROVEMENT] Add option to specify NTP protocol version.
* [IMPROVEMENT] Handle statfs errors gracefully for individual filesystems.
* [IMPROVEMENT] Add load5 and load15 metrics to loadavg collector.
* [IMPROVEMENT] Add end-to-end tests.
* [IMPROVEMENT] Export FreeBSD CPU metrics to seconds.
* [IMPROVEMENT] Add flag to configure sysfs mountpoint.
* [IMPROVEMENT] Add flag to configure procfs mountpoint.
* [IMPROVEMENT] Add metric for last service state change to runit collector.
* [BUGFIX] Fix FreeBSD netdev metrics on 64 bit systems.
* [BUGFIX] Fix mdstat for devices in delayed resync state.
* [BUGFIX] Fix Linux stat metrics on parallel scrapes.
* [BUGFIX] Remove unavailable collectors from defaults.
* [BUGFIX] Fix build errors on FreeBSD, OpenBSD, Darwin and Windows.
* [BUGFIX] Fix build errors on 386, arm, arm64, ppc64 and ppc64le architectures.
* [BUGFIX] Fix export of stale metrics for removed filesystem and network devices.
* [BUGFIX] textfile: Fix mtime reporting.
* [BUGFIX] megacli: prevent crash when drive temperature is N/A

## 0.11.0 / 2015-07-27

* [FEATURE] Add stats from /proc/net/snmp.
* [FEATURE] Add support for FreeBSD.
* [FEATURE] Allow netdev devices to be ignored.
* [MAINTENANCE] New Dockerfile for unified way to dockerize Prometheus exporters.
* [FEATURE] Add device,fstype collection to the filesystem exporter.
* [IMPROVEMENT] Make logging of collector executions less verbose.

## 0.10.0 / 2015-06-10

* [CHANGE] Change logging output format and flags.

## 0.9.0 / 2015-05-26

* [BUGFIX] Fix `/proc/net/dev` parsing.
* [CLEANUP] Remove the `attributes` collector, use `textfile` instead.
* [CLEANUP] Replace last uses of the configuration file with flags.
* [IMPROVEMENT] Remove cgo dependency.
* [IMPROVEMENT] Sort collector names when printing.
* [FEATURE] IPVS stats collector.

## 0.8.1 / 2015-05-17

* [MAINTENANCE] Use the common Prometheus build infrastructure.
* [MAINTENANCE] Update former Google Code imports.
* [IMPROVEMENT] Log the version at startup.
* [FEATURE] TCP stats collector

## 0.8.0 / 2015-03-09
* [CLEANUP] Introduced semantic versioning and changelog. From now on,
  changes will be reported in this file.
