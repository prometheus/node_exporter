# hwmon collector

The hwmon collector exposes metrics about hwmon.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.hwmon.chip-exclude | Regexp of hwmon chip to exclude (mutually exclusive to device-include). |  |
| collector.hwmon.chip-include | Regexp of hwmon chip to include (mutually exclusive to device-exclude). |  |
| collector.hwmon.sensor-exclude | Regexp of hwmon sensor to exclude (mutually exclusive to sensor-include). |  |
| collector.hwmon.sensor-include | Regexp of hwmon sensor to include (mutually exclusive to sensor-exclude). |  |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_hwmon_chip_names | Annotation metric for human-readable chip names | n/a |
| node_hwmon_sensor_label | Label for given chip and sensor | chip, sensor, label |
