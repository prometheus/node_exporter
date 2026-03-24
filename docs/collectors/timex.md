# timex collector

The timex collector exposes metrics about timex.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_timex_estimated_error_seconds | Estimated error in seconds. | n/a |
| node_timex_frequency_adjustment_ratio | Local clock frequency adjustment. | n/a |
| node_timex_loop_time_constant | Phase-locked loop time constant. | n/a |
| node_timex_maxerror_seconds | Maximum error in seconds. | n/a |
| node_timex_offset_seconds | Time offset in between local system and reference clock. | n/a |
| node_timex_pps_calibration_total | Pulse per second count of calibration intervals. | n/a |
| node_timex_pps_error_total | Pulse per second count of calibration errors. | n/a |
| node_timex_pps_frequency_hertz | Pulse per second frequency. | n/a |
| node_timex_pps_jitter_seconds | Pulse per second jitter. | n/a |
| node_timex_pps_jitter_total | Pulse per second count of jitter limit exceeded events. | n/a |
| node_timex_pps_shift_seconds | Pulse per second interval duration. | n/a |
| node_timex_pps_stability_exceeded_total | Pulse per second count of stability limit exceeded events. | n/a |
| node_timex_pps_stability_hertz | Pulse per second stability, average of recent frequency changes. | n/a |
| node_timex_status | Value of the status array bits. | n/a |
| node_timex_sync_status | Is clock synchronized to a reliable server (1 = yes, 0 = no). | n/a |
| node_timex_tai_offset_seconds | International Atomic Time (TAI) offset. | n/a |
| node_timex_tick_seconds | Seconds between clock ticks. | n/a |
