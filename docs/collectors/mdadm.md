# mdadm collector

The mdadm collector exposes metrics about mdadm.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_md_blocks | Total number of blocks on device. | device |
| node_md_blocks_synced | Number of blocks synced on device. | device |
| node_md_degraded | Number of degraded disks on device. | device |
| node_md_disks | Number of active/failed/spare disks of device. | device, state |
| node_md_disks_required | Total number of disks of device. | device |
| node_md_raid_disks | Number of raid disks on device. | device |
| node_md_state | Indicates the state of md-device. | device |
