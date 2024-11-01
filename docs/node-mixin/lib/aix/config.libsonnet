{
  // Rest of the config is imported from linux
  filteringSelector: 'job="aix"',
  dashboardNamePrefix: 'MacOS / ',
  //uid prefix
  uid: 'aix',

  dashboardTags: ['aix-mixin'],


  // Alerts to keep from node-observ-lib:
  alertsKeep: [
    'NodeFilesystemAlmostOutOfSpace',
    'NodeNetworkReceiveErrs',
    'NodeNetworkTransmitErrs',
    'NodeTextFileCollectorScrapeError',
    'NodeFilesystemFilesFillingUp',
    'NodeFilesystemAlmostOutOfFiles',
    'NodeCPUHighUsage',
    'NodeSystemSaturation',
    'NodeMemoryHighUtilization',
    'NodeDiskIOSaturation',
  ],
  // logs lib related
  enableLokiLogs: false,

}
