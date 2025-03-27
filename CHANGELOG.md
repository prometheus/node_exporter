## master / unreleased

* [CHANGE]
* [FEATURE]
* [ENHANCEMENT]
* [BUGFIX]

## 1.9.1 / 2025-04-01

* [BUGFIX] pressure: Fix missing IRQ on older kernels #3263
* [BUGFIX] Fix Darwin memory leak #3277

## 1.9.0 / 2025-02-17

* [CHANGE] meminfo: Convert linux implementation to use procfs lib #3049
* [CHANGE] Update logging to use Go log/slog #3097
* [FEATURE] filesystem: Add `node_filesystem_mount_info` metric #2970
* [FEATURE] btrfs: Add metrics for commit statistics #3010
* [FEATURE] interrupts: Add collector include/exclude filtering #3028
* [FEATURE] interrupts: Add "exclude zeros" filtering #3028
* [FEATURE] slabinfo: Add filters for slab name. #3041
* [FEATURE] pressure: add IRQ PSI metrics #3048
* [FEATURE] hwmon: Add include and exclude filter for sensors #3072
* [FEATURE] filesystem: Add NetBSD support #3082
* [FEATURE] netdev: Add ifAlias label #3087
* [FEATURE] hwmon: Add Support for GPU Clock Frequencies #3093
* [FEATURE] Add `exclude[]` URL parameter #3116
* [FEATURE] Add AIX support #3136
* [FEATURE] filesystem: Add fs-types/mount-points include flags #3171
* [FEATURE] netstat: Add collector for tcp packet counters for FreeBSD. #3177
* [ENHANCEMENT] ethtool: Add logging for filtering flags #2979
* [ENHANCEMENT] netstat: Add TCPRcvQDrop to default metrics #3021
* [ENHANCEMENT] diskstats: Add block device rotational #3022
* [ENHANCEMENT] cpu: Support CPU online status #3032
* [ENHANCEMENT] arp: optimize interface name resolution #3133
* [ENHANCEMENT] textfile: Allow specifiying multiple directory globs #3135
* [ENHANCEMENT] filesystem: Add reporting of purgeable space on MacOS #3206
* [ENHANCEMENT] ethtool: Skip full scan of NetClass directories #3239
* [BUGFIX] zfs: Prevent `procfs` integer underflow #2961
* [BUGFIX] pressure: Fix collection on systems that do not expose a full CPU stat #3054
* [BUGFIX] cpu: Fix FreeBSD 32-bit host support and plug memory leak #3083
* [BUGFIX] hwmon: Add safety check to hwmon read #3134
* [BUGFIX] zfs: Allow space in dataset name #3186

## 1.8.1 / 2024-05-16

* [BUGFIX] Fix CPU seconds on Solaris #2963
* [BUGFIX] Sign Darwin/MacOS binaries #3008
* [BUGFIX] Fix pressure collector nil reference #3016

## 1.8.0 / 2024-04-24

* [CHANGE] exec_bsd: Fix labels for `vm.stats.sys.v_syscall` sysctl #2895
* [CHANGE] diskstats: Ignore zram devices on linux systems #2898
* [CHANGE] textfile: Avoid inconsistent help-texts  #2962
* [CHANGE] os: Removed caching of modtime/filename of os-release file #2987
* [FEATURE] xfrm: Add new collector #2866
* [FEATURE] watchdog: Add new collector #2880
* [ENHANCEMENT] cpu_vulnerabilities: Add mitigation information label #2806
* [ENHANCEMENT] nfsd: Handle new `wdeleg_getattr` attribute #2810
* [ENHANCEMENT] netstat: Add TCPOFOQueue to default netstat metrics #2867
* [ENHANCEMENT] filesystem: surface device errors #2923
* [ENHANCEMENT] os: Add support end parsing #2982
* [ENHANCEMENT] zfs: Log mib when sysctl read fails on FreeBSD #2975
* [ENHANCEMENT] fibre_channel: update procfs to take into account optional attributes #2933
* [BUGFIX] cpu: Fix debug log in cpu collector #2857
* [BUGFIX] hwmon: Fix hwmon nil ptr #2873
* [BUGFIX] hwmon: Fix hwmon error capture #2915
* [BUGFIX] zfs: Revert "Add ZFS freebsd per dataset stats #2925
* [BUGFIX] ethtool: Sanitize ethtool metric name keys #2940
* [BUGFIX] fix: data race of NetClassCollector metrics initialization #2995

