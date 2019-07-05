{
  _config+:: {
    // Selectors are inserted between {} in Prometheus queries.
    nodeExporterSelector: 'job="node-exporter"',

    // Mainly extracted because they are repetitive, but also useful to customize.
    fsSelectors: 'fstype=~"ext.|xfs",mountpoint!="/var/lib/docker/aufs"',

    grafana_prefix: '',
  },
}
