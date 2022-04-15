local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local grafana70 = import 'github.com/grafana/grafonnet-lib/grafonnet-7.0/grafana.libsonnet';
local gaugePanel = grafana70.panel.gauge;

{
  local prometheusDatasourceTemplate = {
    current: {
      text: 'default',
      value: 'default',
    },
    hide: 0,
    label: 'Data Source',
    name: 'prometheus_datasource',
    options: [],
    query: 'prometheus',
    refresh: 1,
    regex: '',
    type: 'datasource',
  },
  local idleCPUPanel =
    graphPanel.new(
      'CPU Usage',
      datasource='$prometheus_datasource',
      span=6,
      format='percentunit',
      max=1,
      min=0,
      stack=true,
    )
    .addTarget(prometheus.target(
      |||
        (
          (1 - sum without (mode) (rate(node_cpu_seconds_total{%(nodeExporterSelector)s, mode=~"idle|iowait|steal", instance="$instance"}[$__rate_interval])))
        / ignoring(cpu) group_left
          count without (cpu, mode) (node_cpu_seconds_total{%(nodeExporterSelector)s, mode="idle", instance="$instance"})
        )
      ||| % $._config,
      legendFormat='{{cpu}}',
      intervalFactor=5,
    )),

  local systemLoadPanel =
    graphPanel.new(
      'Load Average',
      datasource='$prometheus_datasource',
      span=6,
      format='short',
      min=0,
      fill=0,
    )
    .addTarget(prometheus.target('node_load1{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='1m load average'))
    .addTarget(prometheus.target('node_load5{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='5m load average'))
    .addTarget(prometheus.target('node_load15{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='15m load average'))
    .addTarget(prometheus.target('count(node_cpu_seconds_total{%(nodeExporterSelector)s, instance="$instance", mode="idle"})' % $._config, legendFormat='logical cores')),

  local memoryGraphPanel =
    graphPanel.new(
      'Memory Usage',
      datasource='$prometheus_datasource',
      span=9,
      format='bytes',
      stack=true,
      min=0,
    )
    .addTarget(prometheus.target(
      |||
        (
          node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance="$instance"}
        -
          node_memory_MemFree_bytes{%(nodeExporterSelector)s, instance="$instance"}
        -
          node_memory_Buffers_bytes{%(nodeExporterSelector)s, instance="$instance"}
        -
          node_memory_Cached_bytes{%(nodeExporterSelector)s, instance="$instance"}
        )
      ||| % $._config,
      legendFormat='memory used'
    ))
    .addTarget(prometheus.target('node_memory_Buffers_bytes{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='memory buffers'))
    .addTarget(prometheus.target('node_memory_Cached_bytes{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='memory cached'))
    .addTarget(prometheus.target('node_memory_MemFree_bytes{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='memory free')),

  // NOTE: avg() is used to circumvent a label change caused by a node_exporter rollout.
  local memoryGaugePanel =
    gaugePanel.new(
      title='Memory Usage',
      datasource='$prometheus_datasource',
    )
    .addTarget(prometheus.target(
      |||
        100 -
        (
          avg(node_memory_MemAvailable_bytes{%(nodeExporterSelector)s, instance="$instance"})
        /
          avg(node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance="$instance"})
        * 100
        )
      ||| % $._config,
    ))
    .addThresholdStep('rgba(50, 172, 45, 0.97)')
    .addThresholdStep('rgba(237, 129, 40, 0.89)', 80)
    .addThresholdStep('rgba(245, 54, 54, 0.9)', 90)
    .setFieldConfig(max=100, min=0, unit='percent')
    + {
      span: 3,
    },

  local diskIOPanel =
    graphPanel.new(
      'Disk I/O',
      datasource='$prometheus_datasource',
      span=6,
      min=0,
      fill=0,
    )
    // TODO: Does it make sense to have those three in the same panel?
    .addTarget(prometheus.target(
      'rate(node_disk_read_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[$__rate_interval])' % $._config,
      legendFormat='{{device}} read',
    ))
    .addTarget(prometheus.target(
      'rate(node_disk_written_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[$__rate_interval])' % $._config,
      legendFormat='{{device}} written',
    ))
    .addTarget(prometheus.target(
      'rate(node_disk_io_time_seconds_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[$__rate_interval])' % $._config,
      legendFormat='{{device}} io time',
    )) +
    {
      seriesOverrides: [
        {
          alias: '/ read| written/',
          yaxis: 1,
        },
        {
          alias: '/ io time/',
          yaxis: 2,
        },
      ],
      yaxes: [
        self.yaxe(format='bytes'),
        self.yaxe(format='s'),
      ],
    },

  // TODO: Somehow partition this by device while excluding read-only devices.
  local diskSpaceUsagePanel =
    graphPanel.new(
      'Disk Space Usage',
      datasource='$prometheus_datasource',
      span=6,
      format='bytes',
      min=0,
      fill=1,
      stack=true,
    )
    .addTarget(prometheus.target(
      |||
        sum(
          max by (device) (
            node_filesystem_size_bytes{%(nodeExporterSelector)s, instance="$instance", %(fsSelector)s}
          -
            node_filesystem_avail_bytes{%(nodeExporterSelector)s, instance="$instance", %(fsSelector)s}
          )
        )
      ||| % $._config,
      legendFormat='used',
    ))
    .addTarget(prometheus.target(
      |||
        sum(
          max by (device) (
            node_filesystem_avail_bytes{%(nodeExporterSelector)s, instance="$instance", %(fsSelector)s}
          )
        )
      ||| % $._config,
      legendFormat='available',
    )) +
    {
      seriesOverrides: [
        {
          alias: 'used',
          color: '#E0B400',
        },
        {
          alias: 'available',
          color: '#73BF69',
        },
      ],
    },

  local networkReceivedPanel =
    graphPanel.new(
      'Network Received',
      datasource='$prometheus_datasource',
      span=6,
      format='bytes',
      min=0,
      fill=0,
    )
    .addTarget(prometheus.target(
      'rate(node_network_receive_bytes_total{%(nodeExporterSelector)s, instance="$instance", device!="lo"}[$__rate_interval])' % $._config,
      legendFormat='{{device}}',
    )),

  local networkTransmittedPanel =
    graphPanel.new(
      'Network Transmitted',
      datasource='$prometheus_datasource',
      span=6,
      format='bytes',
      min=0,
      fill=0,
    )
    .addTarget(prometheus.target(
      'rate(node_network_transmit_bytes_total{%(nodeExporterSelector)s, instance="$instance", device!="lo"}[$__rate_interval])' % $._config,
      legendFormat='{{device}}',
    )),

  local cpuRow =
    row.new('CPU')
    .addPanel(idleCPUPanel)
    .addPanel(systemLoadPanel),

  local memoryRow =
    row.new('Memory')
    .addPanel(memoryGraphPanel)
    .addPanel(memoryGaugePanel),

  local diskRow =
    row.new('Disk')
    .addPanel(diskIOPanel)
    .addPanel(diskSpaceUsagePanel),

  local networkRow =
    row.new('Network')
    .addPanel(networkReceivedPanel)
    .addPanel(networkTransmittedPanel),

  local instanceTemplate = template.new(
    'instance',
    '$prometheus_datasource',
    'label_values(node_exporter_build_info{%(nodeExporterSelector)s}, instance)' % $._config,
    refresh='time',
    label='Instance',
  ),

  local NodeDashboard =
    dashboard.new(
      '%sNodes' % $._config.dashboardNamePrefix,
      time_from='now-1h',
      tags=($._config.dashboardTags),
      timezone='utc',
      refresh='30s',
      graphTooltip='shared_crosshair'
    )
    .addTemplate(prometheusDatasourceTemplate)
    .addTemplate(instanceTemplate)
    .addRow(cpuRow)
    .addRow(memoryRow)
    .addRow(diskRow)
    .addRow(networkRow),

  grafanaDashboards+::
    if $._config.enableLokiLogs then {
      local lokiMixin = import '../lib/lokimixin/loki-mixin.libsonnet',
      local l = lokiMixin.new('%(nodeExporterSelector)s, instance="$instance"' % $._config),
      'nodes.json':
        NodeDashboard
        .addTemplate(l.lokiDatasourceTemplate)
        .addTemplate(l.unitTemplate)
        .addRow(l.lokiDirectLogRow)
        .addRow(l.lokiJournalLogRow),
    }
    else {
      'nodes.json':
        NodeDashboard,

    },
}
