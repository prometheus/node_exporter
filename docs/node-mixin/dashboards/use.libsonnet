local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;

local c = import '../config.libsonnet';

local datasourceTemplate = {
  current: {
    text: 'default',
    value: 'default',
  },
  hide: 0,
  label: 'Data Source',
  name: 'datasource',
  options: [],
  query: 'prometheus',
  refresh: 1,
  regex: '',
  type: 'datasource',
};

local CPUUtilisation =
  graphPanel.new(
    'CPU Utilisation',
    datasource='$datasource',
    span=6,
    format='percentunit',
    stack=true,
    fill=10,
    legend_show=false,
  ) { tooltip+: { sort: 2 } };

local CPUSaturation =
  // TODO: Is this a useful panel? At least there should be some explanation how load
  // average relates to the "CPU saturation" in the title.
  graphPanel.new(
    'CPU Saturation (Load1 per CPU)',
    datasource='$datasource',
    span=6,
    format='percentunit',
    stack=true,
    fill=10,
    legend_show=false,
  ) { tooltip+: { sort: 2 } };

local memoryUtilisation =
  graphPanel.new(
    'Memory Utilisation',
    datasource='$datasource',
    span=6,
    format='percentunit',
    stack=true,
    fill=10,
    legend_show=false,
  ) { tooltip+: { sort: 2 } };

local memorySaturation =
  graphPanel.new(
    'Memory Saturation (Major Page Faults)',
    datasource='$datasource',
    span=6,
    format='rds',
    stack=true,
    fill=10,
    legend_show=false,
  ) { tooltip+: { sort: 2 } };

local networkUtilisation =
  graphPanel.new(
    'Network Utilisation (Bytes Receive/Transmit)',
    datasource='$datasource',
    span=6,
    format='Bps',
    stack=true,
    fill=10,
    legend_show=false,
  )
  .addSeriesOverride({ alias: '/Receive/', stack: 'A' })
  .addSeriesOverride({ alias: '/Transmit/', stack: 'B', transform: 'negative-Y' })
  { tooltip+: { sort: 2 } };

local networkSaturation =
  graphPanel.new(
    'Network Saturation (Drops Receive/Transmit)',
    datasource='$datasource',
    span=6,
    format='Bps',
    stack=true,
    fill=10,
    legend_show=false,
  )
  .addSeriesOverride({ alias: '/ Receive/', stack: 'A' })
  .addSeriesOverride({ alias: '/ Transmit/', stack: 'B', transform: 'negative-Y' })
  { tooltip+: { sort: 2 } };

local diskIOUtilisation =
  graphPanel.new(
    'Disk IO Utilisation',
    datasource='$datasource',
    span=6,
    format='percentunit',
    stack=true,
    fill=10,
    legend_show=false,
  ) { tooltip+: { sort: 2 } };

local diskIOSaturation =
  graphPanel.new(
    'Disk IO Saturation',
    datasource='$datasource',
    span=6,
    format='percentunit',
    stack=true,
    fill=10,
    legend_show=false,
  ) { tooltip+: { sort: 2 } };

local diskSpaceUtilisation =
  graphPanel.new(
    'Disk Space Utilisation',
    datasource='$datasource',
    span=12,
    format='percentunit',
    stack=true,
    fill=10,
    legend_show=false,
  ) { tooltip+: { sort: 2 } };

