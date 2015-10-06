# Node exporter

[![Build Status](https://travis-ci.org/prometheus/node_exporter.svg)](https://travis-ci.org/prometheus/node_exporter)

Prometheus exporter for machine metrics, written in Go with pluggable metric
collectors.

## Building and running

    make
    ./node_exporter <flags>

## Running tests

    make test

## Available collectors

By default the build will include the native collectors that expose information
from `/proc`.

Which collectors are used is controlled by the `--collectors.enabled` flag.

### Enabled by default

Name     | Description
---------|------------
diskstats | Exposes disk I/O statistics from `/proc/diskstats`.
filesystem | Exposes filesystem statistics, such as disk space used.
loadavg | Exposes load average.
meminfo | Exposes memory statistics from `/proc/meminfo`.
netdev | Exposes network interface statistics from `/proc/netstat`, such as bytes transferred.
netstat | Exposes network statistics from `/proc/net/netstat`. This is the same information as `netstat -s`.
stat | Exposes various statistics from `/proc/stat`. This includes CPU usage, boot time, forks and interrupts.
textfile | Exposes statistics read from local disk. The `--collector.textfile.directory` flag must be set.
time | Exposes the current system time.
mdadm | Exposes statistics about devices in `/proc/mdstat` (does nothing if no /proc/mdstat present).


### Disabled by default

Name     | Description
---------|------------
bonding | Exposes the number of configured and active slaves of Linux bonding interfaces.
gmond | Exposes statistics from Ganglia.
interrupts | Exposes detailed interrupts statistics from `/proc/interrupts`.
ipvs | Exposes IPVS status from `/proc/net/ip_vs` and stats from `/proc/net/ip_vs_stats`.
lastlogin | Exposes the last time there was a login.
megacli | Exposes RAID statistics from MegaCLI.
ntp | Exposes time drift from an NTP server.
runit | Exposes service status from [runit](http://smarden.org/runit/).
supervisord | Exposes service status from [supervisord](http://supervisord.org/).
tcpstat | Exposes TCP connection status information from `/proc/net/tcp` and `/proc/net/tcp6`. (Warning: the current version has potential performance issues in high load situations.)

## Textfile Collector

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

## Using Docker

You can deploy this exporter using the [prom/node-exporter](https://registry.hub.docker.com/u/prom/node-exporter/) Docker image.

For example:

```bash
docker pull prom/node-exporter

docker run -d -p 9100:9100 --net="host" prom/node-exporter
```
