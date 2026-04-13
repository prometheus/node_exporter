# watchdog collector

The watchdog collector exposes metrics about watchdog.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_watchdog_access_cs0 | Value of /sys/class/watchdog/<watchdog>/access_cs0 | name |
| node_watchdog_bootstatus | Value of /sys/class/watchdog/<watchdog>/bootstatus | name |
| node_watchdog_fw_version | Value of /sys/class/watchdog/<watchdog>/fw_version | name |
| node_watchdog_info | Info of /sys/class/watchdog/<watchdog> | name, options, identity, state, status, pretimeout_governor |
| node_watchdog_nowayout | Value of /sys/class/watchdog/<watchdog>/nowayout | name |
| node_watchdog_pretimeout_seconds | Value of /sys/class/watchdog/<watchdog>/pretimeout | name |
| node_watchdog_timeleft_seconds | Value of /sys/class/watchdog/<watchdog>/timeleft | name |
| node_watchdog_timeout_seconds | Value of /sys/class/watchdog/<watchdog>/timeout | name |
