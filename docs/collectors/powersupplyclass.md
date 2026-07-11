# powersupplyclass collector

The powersupplyclass collector exposes metrics about powersupplyclass.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.powersupply.ignored-supplies | Regexp of power supplies to ignore for powersupplyclass collector. | ^$ |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_subsystem_info | IOKit Power Source information for <power_supply>. | n/a |
