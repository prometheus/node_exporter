# powersupplyclass collector

The powersupplyclass collector exposes metrics about power supplies (batteries, AC adapters, etc.) from `/sys/class/power_supply` on Linux and `IOKit` on Darwin.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.powersupply.ignored-supplies | Regexp of power supplies to ignore for powersupplyclass collector. | ^$ |
| collector.powersupply.use-new-names | Use standardized metric names and SI units for powersupplyclass collector. | false |

## Metrics

When `--collector.powersupply.use-new-names` is enabled, the collector exposes both the legacy names and the new, standardized names.

### New Standardized Metrics (Pluralized SI units)

| Metric | Unit | Description |
| --- | --- | --- |
| node_power_supply_current_amperes | Amperes | Current being consumed or provided. |
| node_power_supply_voltage_volts | Volts | Current voltage of the power supply. |
| node_power_supply_power_watts | Watts | Power consumption (if available). |
| node_power_supply_energy_joules | Joules | Energy remaining (converted from Watt-hours). |
| node_power_supply_charge_coulombs | Coulombs | Charge remaining (converted from Ampere-hours). |
| node_power_supply_capacity_ratio | Ratio (0.0-1.0) | Remaining capacity as a fraction of full capacity. |
| node_power_supply_temp_celsius | Celsius | Temperature of the power supply. |

### Info Metric

| Metric | Description | Labels |
| --- | --- | --- |
| node_power_supply_info | General information about the power supply (manufacturer, health, status, etc.) | power_supply, model_name, manufacturer, serial_number, status, etc. |
