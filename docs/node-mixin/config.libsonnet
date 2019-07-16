{
  _config+:: {
    // Selectors are inserted between {} in Prometheus queries.

    // Select the metrics coming from the node exporter.
    nodeExporterSelector: 'job="node"',

    // Select the fstype for filesystem-related queries.
    // TODO: What is a good default selector here?
    fsSelector: 'fstype=~"ext.|xfs|jfs|btrfs|vfat|ntfs"',

    // Select the device for disk-related queries.
    diskDeviceSelector: 'device=~"(sd|xvd).+"',

    grafana_prefix: '',
  },
}
