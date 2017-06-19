## master

* [CHANGE] `node_cpu` metrics moved from `stats` to `cpu` collector on linux (enabled by default).

## v0.14.0 / 2017-03-21

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

## v0.13.0 / 2016-11-26

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
