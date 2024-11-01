{
  // Rest of the config is imported from linux
  filteringSelector: 'job="integrations/macos-node"',
  dashboardNamePrefix: 'MacOS / ',
  uid: 'darwin',

  dashboardTags: [self.uid],


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
  enableLokiLogs: false,
  extraLogLabels: ['filename', 'sender'],
  logsVolumeGroupBy: 'sender',
  showLogsVolume: true,
  logsFilteringSelector: self.filteringSelector,
  logsExtraFilters: '',


}
