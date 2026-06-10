# time collector

The time collector exposes metrics about time.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_time_clocksource_available_info | Available clocksources read from '/sys/devices/system/clocksource'. | device, clocksource |
| node_time_clocksource_current_info | Current clocksource read from '/sys/devices/system/clocksource'. | device, clocksource |
| node_time_seconds | System time in seconds since epoch (1970). | n/a |
| node_time_zone_offset_seconds | System time zone offset in seconds. | time_zone |
