# interrupts collector

The interrupts collector exposes metrics about interrupts.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.interrupts.include-zeros | Include interrupts that have a zero value | true |
| collector.interrupts.name-exclude | Regexp of interrupts name to exclude (mutually exclusive to --collector.interrupts.name-include). |  |
| collector.interrupts.name-include | Regexp of interrupts name to include (mutually exclusive to --collector.interrupts.name-exclude). |  |

