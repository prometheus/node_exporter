# systemd collector

The systemd collector exposes metrics about systemd.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.systemd.enable-restarts-metrics | Enables service unit metric service_restart_total |  |
| collector.systemd.enable-start-time-metrics | Enables service unit metric unit_start_time_seconds |  |
| collector.systemd.enable-task-metrics | Enables service unit tasks metrics unit_tasks_current and unit_tasks_max |  |
| collector.systemd.private | Establish a private, direct connection to systemd without dbus (Strongly discouraged since it requires root. For testing purposes only). |  |
| collector.systemd.unit-blacklist | DEPRECATED: Use collector.systemd.unit-exclude |  |
| collector.systemd.unit-exclude | Regexp of systemd units to exclude. Units must both match include and not match exclude to be included. | .+\\.(automount|device|mount|scope|slice) |
| collector.systemd.unit-include | Regexp of systemd units to include. Units must both match include and not match exclude to be included. | .+ |
| collector.systemd.unit-whitelist | DEPRECATED: Use --collector.systemd.unit-include |  |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_systemd_service_restart_total | Service unit count of Restart triggers | name |
| node_systemd_socket_accepted_connections_total | Total number of accepted socket connections | name |
| node_systemd_socket_current_connections | Current number of socket connections | name |
| node_systemd_socket_refused_connections_total | Total number of refused socket connections | name |
| node_systemd_system_running | Whether the system is operational (see 'systemctl is-system-running') | n/a |
| node_systemd_timer_last_trigger_seconds | Seconds since epoch of last trigger. | name |
| node_systemd_unit_start_time_seconds | Start time of the unit since unix epoch in seconds. | name |
| node_systemd_unit_state | Systemd unit | name, state, type |
| node_systemd_unit_tasks_current | Current number of tasks per Systemd unit | name |
| node_systemd_unit_tasks_max | Maximum number of tasks per Systemd unit | name |
| node_systemd_units | Summary of systemd unit states | state |
| node_systemd_version | Detected systemd version | version |
| node_systemd_virtualization_info | Detected virtualization technology | virtualization_type |
