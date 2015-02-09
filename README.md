# Node exporter

Prometheus exporter for machine metrics, written in Go with pluggable metric
collectors.

## Available collectors

By default the build will include the native collectors that expose information
from /proc.

Which collectors are used is controlled by the `--enabledCollectors` flag.

### Enabled by default

Name     | Description
---------|------------
attributes | Exposes attributes from the configuration file. Deprecated, use textfile module instead.
diskstats | Exposes disk I/O statistics from /proc/diskstats.
filesystem | Exposes filesystem statistics, such as disk space used.
loadavg | Exposes load average.
meminfo | Exposes memory statistics from /proc/meminfo.
netdev | Exposes network interface statistics from /proc/netstat, such as bytes transferred.
netstat | Exposes network statistics from /proc/net/netstat. This is the same information as `netstat -s`.
stat | Exposes various statistics from /proc/stat. This includes CPU usage, boot time, forks and interrupts.
textfile | Exposes statistics read from local disk. The `--textfile.directory` flag must be set.
time | Exposes the current system time.


### Disabled by default

Name     | Description
---------|------------
bonding | Exposes the number of configured and active slaves of Linux bonding interfaces.
gmond | Exposes statistics from Ganglia.
interrupts | Exposes detailed interrupts statistics from /proc/interrupts.
lastlogin | Exposes the last time there was a login.
megacli | Exposes RAID statistics from MegaCLI.
ntp | Exposes time drift from an NTP server.
runit | Exposes service status from [runit](http://smarden.org/runit/).

## Textfile Collector

The textfile collector is similar to the [Pushgateway](https://github.com/prometheus/pushgateway),
in that it allows exporting of statistics from batch jobs. It can also be used
to export static metrics, such as what role a machine has. The Pushgateway
should be used for service-level metrics. The textfile module is for metrics
that are tied to a machine.

To use set the `--textfile.directory` flag on the Node exporter. The collector
will pares all files in that directory matching the glob `*.prom` using the
[text format](http://prometheus.io/docs/instrumenting/exposition_formats/).

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
