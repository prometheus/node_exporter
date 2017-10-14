# Node exporter [![Build Status](https://travis-ci.org/prometheus/node_exporter.svg)][travis]

[![CircleCI](https://circleci.com/gh/prometheus/node_exporter/tree/master.svg?style=shield)][circleci]
[![Buildkite status](https://badge.buildkite.com/94a0c1fb00b1f46883219c256efe9ce01d63b6505f3a942f9b.svg)](https://buildkite.com/prometheus/node-exporter)
[![Docker Repository on Quay](https://quay.io/repository/prometheus/node-exporter/status)][quay]
[![Docker Pulls](https://img.shields.io/docker/pulls/prom/node-exporter.svg?maxAge=604800)][hub]
[![Go Report Card](https://goreportcard.com/badge/github.com/prometheus/node_exporter)][goreportcard]

Prometheus exporter for hardware and OS metrics exposed by \*NIX kernels, written
in Go with pluggable metric collectors.

The [WMI exporter](https://github.com/martinlindhe/wmi_exporter) is recommended for Windows users.

## Collectors

There is varying support for collectors on each operating system. The tables
below list all existing collectors and the supported systems.

Collectors are enabled by providing a `--collector.<name>` flag.
Collectors that are enabled by default can be disabled by providing a `--no-collector.<name>` flag.

### Enabled by default

Name     | Description | OS
---------|-------------|----
arp | Exposes ARP statistics from `/proc/net/arp`. | Linux
bcache | Exposes bcache statistics from `/sys/fs/bcache/`. | Linux
conntrack | Shows conntrack statistics (does nothing if no `/proc/sys/net/netfilter/` present). | Linux
cpu | Exposes CPU statistics | Darwin, Dragonfly, FreeBSD, Linux
diskstats | Exposes disk I/O statistics. | Darwin, Linux
edac | Exposes error detection and correction statistics. | Linux
entropy | Exposes available entropy. | Linux
exec | Exposes execution statistics. | Dragonfly, FreeBSD
filefd | Exposes file descriptor statistics from `/proc/sys/fs/file-nr`. | Linux
filesystem | Exposes filesystem statistics, such as disk space used. | Darwin, Dragonfly, FreeBSD, Linux, OpenBSD
hwmon | Expose hardware monitoring and sensor data from `/sys/class/hwmon/`. | Linux
infiniband | Exposes network statistics specific to InfiniBand and Intel OmniPath configurations. | Linux
ipvs | Exposes IPVS status from `/proc/net/ip_vs` and stats from `/proc/net/ip_vs_stats`. | Linux
loadavg | Exposes load average. | Darwin, Dragonfly, FreeBSD, Linux, NetBSD, OpenBSD, Solaris
mdadm | Exposes statistics about devices in `/proc/mdstat` (does nothing if no `/proc/mdstat` present). | Linux
meminfo | Exposes memory statistics. | Darwin, Dragonfly, FreeBSD, Linux
netdev | Exposes network interface statistics such as bytes transferred. | Darwin, Dragonfly, FreeBSD, Linux, OpenBSD
netstat | Exposes network statistics from `/proc/net/netstat`. This is the same information as `netstat -s`. | Linux
sockstat | Exposes various statistics from `/proc/net/sockstat`. | Linux
stat | Exposes various statistics from `/proc/stat`. This includes boot time, forks and interrupts. | Linux
textfile | Exposes statistics read from local disk. The `--collector.textfile.directory` flag must be set. | _any_
time | Exposes the current system time. | _any_
timex | Exposes selected adjtimex(2) system call stats. | Linux
uname | Exposes system information as provided by the uname system call. | Linux
vmstat | Exposes statistics from `/proc/vmstat`. | Linux
wifi | Exposes WiFi device and station statistics. | Linux
xfs | Exposes XFS runtime statistics. | Linux (kernel 4.4+)
zfs | Exposes [ZFS](http://open-zfs.org/) performance statistics. | [Linux](http://zfsonlinux.org/)

### Disabled by default

Name     | Description | OS
---------|-------------|----
bonding | Exposes the number of configured and active slaves of Linux bonding interfaces. | Linux
buddyinfo | Exposes statistics of memory fragments as reported by /proc/buddyinfo. | Linux
devstat | Exposes device statistics | Dragonfly, FreeBSD
drbd | Exposes Distributed Replicated Block Device statistics (to version 8.4) | Linux
interrupts | Exposes detailed interrupts statistics. | Linux, OpenBSD
ksmd | Exposes kernel and system statistics from `/sys/kernel/mm/ksm`. | Linux
logind | Exposes session counts from [logind](http://www.freedesktop.org/wiki/Software/systemd/logind/). | Linux
meminfo\_numa | Exposes memory statistics from `/proc/meminfo_numa`. | Linux
mountstats | Exposes filesystem statistics from `/proc/self/mountstats`. Exposes detailed NFS client statistics. | Linux
nfs | Exposes NFS client statistics from `/proc/net/rpc/nfs`. This is the same information as `nfsstat -c`. | Linux
ntp | Exposes local NTP daemon health to check [time](./docs/TIME.md) | _any_
qdisc | Exposes [queuing discipline](https://en.wikipedia.org/wiki/Network_scheduler#Linux_kernel) statistics | Linux
runit | Exposes service status from [runit](http://smarden.org/runit/). | _any_
supervisord | Exposes service status from [supervisord](http://supervisord.org/). | _any_
systemd | Exposes service and system status from [systemd](http://www.freedesktop.org/wiki/Software/systemd/). | Linux
tcpstat | Exposes TCP connection status information from `/proc/net/tcp` and `/proc/net/tcp6`. (Warning: the current version has potential performance issues in high load situations.) | Linux

### Deprecated

*These collectors will be (re)moved in the future.*

Name     | Description | OS
---------|-------------|----
gmond | Exposes statistics from Ganglia. | _any_
megacli | Exposes RAID statistics from MegaCLI. | Linux

### Textfile Collector

The textfile collector is similar to the [Pushgateway](https://github.com/prometheus/pushgateway),
in that it allows exporting of statistics from batch jobs. It can also be used
to export static metrics, such as what role a machine has. The Pushgateway
should be used for service-level metrics. The textfile module is for metrics
that are tied to a machine.

To use it, set the `--collector.textfile.directory` flag on the Node exporter. The
collector will parse all files in that directory matching the glob `*.prom`
using the [text
format](http://prometheus.io/docs/instrumenting/exposition_formats/).

To atomically push completion time for a cron job:
```
echo my_batch_job_completion_time $(date +%s) > /path/to/directory/my_batch_job.prom.$$
mv /path/to/directory/my_batch_job.prom.$$ /path/to/directory/my_batch_job.prom
```

To statically set roles for a machine using labels:
```
echo 'role{role="application_server"} 1' > /path/to/directory/role.prom.$$
mv /path/to/directory/role.prom.$$ /path/to/directory/role.prom
```

### Filtering enabled collectors

The node_exporter will expose all metrics from enabled collectors by default, but it can be passed an optional list of collectors to filter metrics. The `collect[]` parameter accepts values matching enabled collector names.

This can be useful for specifying different scrape intervals for different collectors in Prometheus:

```yaml
scrape_configs:
  - job_name: 'node resources'
    scrape_interval: 15s
    static_configs:
      - targets:
        - '192.168.1.2:9100'
    params:
      collect[]:
        - cpu
        - meminfo
        - diskstats
        - netdev
        - netstat

  - job_name: 'node storage'
    scrape_interval: 1m
    static_configs:
      - targets:
        - '192.168.1.2:9100'
    params:
      collect[]:
        - filefd
        - filesystem
        - xfs
```

## Building and running

Prerequisites:

* [Go compiler](https://golang.org/dl/)
* RHEL/CentOS: `glibc-static` package.

Building:

    go get github.com/prometheus/node_exporter
    cd ${GOPATH-$HOME/go}/src/github.com/prometheus/node_exporter
    make
    ./node_exporter <flags>

To see all available configuration flags:

    ./node_exporter -h

## Running tests

    make test


## Using Docker
The node\_exporter is designed to monitor the host system. It's not recommended
to deploy it as Docker container because it requires access to the host system.
If you need to run it on Docker, you can deploy this exporter using the
[node-exporter Docker
image](https://quay.io/repository/prometheus/node-exporter) with the following
options and bind-mounts:

for docker version >= 1.10
```bash
docker run -d -p 9100:9100 \
  -v "/proc:/host/proc:ro" \
  -v "/sys:/host/sys:ro" \
  -v "/:/rootfs:ro,rslave" \
  --net="host" \
  quay.io/prometheus/node-exporter \
    --path.procfs /host/proc \
    --path.sysfs /host/sys \
    --collector.filesystem.ignored-mount-points "^/(sys|proc|dev|host|etc)($|/)"
```

for docker version < 1.10
```bash
docker run -d -p 9100:9100 \
  -v "/proc:/host/proc:ro" \
  -v "/sys:/host/sys:ro" \
  -v "/:/rootfs:ro" \
  --net="host" \
  quay.io/prometheus/node-exporter \
    --path.procfs /host/proc \
    --path.sysfs /host/sys \
    --collector.filesystem.ignored-mount-points "^/(sys|proc|dev|host|etc)($|/)"
```

Be aware though that the mountpoint label in various metrics will now have
`/rootfs` as prefix.

## Using a third-party repository for RHEL/CentOS/Fedora

There is a [community-supplied COPR repository](https://copr.fedorainfracloud.org/coprs/ibotty/prometheus-exporters/). It closely follows upstream releases.

[travis]: https://travis-ci.org/prometheus/node_exporter
[hub]: https://hub.docker.com/r/prom/node-exporter/
[circleci]: https://circleci.com/gh/prometheus/node_exporter
[quay]: https://quay.io/repository/prometheus/node-exporter
[goreportcard]: https://goreportcard.com/report/github.com/prometheus/node_exporter
