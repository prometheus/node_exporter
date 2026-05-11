# collector_name

Brief description of what this collector exposes.

Status: enabled|disabled by default

## Platforms

- Linux
- Darwin
- FreeBSD
- ...

## Configuration

```
--collector.name.flag-name    Description (default: value)
--collector.name.other-flag   Description (default: value)
```

Omit this section if the collector has no flags.

## Data Sources

| Source | Description |
|--------|-------------|
| `/proc/example` | Brief description |
| `/sys/class/example` | Brief description |
| `syscall(2)` | Brief description |

## Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `node_example_total` | counter | `label1`, `label2` | Description |
| `node_example_bytes` | gauge | | Description |
| `node_example_info` | gauge | `key`, `value` | Info metric, always 1 |

For collectors with dynamic metrics (e.g., meminfo), use:

Metrics are derived from `/proc/meminfo`. Each field `FieldName` becomes `node_memory_fieldname_bytes`.

## Labels

| Label | Description |
|-------|-------------|
| `device` | Device name |
| `mountpoint` | Mount path |

Omit this section if metrics have no labels or labels are self-explanatory.

## Notes

- Special behaviors, caveats, kernel version requirements
- Known issues or limitations
- Related collectors

Omit this section if not applicable.
