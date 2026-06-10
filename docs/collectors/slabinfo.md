# slabinfo collector

The slabinfo collector exposes metrics about slabinfo.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.slabinfo.slabs-exclude | Regexp of slabs to exclude in slabinfo collector. |  |
| collector.slabinfo.slabs-include | Regexp of slabs to include in slabinfo collector. | .* |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_subsystem_active_objects | The number of objects that are currently active (i.e., in use). | n/a |
| node_subsystem_object_size_bytes | The size of objects in this slab, in bytes. | n/a |
| node_subsystem_objects | The total number of allocated objects (i.e., objects that are both in use and not in use). | n/a |
| node_subsystem_objects_per_slab | The number of objects stored in each slab. | n/a |
| node_subsystem_pages_per_slab | The number of pages allocated for each slab. | n/a |
