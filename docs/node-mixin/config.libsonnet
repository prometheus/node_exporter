{
  _config+:: {
    // Selectors are inserted between {} in Prometheus queries.

    // Select the metrics coming from the node exporter.
    nodeExporterSelector: 'job="node"',

    // Select the fstype for filesystem-related queries. If left
    // empty, all filesystems are selected. If you have unusual
    // filesystem you don't want to include in dashboards and
    // alerting, you can exclude them here, e.g. 'fstype!="tmpfs"'.
    fsSelector: '',

    // Select the device for disk-related queries. If left empty, all
    // devices are selected. If you have unusual devices you don't
    // want to include in dashboards and alerting, you can exclude
    // them here, e.g. 'device!="tmpfs"'.
    diskDeviceSelector: '',

    grafana_prefix: '',
  },
}