## 1.7.0 / 2023-11-11

* [FEATURE] Add ZFS freebsd per dataset stats #2753
* [FEATURE] Add cpu vulnerabilities reporting from sysfs #2721
* [ENHANCEMENT] Parallelize stat calls in Linux filesystem collector #1772
* [ENHANCEMENT] Add missing linkspeeds to ethtool collector 2711
* [ENHANCEMENT] Add CPU MHz as the value for `node_cpu_info` metric #2778
* [ENHANCEMENT] Improve qdisc collector performance #2779
* [ENHANCEMENT] Add include and exclude filter for hwmon collector #2699
* [ENHANCEMENT] Optionally fetch ARP stats via rtnetlink instead of procfs #2777
* [BUFFIX] Fix ZFS arcstats on FreeBSD 14.0+ 2754
* [BUGFIX] Fallback to 32-bit stats in netdev #2757
* [BUGFIX] Close btrfs.FS handle after use #2780
* [BUGFIX] Move RO status before error return #2807
* [BUFFIX] Fix `promhttp_metric_handler_errors_total` being always active #2808
* [BUGFIX] Fix nfsd v4 index miss #2824

## 1.6.1 / 2023-06-17

Rebuild with latest Go compiler bugfix release.

## 1.6.0 / 2023-05-27

* [CHANGE] Fix cpustat when some cpus are offline #2318
* [CHANGE] Remove metrics of offline CPUs in CPU collector #2605
* [CHANGE] Deprecate ntp collector #2603
* [CHANGE] Remove bcache `cache_readaheads_totals` metrics #2583
* [CHANGE] Deprecate supervisord collector #2685
* [FEATURE] Enable uname collector on NetBSD #2559
* [FEATURE] NetBSD support for the meminfo collector #2570
* [FEATURE] NetBSD support for CPU collector #2626
* [FEATURE] Add FreeBSD collector for netisr subsystem #2668
* [FEATURE] Add softirqs collector #2669
* [ENHANCEMENT] Add suspended as a `node_zfs_zpool_state` #2449
* [ENHANCEMENT] Add administrative state of Linux network interfaces #2515
* [ENHANCEMENT] Log current value of GOMAXPROCS #2537
* [ENHANCEMENT] Add profiler options for perf collector #2542
* [ENHANCEMENT] Allow root path as metrics path #2590
* [ENHANCEMENT] Add cpu frequency governor metrics #2569
* [ENHANCEMENT] Add new landing page #2622
* [ENHANCEMENT] Reduce privileges needed for btrfs device stats #2634
* [ENHANCEMENT] Add ZFS `memory_available_bytes` #2687
* [ENHANCEMENT] Use `SCSI_IDENT_SERIAL` as serial in diskstats #2612
* [ENHANCEMENT] Read missing from netlink netclass attributes from sysfs #2669
* [BUGFIX] perf: fixes for automatically detecting the correct tracefs mountpoints #2553
* [BUGFIX] Fix `thermal_zone` collector noise #2554
* [BUGFIX] Fix a problem fetching the user wire count on FreeBSD #2584
* [BUGFIX] interrupts: Fix fields on linux aarch64 #2631
* [BUGFIX] Remove metrics of offline CPUs in CPU collector #2605
* [BUGFIX] Fix OpenBSD filesystem collector string parsing #2637
* [BUGFIX] Fix bad reporting of `node_cpu_seconds_total` in OpenBSD #2663

## 1.5.0 / 2022-11-29

