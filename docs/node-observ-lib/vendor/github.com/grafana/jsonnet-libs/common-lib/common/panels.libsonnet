local g = import './g.libsonnet';

{
  generic: {
    stat: import './panels/generic/stat/main.libsonnet',
    timeSeries: import './panels/generic/timeSeries/main.libsonnet',
    table: import './panels/generic/table/main.libsonnet',
    statusHistory: import './panels/generic/statusHistory/main.libsonnet',
  },
  network: {
    timeSeries: import './panels/network/timeSeries/main.libsonnet',
    statusHistory: import './panels/network/statusHistory/main.libsonnet',
  },
  system: {
    stat: import './panels/system/stat/main.libsonnet',
    table: import './panels/system/table/main.libsonnet',
    statusHistory: import './panels/system/statusHistory/main.libsonnet',
    timeSeries: import './panels/system/timeSeries/main.libsonnet',
  },
  cpu: {
    stat: import './panels/cpu/stat/main.libsonnet',
    timeSeries: import './panels/cpu/timeSeries/main.libsonnet',
  },
  memory: {
    stat: import './panels/memory/stat/main.libsonnet',
    timeSeries: import './panels/memory/timeSeries/main.libsonnet',
  },
  disk: {
    timeSeries: import './panels/disk/timeSeries/main.libsonnet',
    table: import './panels/disk/table/main.libsonnet',
    stat: import './panels/disk/stat/main.libsonnet',
  },
}
