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

    // Some of the alerts are meant to fire if a critical failure of a
    // node is imminent (e.g. the disk is about to run full). In a
    // true “cloud native” setup, failures of a single node should be
    // tolerated. Hence, even imminent failure of a single node is no
    // reason to create a paging alert. However, in practice there are
    // still many situations where operators like to get paged in time
    // before a node runs out of disk space. nodeCriticalSeverity can
    // be set to the desired severity for this kind of alerts. This
    // can even be templated to depend on labels of the node, e.g. you
    // could make this critical for traditional database masters but
    // just a warning for K8s nodes.
    nodeCriticalSeverity: 'critical',

    grafana_prefix: '',
  },
}