NOTE: This changes the Go runtime "GOMAXPROCS" to 1. This is done to limit the
  concurrency of the exporter to 1 CPU thread at a time in order to avoid a
  race condition problem in the Linux kernel (#2500) and parallel IO issues
  on nodes with high numbers of CPUs/CPU threads (#1880).

NOTE: A command line arg has been changed from `--web.config` to `--web.config.file`.

* [CHANGE] Default GOMAXPROCS to 1 #2530
* [FEATURE] Add multiple listeners and systemd socket listener activation #2393
* [ENHANCEMENT] Add RTNL version of netclass collector #2492, #2528
* [BUGFIX] Fix diskstats exclude flags #2487
* [BUGFIX] Bump go/x/crypt and go/x/net #2488
* [BUGFIX] Fix hwmon label sanitizer #2504
* [BUGFIX] Use native endianness when encoding InetDiagMsg #2508
* [BUGFIX] Fix btrfs device stats always being zero #2516
* [BUGFIX] Security: Update exporter-toolkit (CVE-2022-46146) #2531

## 1.4.1 / 2022-11-29

* [BUGFIX] Fix diskstats exclude flags #2487
* [BUGFIX] Security: Update go/x/crypto and go/x/net (CVE-2022-27191 CVE-2022-27664) #2488
* [BUGFIX] Security: Update exporter-toolkit (CVE-2022-46146) #2531

## 1.4.0 / 2022-09-24

* [CHANGE] Merge metrics descriptions in textfile collector #2475
* [FEATURE] [node-mixin] Add darwin dashboard to mixin #2351
* [FEATURE] Add "isolated" metric on cpu collector on linux #2251
* [FEATURE] Add cgroup summary collector #2408
* [FEATURE] Add selinux collector #2205
* [FEATURE] Add slab info collector #2376
* [FEATURE] Add sysctl collector #2425
* [FEATURE] Also track the CPU Spin time for OpenBSD systems #1971
* [FEATURE] Add support for MacOS version #2471
* [ENHANCEMENT] [node-mixin] Add missing selectors #2426
* [ENHANCEMENT] [node-mixin] Change current datasource to grafana's default #2281
* [ENHANCEMENT] [node-mixin] Change disk graph to disk table #2364
* [ENHANCEMENT] [node-mixin] Change io time units to %util #2375
* [ENHANCEMENT] Ad user_wired_bytes and laundry_bytes on *bsd #2266
* [ENHANCEMENT] Add additional vm_stat memory metrics for darwin #2240
* [ENHANCEMENT] Add device filter flags to arp collector #2254
* [ENHANCEMENT] Add diskstats include and exclude device flags #2417
* [ENHANCEMENT] Add node_softirqs_total metric #2221
* [ENHANCEMENT] Add rapl zone name label option #2401
* [ENHANCEMENT] Add slabinfo collector #1799
* [ENHANCEMENT] Allow user to select port on NTP server to query #2270
* [ENHANCEMENT] collector/diskstats: Add labels and metrics from udev #2404
* [ENHANCEMENT] Enable builds against older macOS SDK #2327
* [ENHANCEMENT] qdisk-linux: Add exclude and include flags for interface name #2432
* [ENHANCEMENT] systemd: Expose systemd minor version #2282
* [ENHANCEMENT] Use netlink for tcpstat collector #2322
* [ENHANCEMENT] Use netlink to get netdev stats #2074
* [ENHANCEMENT] Add additional perf counters for stalled frontend/backend cycles #2191
* [ENHANCEMENT] Add btrfs device error stats #2193
* [BUGFIX] [node-mixin] Fix fsSpaceAvailableCriticalThreshold and fsSpaceAvailableWarning #2352
* [BUGFIX] Fix concurrency issue in ethtool collector #2289
* [BUGFIX] Fix concurrency issue in netdev collector #2267
* [BUGFIX] Fix diskstat reads and write metrics for disks with different sector sizes #2311
* [BUGFIX] Fix iostat on macos broken by deprecation warning #2292
* [BUGFIX] Fix NodeFileDescriptorLimit alerts #2340
* [BUGFIX] Sanitize rapl zone names #2299
* [BUGFIX] Add file descriptor close safely in test #2447
* [BUGFIX] Fix race condition in os_release.go #2454
* [BUGFIX] Skip ZFS IO metrics if their paths are missing #2451

## 1.3.1 / 2021-12-01

* [BUGFIX] Handle nil CPU thermal power status on M1 #2218
* [BUGFIX] bsd: Ignore filesystems flagged as MNT_IGNORE. #2227
* [BUGFIX] Sanitize UTF-8 in dmi collector #2229

## 1.3.0 / 2021-10-20

NOTE: In order to support globs in the textfile collector path, filenames exposed by
      `node_textfile_mtime_seconds` now contain the full path name.

* [CHANGE] Add path label to rapl collector #2146
* [CHANGE] Exclude filesystems under /run/credentials #2157
* [CHANGE] Add TCPTimeouts to netstat default filter #2189
* [FEATURE] Add lnstat collector for metrics from /proc/net/stat/ #1771
* [FEATURE] Add darwin powersupply collector #1777
* [FEATURE] Add support for monitoring GPUs on Linux #1998
* [FEATURE] Add Darwin thermal collector #2032
* [FEATURE] Add os release collector #2094
* [FEATURE] Add netdev.address-info collector #2105
* [FEATURE] Add clocksource metrics to time collector #2197
* [ENHANCEMENT] Support glob textfile collector directories #1985
* [ENHANCEMENT] ethtool: Expose node_ethtool_info metric #2080
* [ENHANCEMENT] Use include/exclude flags for ethtool filtering #2165
* [ENHANCEMENT] Add flag to disable guest CPU metrics #2123
* [ENHANCEMENT] Add DMI collector #2131
* [ENHANCEMENT] Add threads metrics to processes collector #2164
* [ENHANCEMENT] Reduce timer GC delays in the Linux filesystem collector #2169
* [ENHANCEMENT] Add TCPTimeouts to netstat default filter #2189
* [ENHANCEMENT] Use SysctlTimeval for boottime collector on BSD #2208
* [BUGFIX] ethtool: Sanitize metric names #2093
* [BUGFIX] Fix ethtool collector for multiple interfaces #2126
* [BUGFIX] Fix possible panic on macOS #2133
* [BUGFIX] Collect flag_info and bug_info only for one core #2156
* [BUGFIX] Prevent duplicate ethtool metric names #2187

## 1.2.2 / 2021-08-06

* [BUGFIX] Fix processes collector long int parsing #2112

## 1.2.1 / 2021-07-23

* [BUGFIX] Fix zoneinfo parsing prometheus/procfs#386
* [BUGFIX] Fix nvme collector log noise #2091
* [BUGFIX] Fix rapl collector log noise #2092

## 1.2.0 / 2021-07-15

NOTE: Ignoring invalid network speed will be the default in 2.x
NOTE: Filesystem collector flags have been renamed. `--collector.filesystem.ignored-mount-points` is now `--collector.filesystem.mount-points-exclude` and `--collector.filesystem.ignored-fs-types` is now `--collector.filesystem.fs-types-exclude`. The old flags will be removed in 2.x.

* [CHANGE] Rename filesystem collector flags to match other collectors #2012
* [CHANGE] Make node_exporter print usage to STDOUT #2039
* [FEATURE] Add conntrack statistics metrics #1155
* [FEATURE] Add ethtool stats collector #1832
* [FEATURE] Add flag to ignore network speed if it is unknown #1989
* [FEATURE] Add tapestats collector for Linux #2044
* [FEATURE] Add nvme collector #2062
* [ENHANCEMENT] Add ErrorLog plumbing to promhttp #1887
* [ENHANCEMENT] Add more Infiniband counters #2019
* [ENHANCEMENT] netclass: retrieve interface names and filter before parsing #2033
* [ENHANCEMENT] Add time zone offset metric #2060
* [BUGFIX] Handle errors from disabled PSI subsystem #1983
* [BUGFIX] Fix panic when using backwards compatible flags #2000
* [BUGFIX] Fix wrong value for OpenBSD memory buffer cache #2015
* [BUGFIX] Only initiate collectors once #2048
* [BUGFIX] Handle small backwards jumps in CPU idle #2067

## 1.1.2 / 2021-03-05

* [BUGFIX] Handle errors from disabled PSI subsystem #1983
* [BUGFIX] Sanitize strings from /sys/class/power_supply #1984
* [BUGFIX] Silence missing netclass errors #1986

## 1.1.1 / 2021-02-12

* [BUGFIX] Fix ineffassign issue #1957
* [BUGFIX] Fix some noisy log lines #1962

## 1.1.0 / 2021-02-05

NOTE: We have improved some of the flag naming conventions (PR #1743). The old names are
      deprecated and will be removed in 2.0. They will continue to work for backwards
      compatibility.

* [CHANGE] Improve filter flag names #1743
* [CHANGE] Add btrfs and powersupplyclass to list of exporters enabled by default #1897
* [FEATURE] Add fibre channel collector #1786
* [FEATURE] Expose cpu bugs and flags as info metrics. #1788
* [FEATURE] Add network_route collector #1811
* [FEATURE] Add zoneinfo collector #1922
* [ENHANCEMENT] Add more InfiniBand counters #1694
* [ENHANCEMENT] Add flag to aggr ipvs metrics to avoid high cardinality metrics #1709
* [ENHANCEMENT] Adding backlog/current queue length to qdisc collector #1732
* [ENHANCEMENT] Include TCP OutRsts in netstat metrics #1733
* [ENHANCEMENT] Add pool size to entropy collector #1753
* [ENHANCEMENT] Remove CGO dependencies for OpenBSD amd64 #1774
* [ENHANCEMENT] bcache: add writeback_rate_debug stats #1658
* [ENHANCEMENT] Add check state for mdadm arrays via node_md_state metric #1810
* [ENHANCEMENT] Expose XFS inode statistics #1870
* [ENHANCEMENT] Expose zfs zpool state #1878
* [ENHANCEMENT] Added an ability to pass collector.supervisord.url via SUPERVISORD_URL environment variable #1947
* [BUGFIX] filesystem_freebsd: Fix label values #1728
* [BUGFIX] Fix various procfs parsing errors #1735
* [BUGFIX] Handle no data from powersupplyclass #1747
* [BUGFIX] udp_queues_linux.go: change upd to udp in two error strings #1769
* [BUGFIX] Fix node_scrape_collector_success behaviour #1816
* [BUGFIX] Fix NodeRAIDDegraded to not use a string rule expressions #1827
* [BUGFIX] Fix node_md_disks state label from fail to failed #1862
* [BUGFIX] Handle EPERM for syscall in timex collector #1938
* [BUGFIX] bcache: fix typo in a metric name #1943
* [BUGFIX] Fix XFS read/write stats (https://github.com/prometheus/procfs/pull/343)

## 1.0.1 / 2020-06-15

* [BUGFIX] filesystem_freebsd: Fix label values #1728
* [BUGFIX] Update prometheus/procfs to fix log noise #1735
* [BUGFIX] Fix build tags for collectors #1745
* [BUGFIX] Handle no data from powersupplyclass #1747, #1749

## 1.0.0 / 2020-05-25

### **Breaking changes**

* The netdev collector CLI argument `--collector.netdev.ignored-devices` was renamed to `--collector.netdev.device-blacklist` in order to conform with the systemd collector. #1279
* The label named `state` on `node_systemd_service_restart_total` metrics was changed to `name` to better describe the metric. #1393
* Refactoring of the mdadm collector changes several metrics
    - `node_md_disks_active` is removed
    - `node_md_disks` now has a `state` label for "failed", "spare", "active" disks.
    - `node_md_is_active` is replaced by `node_md_state` with a state set of "active", "inactive", "recovering", "resync".
* Additional label `mountaddr` added to NFS device metrics to distinguish mounts from the same URL, but different IP addresses. #1417
* Metrics node_cpu_scaling_frequency_min_hrts and node_cpu_scaling_frequency_max_hrts of the cpufreq collector were renamed to node_cpu_scaling_frequency_min_hertz and node_cpu_scaling_frequency_max_hertz. #1510
* Collectors that are enabled, but are unable to find data to collect, now return 0 for `node_scrape_collector_success`.

### Changes

* [CHANGE] Add `--collector.netdev.device-whitelist`. #1279
* [CHANGE] Ignore iso9600 filesystem on Linux #1355
* [CHANGE] Refactor mdadm collector #1403
* [CHANGE] Add `mountaddr` label to NFS metrics. #1417
* [CHANGE] Don't count empty collectors as success. #1613
* [FEATURE] New flag to disable default collectors #1276
* [FEATURE] Add experimental TLS support #1277, #1687, #1695
* [FEATURE] Add collector for Power Supply Class #1280
* [FEATURE] Add new schedstat collector #1389
* [FEATURE] Add FreeBSD zfs support #1394
* [FEATURE] Add uname support for Darwin and OpenBSD #1433
* [FEATURE] Add new metric node_cpu_info #1489
* [FEATURE] Add new thermal_zone collector #1425
* [FEATURE] Add new cooling_device metrics to thermal zone collector #1445
* [FEATURE] Add swap usage on darwin #1508
* [FEATURE] Add Btrfs collector #1512
* [FEATURE] Add RAPL collector #1523
* [FEATURE] Add new softnet collector #1576
* [FEATURE] Add new udp_queues collector #1503
* [FEATURE] Add basic authentication #1673
* [ENHANCEMENT] Log pid when there is a problem reading the process stats #1341
* [ENHANCEMENT] Collect InfiniBand port state and physical state #1357
* [ENHANCEMENT] Include additional XFS runtime statistics. #1423
* [ENHANCEMENT] Report non-fatal collection errors in the exporter metric. #1439
* [ENHANCEMENT] Expose IPVS firewall mark as a label #1455
* [ENHANCEMENT] Add check for systemd version before attempting to query certain metrics. #1413
* [ENHANCEMENT] Add a flag to adjust mount timeout #1486
* [ENHANCEMENT] Add new counters for flush requests in Linux 5.5 #1548
* [ENHANCEMENT] Add metrics and tests for UDP receive and send buffer errors #1534
* [ENHANCEMENT] The sockstat collector now exposes IPv6 statistics in addition to the existing IPv4 support. #1552
* [ENHANCEMENT] Add infiniband info metric #1563
* [ENHANCEMENT] Add unix socket support for supervisord collector #1592
* [ENHANCEMENT] Implement loadavg on all BSDs without cgo #1584
* [ENHANCEMENT] Add model_name and stepping to node_cpu_info metric #1617
* [ENHANCEMENT] Add `--collector.perf.cpus` to allow setting the CPU list for perf stats. #1561
* [ENHANCEMENT] Add metrics for IO errors and retires on Darwin. #1636
* [ENHANCEMENT] Add perf tracepoint collection flag #1664
* [ENHANCEMENT] ZFS: read contents of objset file #1632
* [ENHANCEMENT] Linux CPU: Cache CPU metrics to make them monotonically increasing #1711
* [BUGFIX] Read /proc/net files with a single read syscall #1380
* [BUGFIX] Renamed label `state` to `name` on `node_systemd_service_restart_total`. #1393
* [BUGFIX] Fix netdev nil reference on Darwin #1414
* [BUGFIX] Strip path.rootfs from mountpoint labels #1421
* [BUGFIX] Fix seconds reported by schedstat #1426
* [BUGFIX] Fix empty string in path.rootfs #1464
* [BUGFIX] Fix typo in cpufreq metric names #1510
* [BUGFIX] Read /proc/stat in one syscall #1538
* [BUGFIX] Fix OpenBSD cache memory information #1542
* [BUGFIX] Refactor textfile collector to avoid looping defer #1549
* [BUGFIX] Fix network speed math #1580
* [BUGFIX] collector/systemd: use regexp to extract systemd version #1647
* [BUGFIX] Fix initialization in perf collector when using multiple CPUs #1665
* [BUGFIX] Fix accidentally empty lines in meminfo_linux #1671

## 0.18.1 / 2019-06-04

### Changes
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

* [CHANGE] Use /proc/mounts instead of statfs(2) for ro state #1002
* [CHANGE] Exclude only subdirectories of /var/lib/docker #1003
* [CHANGE] Filter out non-installed units when collecting all systemd units #1011
* [CHANGE] `service_restart_total` and `socket_refused_connections_total` will not be reported if you're running an older version of systemd
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
* [CHANGE] Remove deprecated gmond collector #852
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
* The flag library has been changed, all flags now require double-dashes. (`-foo` becomes `--foo`).
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
* [FEATURE] Add XFS collector for Linux #568, #575
* [FEATURE] Add qdisc collector for Linux #580
* [FEATURE] Add cpufreq stats for Linux #548
* [FEATURE] Add diskstats for Darwin #593
* [FEATURE] Add bcache collector for Linux #597
* [FEATURE] Add parsing /proc/net/snmp6 file for Linux #615
* [FEATURE] Add timex collector for Linux #664
* [ENHANCEMENT] Include overall health status in smartmon.sh example script #546
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
