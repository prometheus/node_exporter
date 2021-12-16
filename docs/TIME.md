# Monitoring time sync with node_exporter

## `ntp` collector

This collector is intended for usage with local NTP daemons including [ntp.org](http://ntp.org/), [chrony](https://chrony.tuxfamily.org/comparison.html), and [OpenNTPD](http://www.openntpd.org/).

Note, some chrony packages have `local stratum 10` configuration value making chrony a valid server when it is unsynchronised. This configuration makes one of the heuristics that derive `node_ntp_sanity` unreliable.

Note, OpenNTPD does not listen for SNTP queries by default. Add `listen on 127.0.0.1` to the OpenNTPD configuration when using this collector with that package.

### `node_ntp_stratum`

This metric shows the [stratum](https://en.wikipedia.org/wiki/Network_Time_Protocol#Clock_strata) of the local NTP daemon.

Stratum `16` means that clock are unsynchronised. See also aforementioned note about default local stratum in chrony.

### `node_ntp_leap`

Raw leap flag value. 0 – OK, 1 – add leap second at UTC midnight, 2 – delete leap second at UTC midnight, 3 – unsynchronised.

OpenNTPD ignores leap seconds and never sets leap flag to `1` or `2`.

### `node_ntp_rtt`

RTT (round-trip time) from node_exporter collector to local NTPD. This value is
used in sanity check as part of causality violation estimate.

### `node_ntp_offset`

[Clock offset](https://en.wikipedia.org/wiki/Network_Time_Protocol#Clock_synchronization_algorithm) between local time and NTPD time.

ntp.org always sets NTPD time to local clock instead of relaying remote NTP
time, so this offset is irrelevant for this NTPD.

This value is used in sanity check as part of causality violation estimate.

### `node_ntp_reference_timestamp_seconds`

Reference Time. This field show time when the last adjustment was made, but
implementation details vary from "**local** wall-clock time" to "Reference Time
field in incoming SNTP packet".

`time() - node_ntp_reference_timestamp_seconds` and
`node_time_seconds - node_ntp_reference_timestamp_seconds` represent some estimate of
"freshness" of synchronization.

### `node_ntp_root_delay` and `node_ntp_root_dispersion`

These values are used to calculate synchronization distance that is limited by
`collector.ntp.max-distance`.

ntp.org adds known local offset to announced root dispersion and linearly
increases dispersion in case of NTP connectivity problems, OpenNTPD does not
account dispersion at all and always reports `0`.

### `node_ntp_sanity`

Aggregate NTPD health including stratum, leap flag, sane freshness, root
distance being less than `collector.ntp.max-distance` and causality violation
being less than `collector.ntp.local-offset-tolerance`.

Causality violation is lower bound estimate of clock error done using SNTP,
it's calculated as positive portion of `abs(node_ntp_offset) - node_ntp_rtt / 2`.

## `timex` collector

This collector exports state of kernel time synchronization flag that should be
maintained by time-keeping daemon and is eventually raised by Linux kernel if
time-keeping daemon does not update it regularly.

Unfortunately some daemons do not handle this flag properly, e.g. chrony-1.30
from Debian/jessie clears `STA_UNSYNC` flag during daemon initialisation and
does not indicate clock synchronization status using this flag. Modern chrony
versions should work better. All chrony versions require `rtcsync` option to
maintain this flag. OpenNTPD does not touch this flag at all till
OpenNTPD-5.9p1.

On the other hand combination of `sync_status` and `offset` exported by `timex`
module is the way to monitor if systemd-timesyncd does its job.
