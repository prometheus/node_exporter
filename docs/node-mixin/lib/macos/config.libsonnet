{
  // Rest of the config is imported from linux
  filteringSelector: 'job=macos"',
  dashboardNamePrefix: 'MacOS / ',
  //uid prefix
  uid: 'darwin',

  dashboardTags: ['macos-mixin'],


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
