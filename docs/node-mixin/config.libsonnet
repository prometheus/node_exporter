{
  _config+:: {
    // Selectors are inserted between {} in Prometheus queries.

    // Select the metrics coming from the node exporter.
    nodeExporterSelector: 'job="node-exporter"',

    // Select the fstype for filesystem-related queries.
    fsSelector: 'fstype=~"ext.|xfs",mountpoint!="/var/lib/docker/aufs"',

    // Select the device for disk-related queries.
    diskDeviceSelector: 'device=~"(sd|xvd).+"',

    grafana_prefix: '',
  },
}
