# hwmon

Exposes hardware monitoring statistics from `/sys/class/hwmon/`, similar to `lm-sensors`.

Status: enabled by default

## Platforms

- Linux

## Configuration

```
--collector.hwmon.chip-include    Regexp of chips to include (mutually exclusive to chip-exclude)
--collector.hwmon.chip-exclude    Regexp of chips to exclude (mutually exclusive to chip-include)
--collector.hwmon.sensor-include  Regexp of sensors to include (mutually exclusive to sensor-exclude)
--collector.hwmon.sensor-exclude  Regexp of sensors to exclude (mutually exclusive to sensor-include)
```

### Examples

Exclude a specific chip:
```
--collector.hwmon.chip-exclude="^platform_thinkpad_hwmon$"
```

Monitor only coretemp sensors:
```
--collector.hwmon.chip-include="^platform_coretemp.*"
```

Exclude specific sensor on specific chip (format: `chip;sensor`):
```
--collector.hwmon.sensor-exclude="platform_coretemp_0;temp3"
```

Monitor only temperature sensors:
```
--collector.hwmon.sensor-include=";temp[0-9]+"
```

## Data Sources

| Source | Description |
|--------|-------------|
| `/sys/class/hwmon/` | Hardware monitoring chips and sensors |

Documentation:
- https://www.kernel.org/doc/Documentation/hwmon/sysfs-interface
- `sensors(1)` manpage (lm-sensors)

## Metrics

All metrics have `chip` and `sensor` labels.

### Metadata

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `node_hwmon_chip_names` | gauge | `chip`, `chip_name` | Human-readable chip name annotation (always 1) |
| `node_hwmon_sensor_label` | gauge | `chip`, `sensor`, `label` | Sensor label annotation (always 1) |

### Temperature

| Metric | Type | Description |
|--------|------|-------------|
| `node_hwmon_temp_celsius` | gauge | Temperature reading in Celsius |
| `node_hwmon_temp_crit_celsius` | gauge | Critical temperature threshold |
| `node_hwmon_temp_crit_alarm_celsius` | gauge | Critical alarm temperature |
| `node_hwmon_temp_max_celsius` | gauge | Maximum temperature threshold |
| `node_hwmon_temp_min_celsius` | gauge | Minimum temperature threshold |

### Voltage

| Metric | Type | Description |
|--------|------|-------------|
| `node_hwmon_in_volts` | gauge | Voltage reading in volts |
| `node_hwmon_in_min_volts` | gauge | Minimum voltage threshold |
| `node_hwmon_in_max_volts` | gauge | Maximum voltage threshold |
| `node_hwmon_in_crit_volts` | gauge | Critical voltage threshold |
| `node_hwmon_cpu_volts` | gauge | CPU voltage in volts |

### Fan

| Metric | Type | Description |
|--------|------|-------------|
| `node_hwmon_fan_rpm` | gauge | Fan speed in RPM |
| `node_hwmon_fan_min_rpm` | gauge | Minimum fan speed |
| `node_hwmon_fan_max_rpm` | gauge | Maximum fan speed |
| `node_hwmon_fan_target_rpm` | gauge | Target fan speed |
| `node_hwmon_fan_alarm` | gauge | Fan alarm status |
| `node_hwmon_fan_fault` | gauge | Fan fault status |

### Power

| Metric | Type | Description |
|--------|------|-------------|
| `node_hwmon_power_watt` | gauge | Power usage in watts |
| `node_hwmon_power_max_watt` | gauge | Maximum power |
| `node_hwmon_power_crit_watt` | gauge | Critical power threshold |
| `node_hwmon_power_accuracy` | gauge | Power meter accuracy ratio |
| `node_hwmon_power_average_interval_seconds` | gauge | Power averaging interval |

### Current

| Metric | Type | Description |
|--------|------|-------------|
| `node_hwmon_curr_amps` | gauge | Current in amperes |
| `node_hwmon_curr_min_amps` | gauge | Minimum current threshold |
| `node_hwmon_curr_max_amps` | gauge | Maximum current threshold |
| `node_hwmon_curr_crit_amps` | gauge | Critical current threshold |

### Energy

| Metric | Type | Description |
|--------|------|-------------|
| `node_hwmon_energy_joule_total` | counter | Total energy consumed in joules |

### PWM

| Metric | Type | Description |
|--------|------|-------------|
| `node_hwmon_pwm` | gauge | PWM value (0-255) |
| `node_hwmon_pwm_enable` | gauge | PWM control mode |

### Other

| Metric | Type | Description |
|--------|------|-------------|
| `node_hwmon_humidity` | gauge | Humidity as ratio (multiply by 100 for percentage) |
| `node_hwmon_intrusion_alarm` | gauge | Chassis intrusion detection |
| `node_hwmon_freq_freq_mhz` | gauge | GPU frequency in MHz |
| `node_hwmon_beep_enabled` | gauge | Beep enabled status |
| `node_hwmon_voltage_regulator_version` | gauge | VRM version |
| `node_hwmon_update_interval_seconds` | gauge | Sensor update interval |

## Labels

| Label | Description |
|-------|-------------|
| `chip` | Chip identifier derived from device path or name (e.g., `platform_coretemp_0`, `pci0000_00_1f_3`) |
| `sensor` | Sensor identifier (e.g., `temp1`, `fan2`, `in0`) |
| `chip_name` | Human-readable chip name from sysfs (chip_names metric only) |
| `label` | Sensor label from sysfs if available (sensor_label metric only) |

## Notes

- Chip names are derived from device paths to ensure stability across reboots (hwmon numbering can change)
- Sensor filtering uses format `chip;sensor` to allow per-chip sensor exclusion
- Raw sysfs values are converted to standard units (millivolts -> volts, millidegrees -> degrees)
- Some drivers return EAGAIN; the collector handles this gracefully
- Use `sensors` command from lm-sensors to explore available sensors
