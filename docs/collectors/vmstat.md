# vmstat collector

The vmstat collector exposes metrics about vmstat.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.vmstat.fields | Regexp of fields to return for vmstat collector. | ^(oom_kill|pgpg|pswp|pg.*fault).* |

