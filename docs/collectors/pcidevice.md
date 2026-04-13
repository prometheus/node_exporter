# pcidevice collector

The pcidevice collector exposes metrics about pcidevice.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.pcidevice.idsfile | Path to pci.ids file to use for PCI device identification. |  |
| collector.pcidevice.names | Enable PCI device name resolution (requires pci.ids file). | false |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_pcidevice_current_link_transfers_per_second | Value of current link's transfers per second (T/s) | n/a |
| node_pcidevice_current_link_width | Value of current link's width (number of lanes) | n/a |
| node_pcidevice_d3cold_allowed | Whether the PCIe device supports D3cold power state (0/1). | n/a |
| node_pcidevice_info | Non-numeric data from /sys/bus/pci/devices/<location>, value is always 1. | n/a |
| node_pcidevice_max_link_transfers_per_second | Value of maximum link's transfers per second (T/s) | n/a |
| node_pcidevice_max_link_width | Value of maximum link's width (number of lanes) | n/a |
| node_pcidevice_numa_node | NUMA node number for the PCI device. -1 indicates unknown or not available. | n/a |
| node_pcidevice_power_state | PCIe device power state, one of: D0, D1, D2, D3hot, D3cold, unknown or error. | n/a |
| node_pcidevice_sriov_drivers_autoprobe | Whether SR-IOV drivers autoprobe is enabled for the device (0/1). | n/a |
| node_pcidevice_sriov_numvfs | Number of Virtual Functions (VFs) currently enabled for SR-IOV. | n/a |
| node_pcidevice_sriov_totalvfs | Total number of Virtual Functions (VFs) supported by the device. | n/a |
| node_pcidevice_sriov_vf_total_msix | Total number of MSI-X vectors for Virtual Functions. | n/a |
