local grafana = import 'github.com/grafana/grafonnet/gen/grafonnet-latest/main.libsonnet';
local dashboard = grafana.dashboard;
local variable = dashboard.variable;
local row = grafana.panel.row;
local prometheus = grafana.query.prometheus;

local timeSeriesPanel = grafana.panel.timeSeries;
local tsOptions = timeSeriesPanel.options;
local tsStandardOptions = timeSeriesPanel.standardOptions;
local tsQueryOptions = timeSeriesPanel.queryOptions;
local tsCustom = timeSeriesPanel.fieldConfig.defaults.custom;
local tsLegend = tsOptions.legend;

local c = import '../config.libsonnet';

local datasource = variable.datasource.new(
  'datasource', 'prometheus'
);

local tsCommonPanelOptions =
  variable.query.withDatasourceFromVariable(datasource)
  + tsCustom.stacking.withMode('normal')
  + tsCustom.withFillOpacity(100)
  + tsCustom.withShowPoints('never')
  + tsLegend.withShowLegend(false)
  + tsOptions.tooltip.withMode('multi')
  + tsOptions.tooltip.withSort('desc');

local CPUUtilisation =
  timeSeriesPanel.new(
    'CPU Utilisation',
  )
  + tsCommonPanelOptions
  + tsStandardOptions.withUnit('percentunit');

local CPUSaturation =
  // TODO: Is this a useful panel? At least there should be some explanation how load
  // average relates to the "CPU saturation" in the title.
  timeSeriesPanel.new(
    'CPU Saturation (Load1 per CPU)',
  )
  + tsCommonPanelOptions
  + tsStandardOptions.withUnit('percentunit');

local memoryUtilisation =
  timeSeriesPanel.new(
    'Memory Utilisation',
  )
  + tsCommonPanelOptions
  + tsStandardOptions.withUnit('percentunit');

local memorySaturation =
  timeSeriesPanel.new(
    'Memory Saturation (Major Page Faults)',
  )
  + tsCommonPanelOptions
  + tsStandardOptions.withUnit('rds');

local networkOverrides = tsStandardOptions.withOverrides(
  [
    tsStandardOptions.override.byRegexp.new('/Transmit/')
    + tsStandardOptions.override.byRegexp.withPropertiesFromOptions(
      tsCustom.withTransform('negative-Y')
    ),
  ]
);

local networkUtilisation =
  timeSeriesPanel.new(
    'Network Utilisation (Bytes Receive/Transmit)',
  )
  + tsCommonPanelOptions
  + tsStandardOptions.withUnit('Bps')
  + networkOverrides;

local networkSaturation =
  timeSeriesPanel.new(
    'Network Saturation (Drops Receive/Transmit)',
  )
  + tsCommonPanelOptions
  + tsStandardOptions.withUnit('Bps')
  + networkOverrides;

local diskIOUtilisation =
  timeSeriesPanel.new(
    'Disk IO Utilisation',
  )
  + tsCommonPanelOptions
  + tsStandardOptions.withUnit('percentunit');

local diskIOSaturation =
  timeSeriesPanel.new(
    'Disk IO Saturation',
  )
  + tsCommonPanelOptions
  + tsStandardOptions.withUnit('percentunit');

local diskSpaceUtilisation =
  timeSeriesPanel.new(
    'Disk Space Utilisation',
  )
  + tsCommonPanelOptions
  + tsStandardOptions.withUnit('percentunit');