{
  _clusterTemplate:: template.new(
    name='cluster',
    datasource='$datasource',
    query='label_values(node_time_seconds, %s)' % $._config.clusterLabel,
    current='',
    hide=if $._config.showMultiCluster then '' else '2',
    refresh=2,
    includeAll=false,
    sort=1
  ),

  grafanaDashboards+:: {
                         'node-rsrc-use.json':

                           dashboard.new(
                             '%sUSE Method / Node' % $._config.dashboardNamePrefix,
                             time_from='now-1h',
                             tags=($._config.dashboardTags),
                             timezone='utc',
                             refresh='30s',
                             graphTooltip='shared_crosshair'
                           )
                           .addTemplate(datasourceTemplate)
                           .addTemplate($._clusterTemplate)
                           .addTemplate(
                             template.new(
                               'instance',
                               '$datasource',
                               'label_values(node_exporter_build_info{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}, instance)' % $._config,
                               refresh='time',
                               sort=1
                             )
                           )
                           .addRow(
                             row.new('CPU')
                             .addPanel(CPUUtilisation.addTarget(prometheus.target('instance:node_cpu_utilisation:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config, legendFormat='Utilisation')))
                             .addPanel(CPUSaturation.addTarget(prometheus.target('instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config, legendFormat='Saturation')))
                           )
                           .addRow(
                             row.new('Memory')
                             .addPanel(memoryUtilisation.addTarget(prometheus.target('instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config, legendFormat='Utilisation')))
                             .addPanel(memorySaturation.addTarget(prometheus.target('instance:node_vmstat_pgmajfault:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config, legendFormat='Major page Faults')))
                           )
                           .addRow(
                             row.new('Network')
                             .addPanel(
                               networkUtilisation
                               .addTarget(prometheus.target('instance:node_network_receive_bytes_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config, legendFormat='Receive'))
                               .addTarget(prometheus.target('instance:node_network_transmit_bytes_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config, legendFormat='Transmit'))
                             )
                             .addPanel(
                               networkSaturation
                               .addTarget(prometheus.target('instance:node_network_receive_drop_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config, legendFormat='Receive'))
                               .addTarget(prometheus.target('instance:node_network_transmit_drop_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config, legendFormat='Transmit'))
                             )
                           )
                           .addRow(
                             row.new('Disk IO')
                             .addPanel(diskIOUtilisation.addTarget(prometheus.target('instance_device:node_disk_io_time_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config, legendFormat='{{device}}')))
                             .addPanel(diskIOSaturation.addTarget(prometheus.target('instance_device:node_disk_io_time_weighted_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} != 0' % $._config, legendFormat='{{device}}')))
                           )
                           .addRow(
                             row.new('Disk Space')
                             .addPanel(
                               diskSpaceUtilisation.addTarget(prometheus.target(
                                 |||
                                   sort_desc(1 -
                                     (
                                      max without (mountpoint, fstype) (node_filesystem_avail_bytes{%(nodeExporterSelector)s, fstype!="", instance="$instance", %(clusterLabel)s="$cluster"})
                                      /
                                      max without (mountpoint, fstype) (node_filesystem_size_bytes{%(nodeExporterSelector)s, fstype!="", instance="$instance", %(clusterLabel)s="$cluster"})
                                     ) != 0
                                   )
                                 ||| % $._config, legendFormat='{{device}}'
                               ))
                             )
                           ),

                         'node-cluster-rsrc-use.json':
                           dashboard.new(
                             '%sUSE Method / Cluster' % $._config.dashboardNamePrefix,
                             time_from='now-1h',
                             tags=($._config.dashboardTags),
                             timezone='utc',
                             refresh='30s',
                             graphTooltip='shared_crosshair'
                           )
                           .addTemplate(datasourceTemplate)
                           .addTemplate($._clusterTemplate)
                           .addRow(
                             row.new('CPU')
                             .addPanel(
                               CPUUtilisation
                               .addTarget(prometheus.target(
                                 |||
                                   ((
                                     instance:node_cpu_utilisation:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}
                                     *
                                     instance:node_num_cpu:sum{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}
                                   ) != 0 )
                                   / scalar(sum(instance:node_num_cpu:sum{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}))
                                 ||| % $._config, legendFormat='{{ instance }}'
                               ))
                             )
                             .addPanel(
                               CPUSaturation
                               .addTarget(prometheus.target(
                                 |||
                                   (
                                     instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}
                                     / scalar(count(instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}))
                                   )  != 0
                                 ||| % $._config, legendFormat='{{instance}}'
                               ))
                             )
                           )
                           .addRow(
                             row.new('Memory')
                             .addPanel(
                               memoryUtilisation
                               .addTarget(prometheus.target(
                                 |||
                                   (
                                     instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}
                                     / scalar(count(instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}))
                                   ) != 0
                                 ||| % $._config, legendFormat='{{instance}}',
                               ))
                             )
                             .addPanel(memorySaturation.addTarget(prometheus.target('instance:node_vmstat_pgmajfault:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}' % $._config, legendFormat='{{instance}}')))
                           )
                           .addRow(
                             row.new('Network')
                             .addPanel(
                               networkUtilisation
                               .addTarget(prometheus.target('instance:node_network_receive_bytes_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"} != 0' % $._config, legendFormat='{{instance}} Receive'))
                               .addTarget(prometheus.target('instance:node_network_transmit_bytes_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"} != 0' % $._config, legendFormat='{{instance}} Transmit'))
                             )
                             .addPanel(
                               networkSaturation
                               .addTarget(prometheus.target('instance:node_network_receive_drop_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"} != 0' % $._config, legendFormat='{{instance}} Receive'))
                               .addTarget(prometheus.target('instance:node_network_transmit_drop_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"} != 0' % $._config, legendFormat='{{instance}} Transmit'))
                             )
                           )
                           .addRow(
                             row.new('Disk IO')
                             .addPanel(
                               diskIOUtilisation
                               .addTarget(prometheus.target(
                                 |||
                                   (
                                     instance_device:node_disk_io_time_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}
                                     / scalar(count(instance_device:node_disk_io_time_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}))
                                   ) != 0
                                 ||| % $._config, legendFormat='{{instance}} {{device}}'
                               ))
                             )
                             .addPanel(
                               diskIOSaturation
                               .addTarget(prometheus.target(
                                 |||
                                   (
                                     instance_device:node_disk_io_time_weighted_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}
                                     / scalar(count(instance_device:node_disk_io_time_weighted_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster"}))
                                   ) != 0
                                 ||| % $._config, legendFormat='{{instance}} {{device}}'
                               ))
                             )
                           )
                           .addRow(
                             row.new('Disk Space')
                             .addPanel(
                               diskSpaceUtilisation
                               .addTarget(prometheus.target(
                                 |||
                                   sum without (device) (
                                     max without (fstype, mountpoint) ((
                                       node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s, %(fsMountpointSelector)s, %(clusterLabel)s="$cluster"}
                                       -
                                       node_filesystem_avail_bytes{%(nodeExporterSelector)s, %(fsSelector)s, %(fsMountpointSelector)s, %(clusterLabel)s="$cluster"}
                                     ) != 0)
                                   )
                                   / scalar(sum(max without (fstype, mountpoint) (node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s, %(fsMountpointSelector)s, %(clusterLabel)s="$cluster"})))
                                 ||| % $._config, legendFormat='{{instance}}'
                               ))
                             )
                           ),
                       } +
                       if $._config.showMultiCluster then {
                         'node-multicluster-rsrc-use.json':
                           dashboard.new(
                             '%sUSE Method / Multi-cluster' % $._config.dashboardNamePrefix,
                             time_from='now-1h',
                             tags=($._config.dashboardTags),
                             timezone='utc',
                             refresh='30s',
                             graphTooltip='shared_crosshair'
                           )
                           .addTemplate(datasourceTemplate)
                           .addRow(
                             row.new('CPU')
                             .addPanel(
                               CPUUtilisation
                               .addTarget(prometheus.target(
                                 |||
                                   sum(
                                     ((
                                       instance:node_cpu_utilisation:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                       *
                                       instance:node_num_cpu:sum{%(nodeExporterSelector)s}
                                     ) != 0)
                                     / scalar(sum(instance:node_num_cpu:sum{%(nodeExporterSelector)s}))
                                   ) by (%(clusterLabel)s)
                                 ||| % $._config, legendFormat='{{%(clusterLabel)s}}' % $._config
                               ))
                             )
                             .addPanel(
                               CPUSaturation
                               .addTarget(prometheus.target(
                                 |||
                                   sum((
                                     instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s}
                                     / scalar(count(instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s}))
                                   ) != 0) by (%(clusterLabel)s)
                                 ||| % $._config, legendFormat='{{%(clusterLabel)s}}' % $._config
                               ))
                             )
                           )
                           .addRow(
                             row.new('Memory')
                             .addPanel(
                               memoryUtilisation
                               .addTarget(prometheus.target(
                                 |||
                                   sum((
                                       instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s}
                                       / scalar(count(instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s}))
                                   ) != 0) by (%(clusterLabel)s)
                                 ||| % $._config, legendFormat='{{%(clusterLabel)s}}' % $._config
                               ))
                             )
                             .addPanel(
                               memorySaturation
                               .addTarget(prometheus.target(
                                 |||
                                   sum((
                                       instance:node_vmstat_pgmajfault:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                   ) != 0) by (%(clusterLabel)s)
                                 ||| % $._config, legendFormat='{{%(clusterLabel)s}}' % $._config
                               ))
                             )
                           )
                           .addRow(
                             row.new('Network')
                             .addPanel(
                               networkUtilisation
                               .addTarget(prometheus.target(
                                 |||
                                   sum((
                                       instance:node_network_receive_bytes_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                   ) != 0) by (%(clusterLabel)s)
                                 ||| % $._config, legendFormat='{{%(clusterLabel)s}} Receive' % $._config
                               ))
                               .addTarget(prometheus.target(
                                 |||
                                   sum((
                                       instance:node_network_transmit_bytes_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                   ) != 0) by (%(clusterLabel)s)
                                 ||| % $._config, legendFormat='{{%(clusterLabel)s}} Transmit' % $._config
                               ))
                             )
                             .addPanel(
                               networkSaturation
                               .addTarget(prometheus.target(
                                 |||
                                   sum((
                                       instance:node_network_receive_drop_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                   ) != 0) by (%(clusterLabel)s)
                                 ||| % $._config, legendFormat='{{%(clusterLabel)s}} Receive' % $._config
                               ))
                               .addTarget(prometheus.target(
                                 |||
                                   sum((
                                       instance:node_network_transmit_drop_excluding_lo:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                   ) != 0) by (%(clusterLabel)s)
                                 ||| % $._config, legendFormat='{{%(clusterLabel)s}} Transmit' % $._config
                               ))
                             )
                           )
                           .addRow(
                             row.new('Disk IO')
                             .addPanel(
                               diskIOUtilisation
                               .addTarget(prometheus.target(
                                 |||
                                   sum((
                                       instance_device:node_disk_io_time_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                       / scalar(count(instance_device:node_disk_io_time_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s}))
                                   ) != 0) by (%(clusterLabel)s, device)
                                 ||| % $._config, legendFormat='{{%(clusterLabel)s}} {{device}}' % $._config
                               ))
                             )
                             .addPanel(
                               diskIOSaturation
                               .addTarget(prometheus.target(
                                 |||
                                   sum((
                                     instance_device:node_disk_io_time_weighted_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s}
                                     / scalar(count(instance_device:node_disk_io_time_weighted_seconds:rate%(rateInterval)s{%(nodeExporterSelector)s}))
                                   ) != 0) by (%(clusterLabel)s, device)
                                 ||| % $._config, legendFormat='{{%(clusterLabel)s}} {{device}}' % $._config
                               ))
                             )
                           )
                           .addRow(
                             row.new('Disk Space')
                             .addPanel(
                               diskSpaceUtilisation
                               .addTarget(prometheus.target(
                                 |||
                                   sum (
                                     sum without (device) (
                                       max without (fstype, mountpoint, instance, pod) ((
                                         node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s, %(fsMountpointSelector)s} - node_filesystem_avail_bytes{%(nodeExporterSelector)s, %(fsSelector)s, %(fsMountpointSelector)s}
                                       ) != 0)
                                     )
                                     / scalar(sum(max without (fstype, mountpoint) (node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s, %(fsMountpointSelector)s})))
                                   ) by (%(clusterLabel)s)
                                 ||| % $._config, legendFormat='{{%(clusterLabel)s}}' % $._config
                               ))
                             )
                           ),
                       } else {},
}
