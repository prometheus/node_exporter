local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local grafana70 = import 'github.com/grafana/grafonnet-lib/grafonnet-7.0/grafana.libsonnet';
local gaugePanel = grafana70.panel.gauge;
local loki = grafana.loki;
local logPanel = grafana.logPanel;

local c = import '../config.libsonnet';

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
};

{
  _idleCPUPanel::
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

  _systemLoadPanel::
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

  _memoryGraphPanel::
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
  _memoryGaugePanel::
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

  _diskIOPanel:
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
  _diskSpaceUsagePanel:
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

  _networkReceivedPanel::
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

  _networkTransmittedPanel::
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

  _cpuRow::
    row.new('CPU')
    .addPanel($._idleCPUPanel)
    .addPanel($._systemLoadPanel),

  _memoryRow::
    row.new('Memory')
    .addPanel($._memoryGraphPanel)
    .addPanel($._memoryGaugePanel),

  _diskRow::
    row.new('Disk')
    .addPanel($._diskIOPanel)
    .addPanel($._diskSpaceUsagePanel),

  _networkRow::
    row.new('Network')
    .addPanel($._networkReceivedPanel)
    .addPanel($._networkTransmittedPanel),

  _instanceTemplate:: template.new(
    'instance',
    '$prometheus_datasource',
    'label_values(node_exporter_build_info{%(nodeExporterSelector)s}, instance)' % $._config,
    refresh='time',
    label='Instance',
  ),

  _NodeDashboard:: dashboard.new(
    '%sNodes' % $._config.dashboardNamePrefix,
    time_from='now-1h',
    tags=($._config.dashboardTags),
    timezone='utc',
    refresh='30s',
    graphTooltip='shared_crosshair'
  )
                   .addTemplate(prometheusDatasourceTemplate)
                   .addRow($._cpuRow)
                   .addRow($._memoryRow)
                   .addRow($._diskRow)
                   .addRow($._networkRow),

  grafanaDashboards+::
    if !$._config.enableLokiLogs then {
      'nodes.json':

        $._NodeDashboard
        .addTemplate($._instanceTemplate),
    }
    else {
      'nodes.json':

        local lokiDatasourceTemplate = {
          current:
            {
              text: 'Loki',
              value: 'Loki',
            },
          name: 'loki_datasource',
          label: 'Loki Data Source',
          options: [],
          query: 'loki',
          hide: 0,
          refresh: 1,
          regex: '',
          type: 'datasource',
        };

        local unitTemplate = template.new(
          'unit',
          '$loki_datasource',
          'label_values(unit)',
          label='Systemd Unit',
          refresh='time',
          includeAll=true,
          multi=true,
          allValues='.+',
        );

        local syslog =
          logPanel.new(
            'syslog Errors',
            datasource='$loki_datasource',
          )
          .addTarget(
            loki.target('{filename=~"/var/log/syslog*|/var/log/messages*", %(nodeExporterSelector)s, instance=~"$instance"} |~".+(?i)error(?-i).+"' % $._config)
          );

        local authlog =
          logPanel.new(
            'authlog',
            datasource='$loki_datasource',
          )
          .addTarget(
            loki.target('{filename=~"/var/log/auth.log|/var/log/secure", %(nodeExporterSelector)s, instance=~"$instance"}' % $._config)
          );

        local kernellog =
          logPanel.new(
            'Kernel logs',
            datasource='$loki_datasource',
          )
          .addTarget(
            loki.target('{filename=~"/var/log/kern.log*", %(nodeExporterSelector)s, instance=~"$instance"}' % $._config)
          );

        local journalsyslog =
          logPanel.new(
            'Journal syslogs',
            datasource='$loki_datasource',
          )
          .addTarget(
            loki.target('{transport="syslog", %(nodeExporterSelector)s, instance=~"$instance"}' % $._config)
          );

        local journalkernel =
          logPanel.new(
            'Journal Kernel logs',
            datasource='$loki_datasource',
          )
          .addTarget(
            loki.target('{transport="kernel", %(nodeExporterSelector)s, instance=~"$instance"}' % $._config)
          );

        local journalstdout =
          logPanel.new(
            'Journal stdout Errors',
            datasource='$loki_datasource',
          )
          .addTarget(
            loki.target('{transport="stdout", %(nodeExporterSelector)s, instance=~"$instance", unit=~"$unit"} |~".+(?i)error(?-i).+"' % $._config)
          );

        local lokiDirectLogRow =
          row.new(
            'Loki Direct Log Scrapes'
          )
          .addPanel(syslog)
          .addPanel(authlog)
          .addPanel(kernellog);

        local lokiJournalLogRow =
          row.new(
            'Loki Journal Log Scrapes'
          )
          .addPanel(journalsyslog)
          .addPanel(journalkernel)
          .addPanel(journalstdout);

        $._NodeDashboard
        .addTemplate(lokiDatasourceTemplate)
        .addTemplate($._instanceTemplate)
        .addTemplate(unitTemplate)
        .addRow(lokiDirectLogRow)
        .addRow(lokiJournalLogRow),
    },
}