{
  _clusterVariable::
    variable.query.new('cluster')
    + variable.query.withDatasourceFromVariable(datasource)
    + variable.query.queryTypes.withLabelValues(
      $._config.clusterLabel,
      'node_time_seconds',
    )
    + (if $._config.showMultiCluster then variable.query.generalOptions.showOnDashboard.withLabelAndValue() else variable.query.generalOptions.showOnDashboard.withNothing())
    + variable.query.refresh.onTime()
    + variable.query.selectionOptions.withIncludeAll(false)
    + variable.query.withSort(asc=true),

  grafanaDashboards+:: {
                         'node-rsrc-use.json':
                           dashboard.new(
                             '%sUSE Method / Node' % $._config.dashboardNamePrefix,
                           )
                           + dashboard.time.withFrom('now-1h')
                           + dashboard.withTags($._config.dashboardTags)
                           + dashboard.withTimezone('utc')
                           + dashboard.withRefresh('30s')
                           + dashboard.graphTooltip.withSharedCrosshair()
                           + dashboard.withUid(std.md5('node-rsrc-use.json'))
                           + dashboard.withVariables([
                             datasource,
                             $._clusterVariable,
                             variable.query.new('instance')
                             + variable.query.withDatasourceFromVariable(datasource)
                             + variable.query.queryTypes.withLabelValues(
                               'instance',
                               'node_exporter_build_info{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}' % $._config,
                             )
                             + variable.query.refresh.onTime()
                             + variable.query.withSort(asc=true),
                           ])
                           + dashboard.withPanels(
                             grafana.util.grid.makeGrid([
                               row.new('CPU')
                               + row.withPanels([
                                 CPUUtilisation + tsQueryOptions.withTargets([prometheus.new('$datasource', 'instance:node_cpu_utilisation:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config) + prometheus.withLegendFormat('Utilisation')]),
                                 CPUSaturation + tsQueryOptions.withTargets([prometheus.new('$datasource', 'instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config) + prometheus.withLegendFormat('Saturation')]),
                               ]),
                               row.new('Memory')
                               + row.withPanels([
                                 memoryUtilisation + tsQueryOptions.withTargets([prometheus.new('$datasource', 'instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config) + prometheus.withLegendFormat('Utilisation')]),
                                 memorySaturation + tsQueryOptions.withTargets([prometheus.new('$datasource', 'instance:node_vmstat_pgmajfault:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config) + prometheus.withLegendFormat('Major page Faults')]),
                               ]),
                               row.new('Network')
                               + row.withPanels([
                                 networkUtilisation + tsQueryOptions.withTargets([
                                   prometheus.new('$datasource', 'instance:node_network_receive_bytes_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config) + prometheus.withLegendFormat('Receive'),
                                   prometheus.new('$datasource', 'instance:node_network_transmit_bytes_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config) + prometheus.withLegendFormat('Transmit'),
                                 ]),
                                 networkSaturation + tsQueryOptions.withTargets([
                                   prometheus.new('$datasource', 'instance:node_network_receive_drop_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config) + prometheus.withLegendFormat('Receive'),
                                   prometheus.new('$datasource', 'instance:node_network_transmit_drop_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config) + prometheus.withLegendFormat('Transmit'),
                                 ]),
                               ]),
                               row.new('Disk IO')
                               + row.withPanels([
                                 diskIOUtilisation + tsQueryOptions.withTargets([prometheus.new('$datasource', 'instance_device:node_disk_io_time_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config) + prometheus.withLegendFormat('{{device}}')]),
                                 diskIOSaturation + tsQueryOptions.withTargets([prometheus.new('$datasource', 'instance_device:node_disk_io_time_weighted_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config) + prometheus.withLegendFormat('{{device}}')]),
                               ]),
                             ], panelWidth=12, panelHeight=7)
                             + grafana.util.grid.makeGrid([
                               row.new('Disk Space')
                               + row.withPanels([
                                 diskSpaceUtilisation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       sort_desc(1 -
                                         (
                                           max without (mountpoint, fstype) (node_filesystem_avail_bytes{%(nodeExporterSelector)s, fstype!="", instance="$instance", %(clusterLabel)s="$cluster"})
                                           /
                                           max without (mountpoint, fstype) (node_filesystem_size_bytes{%(nodeExporterSelector)s, fstype!="", instance="$instance", %(clusterLabel)s="$cluster"})
                                         ) != 0
                                       )
                                     ||| % $._config
                                   ) + prometheus.withLegendFormat('{{device}}'),
                                 ]),
                               ]),
                             ], panelWidth=24, panelHeight=7, startY=34),
                           ),
                         'node-cluster-rsrc-use.json':
                           dashboard.new(
                             '%sUSE Method / Cluster' % $._config.dashboardNamePrefix,
                           )
                           + dashboard.time.withFrom('now-1h')
                           + dashboard.withTags($._config.dashboardTags)
                           + dashboard.withTimezone('utc')
                           + dashboard.withRefresh('30s')
                           + dashboard.graphTooltip.withSharedCrosshair()
                           + dashboard.withUid(std.md5('node-cluster-rsrc-use.json'))
                           + dashboard.withVariables([
                             datasource,
                             $._clusterVariable,
                             variable.query.withDatasourceFromVariable(datasource)
                             + variable.query.refresh.onTime()
                             + variable.query.withSort(asc=true),
                           ])
                           + dashboard.withPanels(
                             grafana.util.grid.makeGrid([
                               row.new('CPU')
                               + row.withPanels([
                                 CPUUtilisation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       ((
                                         instance:node_cpu_utilisation:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}
                                         *
                                         instance:node_num_cpu:sum{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}
                                       ) != 0 )
                                       / scalar(sum(instance:node_num_cpu:sum{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}))
                                     ||| % $._config
                                   ) + prometheus.withLegendFormat('{{ instance }}'),
                                 ]),
                                 CPUSaturation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       (
                                         instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}
                                         / scalar(count(instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}))
                                       )  != 0
                                     ||| % $._config
                                   ) + prometheus.withLegendFormat('{{ instance }}'),
                                 ]),
                               ]),
                               row.new('Memory')
                               + row.withPanels([
                                 memoryUtilisation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       (
                                         instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}
                                         / scalar(count(instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}))
                                       ) != 0
                                     ||| % $._config
                                   ) + prometheus.withLegendFormat('{{ instance }}'),
                                 ]),
                                 memorySaturation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     'instance:node_vmstat_pgmajfault:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}' % $._config
                                   ) + prometheus.withLegendFormat('{{ instance }}'),
                                 ]),
                               ]),
                               row.new('Network')
                               + row.withPanels([
                                 networkUtilisation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     'instance:node_network_receive_bytes_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"} != 0' % $._config
                                   ) + prometheus.withLegendFormat('{{ instance }} Receive'),
                                   prometheus.new(
                                     '$datasource',
                                     'instance:node_network_transmit_bytes_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"} != 0' % $._config
                                   ) + prometheus.withLegendFormat('{{ instance }} Transmit'),
                                 ]),
                                 networkSaturation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     'instance:node_network_receive_drop_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"} != 0' % $._config
                                   ) + prometheus.withLegendFormat('{{ instance }} Receive'),
                                   prometheus.new(
                                     '$datasource',
                                     'instance:node_network_transmit_drop_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"} != 0' % $._config
                                   ) + prometheus.withLegendFormat('{{ instance }} Transmit'),
                                 ]),
                               ]),
                               row.new('Disk IO')
                               + row.withPanels([
                                 diskIOUtilisation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       instance_device:node_disk_io_time_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}
                                       / scalar(count(instance_device:node_disk_io_time_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}))
                                     ||| % $._config
                                   ) + prometheus.withLegendFormat('{{ instance }} {{device}}'),
                                 ]),
                                 diskIOSaturation + tsQueryOptions.withTargets([prometheus.new(
                                   '$datasource',
                                   |||
                                     instance_device:node_disk_io_time_weighted_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}
                                     / scalar(count(instance_device:node_disk_io_time_weighted_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}))
                                   ||| % $._config
                                 ) + prometheus.withLegendFormat('{{ instance }} {{device}}')]),
                               ]),
                             ], panelWidth=12, panelHeight=7)
                             + grafana.util.grid.makeGrid([
                               row.new('Disk Space')
                               + row.withPanels([
                                 diskSpaceUtilisation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       sum without (device) (
                                         max without (fstype, mountpoint) ((
                                           node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s, %(fsMountpointSelector)s, %(clusterLabel)s="$cluster"}
                                           -
                                           node_filesystem_avail_bytes{%(nodeExporterSelector)s, %(fsSelector)s, %(fsMountpointSelector)s, %(clusterLabel)s="$cluster"}
                                         ) != 0)
                                       )
                                       / scalar(sum(max without (fstype, mountpoint) (node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s, %(fsMountpointSelector)s, %(clusterLabel)s="$cluster"})))
                                     ||| % $._config
                                   ) + prometheus.withLegendFormat('{{ instance }}'),
                                 ]),
                               ]),
                             ], panelWidth=24, panelHeight=7, startY=34),
                           ),
                       } +
                       if $._config.showMultiCluster then {
                         'node-multicluster-rsrc-use.json':
                           dashboard.new(
                             '%sUSE Method / Multi-cluster' % $._config.dashboardNamePrefix,
                           )
                           + dashboard.time.withFrom('now-1h')
                           + dashboard.withTags($._config.dashboardTags)
                           + dashboard.withTimezone('utc')
                           + dashboard.withRefresh('30s')
                           + dashboard.graphTooltip.withSharedCrosshair()
                           + dashboard.withUid(std.md5('node-multicluster-rsrc-use.json'))
                           + dashboard.withVariables([
                             datasource,
                             variable.query.withDatasourceFromVariable(datasource)
                             + variable.query.refresh.onTime()
                             + variable.query.withSort(asc=true),
                           ])
                           + dashboard.withPanels(
                             grafana.util.grid.makeGrid([
                               row.new('CPU')
                               + row.withPanels([
                                 CPUUtilisation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       sum(
                                         ((
                                           instance:node_cpu_utilisation:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                           *
                                           instance:node_num_cpu:sum{%(nodeExporterSelector)s}
                                         ) != 0)
                                         / scalar(sum(instance:node_num_cpu:sum{%(nodeExporterSelector)s}))
                                       ) by (%(clusterLabel)s)
                                     ||| % $._config
                                   ) + prometheus.withLegendFormat('{{%(clusterLabel)s}}'),
                                 ]),
                                 CPUSaturation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       sum((
                                           instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s}
                                           / scalar(count(instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s}))
                                       ) != 0) by (%(clusterLabel)s)
                                     ||| % $._config
                                   ) + prometheus.withLegendFormat('{{%(clusterLabel)s}}'),
                                 ]),
                               ]),
                               row.new('Memory')
                               + row.withPanels([
                                 memoryUtilisation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       sum((
                                           instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s}
                                           / scalar(count(instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s}))
                                       ) != 0) by (%(clusterLabel)s)
                                     ||| % $._config
                                   ) + prometheus.withLegendFormat('{{%(clusterLabel)s}}'),
                                 ]),
                                 memorySaturation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       sum((
                                           instance:node_vmstat_pgmajfault:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                       ) != 0) by (%(clusterLabel)s)
                                     |||
                                     % $._config
                                   ) + prometheus.withLegendFormat('{{%(clusterLabel)s}}'),
                                 ]),
                               ]),
                               row.new('Network')
                               + row.withPanels([
                                 networkUtilisation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       sum((
                                           instance:node_network_receive_bytes_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                       ) != 0) by (%(clusterLabel)s)
                                     ||| % $._config
                                   ) + prometheus.withLegendFormat('{{%(clusterLabel)s}} Receive'),
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       sum((
                                           instance:node_network_transmit_bytes_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                       ) != 0) by (%(clusterLabel)s)
                                     ||| % $._config
                                   ) + prometheus.withLegendFormat('{{%(clusterLabel)s}} Transmit'),
                                 ]),
                                 networkSaturation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       sum((
                                           instance:node_network_receive_drop_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                       ) != 0) by (%(clusterLabel)s)
                                     ||| % $._config
                                   ) + prometheus.withLegendFormat('{{%(clusterLabel)s}} Receive'),
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       sum((
                                           instance:node_network_transmit_drop_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                       ) != 0) by (%(clusterLabel)s)
                                     ||| % $._config
                                   ) + prometheus.withLegendFormat('{{%(clusterLabel)s}} Transmit'),
                                 ]),
                               ]),
                               row.new('Disk IO')
                               + row.withPanels([
                                 diskIOUtilisation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       sum((
                                           instance_device:node_disk_io_time_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                           / scalar(count(instance_device:node_disk_io_time_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s}))
                                       ) != 0) by (%(clusterLabel)s, device)
                                     ||| % $._config
                                   ) + prometheus.withLegendFormat('{{%(clusterLabel)s}} {{device}}'),
                                 ]),
                                 diskIOSaturation + tsQueryOptions.withTargets([prometheus.new(
                                   '$datasource',
                                   |||
                                     sum((
                                       instance_device:node_disk_io_time_weighted_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                       / scalar(count(instance_device:node_disk_io_time_weighted_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s}))
                                     ) != 0) by (%(clusterLabel)s, device)
                                   ||| % $._config
                                 ) + prometheus.withLegendFormat('{{%(clusterLabel)s}} {{device}}')]),
                               ]),

                             ], panelWidth=12, panelHeight=7)
                             + grafana.util.grid.makeGrid([
                               row.new('Disk Space')
                               + row.withPanels([
                                 diskSpaceUtilisation + tsQueryOptions.withTargets([
                                   prometheus.new(
                                     '$datasource',
                                     |||
                                       sum (
                                         sum without (device) (
                                           max without (fstype, mountpoint, instance, pod) ((
                                             node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s, %(fsMountpointSelector)s} - node_filesystem_avail_bytes{%(nodeExporterSelector)s, %(fsSelector)s, %(fsMountpointSelector)s}
                                           ) != 0)
                                         )
                                         / scalar(sum(max without (fstype, mountpoint) (node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s, %(fsMountpointSelector)s})))
                                       ) by (%(clusterLabel)s)
                                     ||| % $._config
                                   ) + prometheus.withLegendFormat('{{%(clusterLabel)s}}'),
                                 ]),
                               ]),
                             ], panelWidth=24, panelHeight=7, startY=34),
                           ),
                       } else {},
}
