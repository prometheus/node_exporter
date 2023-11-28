{

  //  any modular observability library should inlcude as inputs:
  // 'dashboardNamePrefix' - Use as prefix for all Dashboards and (optional) rule groups
  // 'filteringSelector' - Static selector to apply to ALL dashboard variables of type query, panel queries, alerts and recording rules.
  // 'groupLabels' - one or more labels that can be used to identify 'group' of instances. In simple cases, can be 'job' or 'cluster'.
  // 'instanceLabels' - one or more labels that can be used to identify single entity of instances. In simple cases, can be 'instance' or 'pod'.
  // 'uid' - UID to prefix all dashboards original uids

  filteringSelector: 'job="integrations/macos-node"',
  groupLabels: ['job'],
  instanceLabels: ['instance'],
  dashboardNamePrefix: 'MacOS / ',
  uid: 'darwin',

  dashboardTags: [self.uid],

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
  dashboardPeriod: 'now-1h',
  dashboardTimezone: 'default',
  dashboardRefresh: '1m',

  // Alerts to keep from node-observ-lib:
  alertsMacKeep: [
    'NodeFilesystemAlmostOutOfSpace',
    'NodeNetworkReceiveErrs',
    'NodeNetworkTransmitErrs',
    'NodeTextFileCollectorScrapeError',
    'NodeFilesystemFilesFillingUp',
    'NodeFilesystemAlmostOutOfFiles',
  ],
  // logs lib related
  enableLokiLogs: true,
  extraLogLabels: ['filename', 'sender'],
  logsVolumeGroupBy: 'sender',
  showLogsVolume: true,
  logsFilteringSelector: self.filteringSelector,
  logsExtraFilters: '',


}
