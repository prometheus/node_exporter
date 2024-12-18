{
  _config+:: {
    // Selectors are inserted between {} in Prometheus queries.

    // Select the metrics coming from the node exporter. Note that all
    // the selected metrics are shown stacked on top of each other in
    // the 'USE Method / Cluster' dashboard. Consider disabling that
    // dashboard if mixing up all those metrics in the same dashboard
    // doesn't make sense (e.g. because they are coming from different
    // clusters).
    nodeExporterSelector: 'job="node"',

    // Select the fstype for filesystem-related queries. If left
    // empty, all filesystems are selected. If you have unusual
    // filesystem you don't want to include in dashboards and
    // alerting, you can exclude them here, e.g. 'fstype!="tmpfs"'.
    fsSelector: 'fstype!=""',

    // Select the mountpoint for filesystem-related queries. If left
    // empty, all mountpoints are selected. For example if you have a
    // special purpose tmpfs instance that has a fixed size and will
    // always be 100% full, but you still want alerts and dashboards for
    // other tmpfs instances, you can exclude those by mountpoint prefix
    // like so: 'mountpoint!~"/var/lib/foo.*"'.
    fsMountpointSelector: 'mountpoint!=""',

    // Select the device for disk-related queries. If left empty, all
    // devices are selected. If you have unusual devices you don't
    // want to include in dashboards and alerting, you can exclude
    // them here, e.g. 'device!="tmpfs"'.
    diskDeviceSelector: 'device!=""',

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

    // CPU utilization (%) on which to trigger the
    // 'NodeCPUHighUsage' alert.
    cpuHighUsageThreshold: 90,
    // Load average 1m (per core) on which to trigger the
    // 'NodeSystemSaturation' alert.
    systemSaturationPerCoreThreshold: 2,

    // Some of the alerts use predict_linear() to fire alerts ahead of time to
    // prevent unrecoverable situations (eg. no more disk space). However, the
    // node may have automatic processes (cronjobs) in place to prevent that
    // within a certain time window, this may not align with the default time
    // window of these alerts. This can cause these alerts to start flapping.
    // By reducing the time window, the system gets more time to
    // resolve this before problems occur.
    nodeWarningWindowHours: '24',
    nodeCriticalWindowHours: '4',

    // Available disk space (%) thresholds on which to trigger the
    // 'NodeFilesystemSpaceFillingUp' alerts. These alerts fire if the disk
    // usage grows in a way that it is predicted to run out in 4h or 1d
    // and if the provided thresholds have been reached right now.
    // In some cases you'll want to adjust these, e.g., by default, Kubernetes
    // runs the image garbage collection when the disk usage reaches 85%
    // of its available space. In that case, you'll want to reduce the
    // critical threshold below to something like 14 or 15, otherwise
    // the alert could fire under normal node usage.
    // Additionally, the prediction window for the alert can be configured
    // to account for environments where disk usage can fluctuate within
    // a short time frame. By extending the prediction window, you can
    // reduce false positives caused by temporary spikes, providing a
    // more accurate prediction of disk space issues.
    fsSpaceFillingUpWarningThreshold: 40,
    fsSpaceFillingUpCriticalThreshold: 20,
    fsSpaceFillingUpPredictionWindow: '6h',

    // Available disk space (%) thresholds on which to trigger the
    // 'NodeFilesystemAlmostOutOfSpace' alerts.
    fsSpaceAvailableWarningThreshold: 5,
    fsSpaceAvailableCriticalThreshold: 3,

    // Memory utilization (%) level on which to trigger the
    // 'NodeMemoryHighUtilization' alert.
    memoryHighUtilizationThreshold: 90,

    // Threshold for the rate of memory major page faults to trigger
    // 'NodeMemoryMajorPagesFaults' alert.
    memoryMajorPagesFaultsThreshold: 500,

    // Disk IO queue level above which to trigger
    // 'NodeDiskIOSaturation' alert.
    diskIOSaturationThreshold: 10,

    rateInterval: '5m',
    // Opt-in for multi-cluster support.
    showMultiCluster: false,
    clusterLabel: 'cluster',

    dashboardNamePrefix: 'Node Exporter / ',
    dashboardTags: ['node-exporter-mixin'],
  },
}
