local grafana = import 'github.com/grafana/grafonnet/gen/grafonnet-latest/main.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.panel.row;
local prometheus = grafana.query.prometheus;
local variable = dashboard.variable;
local template = grafana.template;

local graphPanel = grafana.graphPanel;
local gaugePanel = grafana70.panel.gauge;
local table = grafana70.panel.table;

local timeSeriesPanel = grafana.panel.timeSeries;


{

  new(config=null, platform=null, uid=null):: {

    local datasource = variable.datasource.new(
      'datasource', 'prometheus'
    ),

    local clusterVariablePrototype = variable.query.new('cluster')
    + variable.query.withDatasourceFromVariable(datasource)
    + (if config.showMultiCluster then variable.query.generalOptions.showOnDashboard.withLabelAndValue() else variable.query.generalOptions.showOnDashboard.withNothing())
    + variable.query.refresh.onTime()
    + variable.generalOptions.withLabel("Cluster"),

    local clusterVariable = 
      if platform == 'Darwin' then
        clusterVariablePrototype + 
        variable.query.queryTypes.withLabelValues(
          config,
          'node_uname_info{%(nodeExporterSelector)s, sysname="Darwin"}, %(clusterLabel)s',
        )
      else
      clusterVariablePrototype + 
        variable.query.queryTypes.withLabelValues(
          config,
          'node_uname_info{%(nodeExporterSelector)s, sysname!="Darwin"}, %(clusterLabel)s',
        ),

    local instanceTemplatePrototype =
      template.new(
        'instance',
        '$datasource',
        '',
        refresh='time',
        label='Instance',
      ),
    local instanceTemplate =
      if platform == 'Darwin' then
        instanceTemplatePrototype
        { query: 'label_values(node_uname_info{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster", sysname="Darwin"}, instance)' % config }
      else
        instanceTemplatePrototype
        { query: 'label_values(node_uname_info{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster", sysname!="Darwin"}, instance)' % config },

    local idleCPU =
      graphPanel.new(
        'CPU Usage',
        datasource='$datasource',
        span=6,
        format='percentunit',
        max=1,
        min=0,
        stack=true,
      )
      .addTarget(prometheus.target(
        |||
          (
            (1 - sum without (mode) (rate(node_cpu_seconds_total{%(nodeExporterSelector)s, mode=~"idle|iowait|steal", instance="$instance", %(clusterLabel)s="$cluster"}[$__rate_interval])))
          / ignoring(cpu) group_left
            count without (cpu, mode) (node_cpu_seconds_total{%(nodeExporterSelector)s, mode="idle", instance="$instance", %(clusterLabel)s="$cluster"})
          )
        ||| % config,
        legendFormat='{{cpu}}',
        intervalFactor=5,
      )),

    local systemLoad =
      graphPanel.new(
        'Load Average',
        datasource='$datasource',
        span=6,
        format='short',
        min=0,
        fill=0,
      )
      .addTarget(prometheus.target('node_load1{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config, legendFormat='1m load average'))
      .addTarget(prometheus.target('node_load5{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config, legendFormat='5m load average'))
      .addTarget(prometheus.target('node_load15{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config, legendFormat='15m load average'))
      .addTarget(prometheus.target('count(node_cpu_seconds_total{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", mode="idle"})' % config, legendFormat='logical cores')),

    local memoryGraphPanelPrototype =
      graphPanel.new(
        'Memory Usage',
        datasource='$datasource',
        span=9,
        format='bytes',
        min=0,
      ),
    local memoryGraph =
      if platform == 'Linux' then
        memoryGraphPanelPrototype { stack: true }
        .addTarget(prometheus.target(
          |||
            (
              node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}
            -
              node_memory_MemFree_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}
            -
              node_memory_Buffers_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}
            -
              node_memory_Cached_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}
            )
          ||| % config,
          legendFormat='memory used'
        ))
        .addTarget(prometheus.target('node_memory_Buffers_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config, legendFormat='memory buffers'))
        .addTarget(prometheus.target('node_memory_Cached_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config, legendFormat='memory cached'))
        .addTarget(prometheus.target('node_memory_MemFree_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config, legendFormat='memory free'))
      else if platform == 'Darwin' then
        // not useful to stack
        memoryGraphPanelPrototype { stack: false }
        .addTarget(prometheus.target('node_memory_total_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config, legendFormat='Physical Memory'))
        .addTarget(prometheus.target(
          |||
            (
                node_memory_internal_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} -
                node_memory_purgeable_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} +
                node_memory_wired_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} +
                node_memory_compressed_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}
            )
          ||| % config, legendFormat='Memory Used'
        ))
        .addTarget(prometheus.target(
          |||
            (
                node_memory_internal_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} -
                node_memory_purgeable_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}
            )
          ||| % config, legendFormat='App Memory'
        ))
        .addTarget(prometheus.target('node_memory_wired_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config, legendFormat='Wired Memory'))
        .addTarget(prometheus.target('node_memory_compressed_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config, legendFormat='Compressed'))
      else if platform == 'AIX' then
        memoryGraphPanelPrototype { stack: false }
        .addTarget(prometheus.target('node_memory_total_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config, legendFormat='Physical Memory'))
        .addTarget(prometheus.target(
          |||
            (
                node_memory_total_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} -
                node_memory_available_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}
            )
          ||| % config, legendFormat='Memory Used'
        )),


    // NOTE: avg() is used to circumvent a label change caused by a node_exporter rollout.
    local memoryGaugePanelPrototype =
      gaugePanel.new(
        title='Memory Usage',
        datasource='$datasource',
      )
      .addThresholdStep('rgba(50, 172, 45, 0.97)')
      .addThresholdStep('rgba(237, 129, 40, 0.89)', 80)
      .addThresholdStep('rgba(245, 54, 54, 0.9)', 90)
      .setFieldConfig(max=100, min=0, unit='percent')
      + {
        span: 3,
      },

    local memoryGauge =
      if platform == 'Linux' then
        memoryGaugePanelPrototype

        .addTarget(prometheus.target(
          |||
            100 -
            (
              avg(node_memory_MemAvailable_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}) /
              avg(node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"})
            * 100
            )
          ||| % config,
        ))

      else if platform == 'Darwin' then
        memoryGaugePanelPrototype
        .addTarget(prometheus.target(
          |||
            (
                (
                  avg(node_memory_internal_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}) -
                  avg(node_memory_purgeable_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}) +
                  avg(node_memory_wired_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}) +
                  avg(node_memory_compressed_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"})
                ) /
                avg(node_memory_total_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"})
            )
            *
            100
          ||| % config
        ))
      else if platform == 'AIX' then
        memoryGaugePanelPrototype
        .addTarget(prometheus.target(
          |||
            100 -
            (
              avg(node_memory_available_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}) /
              avg(node_memory_total_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"})
              * 100
            )
          ||| % config
        )),


    local diskIO =
      graphPanel.new(
        'Disk I/O',
        datasource='$datasource',
        span=6,
        min=0,
        fill=0,
      )
      // TODO: Does it make sense to have those three in the same panel?
      .addTarget(prometheus.target(
        'rate(node_disk_read_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", %(diskDeviceSelector)s}[$__rate_interval])' % config,
        legendFormat='{{device}} read',
        intervalFactor=1,
      ))
      .addTarget(prometheus.target(
        'rate(node_disk_written_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", %(diskDeviceSelector)s}[$__rate_interval])' % config,
        legendFormat='{{device}} written',
        intervalFactor=1,
      ))
      .addTarget(prometheus.target(
        'rate(node_disk_io_time_seconds_total{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", %(diskDeviceSelector)s}[$__rate_interval])' % config,
        legendFormat='{{device}} io time',
        intervalFactor=1,
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
          self.yaxe(format='Bps'),
          self.yaxe(format='percentunit'),
        ],
      },

    local diskSpaceUsage =
      table.new(
        title='Disk Space Usage',
        datasource='$datasource',
      )
      .setFieldConfig(unit='decbytes')
      .addThresholdStep(color='green', value=null)
      .addThresholdStep(color='yellow', value=0.8)
      .addThresholdStep(color='red', value=0.9)
      .addTarget(prometheus.target(
        |||
          max by (mountpoint) (node_filesystem_size_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", %(fsSelector)s, %(fsMountpointSelector)s})
        ||| % config,
        legendFormat='',
        instant=true,
        format='table'
      ))
      .addTarget(prometheus.target(
        |||
          max by (mountpoint) (node_filesystem_avail_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", %(fsSelector)s, %(fsMountpointSelector)s})
        ||| % config,
        legendFormat='',
        instant=true,
        format='table'
      ))
      .addOverride(
        matcher={
          id: 'byName',
          options: 'Mounted on',
        },
        properties=[
          {
            id: 'custom.width',
            value: 260,
          },
        ],
      )
      .addOverride(
        matcher={
          id: 'byName',
          options: 'Size',
        },
        properties=[

          {
            id: 'custom.width',
            value: 93,
          },

        ],
      )
      .addOverride(
        matcher={
          id: 'byName',
          options: 'Used',
        },
        properties=[
          {
            id: 'custom.width',
            value: 72,
          },
        ],
      )
      .addOverride(
        matcher={
          id: 'byName',
          options: 'Available',
        },
        properties=[
          {
            id: 'custom.width',
            value: 88,
          },
        ],
      )

      .addOverride(
        matcher={
          id: 'byName',
          options: 'Used, %',
        },
        properties=[
          {
            id: 'unit',
            value: 'percentunit',
          },
          {
            id: 'custom.displayMode',
            value: 'gradient-gauge',
          },
          {
            id: 'max',
            value: 1,
          },
          {
            id: 'min',
            value: 0,
          },
        ]
      )
      + { span: 6 }
      + {
        transformations: [
          {
            id: 'groupBy',
            options: {
              fields: {
                'Value #A': {
                  aggregations: [
                    'lastNotNull',
                  ],
                  operation: 'aggregate',
                },
                'Value #B': {
                  aggregations: [
                    'lastNotNull',
                  ],
                  operation: 'aggregate',
                },
                mountpoint: {
                  aggregations: [],
                  operation: 'groupby',
                },
              },
            },
          },
          {
            id: 'merge',
            options: {},
          },
          {
            id: 'calculateField',
            options: {
              alias: 'Used',
              binary: {
                left: 'Value #A (lastNotNull)',
                operator: '-',
                reducer: 'sum',
                right: 'Value #B (lastNotNull)',
              },
              mode: 'binary',
              reduce: {
                reducer: 'sum',
              },
            },
          },
          {
            id: 'calculateField',
            options: {
              alias: 'Used, %',
              binary: {
                left: 'Used',
                operator: '/',
                reducer: 'sum',
                right: 'Value #A (lastNotNull)',
              },
              mode: 'binary',
              reduce: {
                reducer: 'sum',
              },
            },
          },
          {
            id: 'organize',
            options: {
              excludeByName: {},
              indexByName: {},
              renameByName: {
                'Value #A (lastNotNull)': 'Size',
                'Value #B (lastNotNull)': 'Available',
                mountpoint: 'Mounted on',
              },
            },
          },
          {
            id: 'sortBy',
            options: {
              fields: {},
              sort: [
                {
                  field: 'Mounted on',
                },
              ],
            },
          },
        ],
      },


    local networkReceived =
      graphPanel.new(
        'Network Received',
        description='Network received (bits/s)',
        datasource='$datasource',
        span=6,
        format='bps',
        min=0,
        fill=0,
      )
      .addTarget(prometheus.target(
        'rate(node_network_receive_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", device!="lo"}[$__rate_interval]) * 8' % config,
        legendFormat='{{device}}',
        intervalFactor=1,
      )),

    local networkTransmitted =
      graphPanel.new(
        'Network Transmitted',
        description='Network transmitted (bits/s)',
        datasource='$datasource',
        span=6,
        format='bps',
        min=0,
        fill=0,
      )
      .addTarget(prometheus.target(
        'rate(node_network_transmit_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", device!="lo"}[$__rate_interval]) * 8' % config,
        legendFormat='{{device}}',
        intervalFactor=1,
      )),

    local cpuRow =
      row.new('CPU')
      .addPanel(idleCPU)
      .addPanel(systemLoad),

    local memoryRow =
      row.new('Memory')
      .addPanel(memoryGraph)
      .addPanel(memoryGauge),

    local diskRow =
      row.new('Disk')
      .addPanel(diskIO)
      .addPanel(diskSpaceUsage),

    local networkRow =
      row.new('Network')
      .addPanel(networkReceived)
      .addPanel(networkTransmitted),

    local rows =
      [
        cpuRow,
        memoryRow,
        diskRow,
        networkRow,
      ],

    local templates =
      [
        prometheusDatasourceTemplate,
        clusterTemplate,
        instanceTemplate,
      ],


    dashboard: if platform == 'Linux' then
      dashboard.new(
        '%sNodes' % config.dashboardNamePrefix,
        time_from='now-1h',
        tags=(config.dashboardTags),
        timezone='utc',
        refresh='30s',
        uid=std.md5(uid),
        graphTooltip='shared_crosshair'
      )
      .addTemplates(templates)
      .addRows(rows)
    else if platform == 'Darwin' then
      dashboard.new(
        '%sMacOS' % config.dashboardNamePrefix,
        time_from='now-1h',
        tags=(config.dashboardTags),
        timezone='utc',
        refresh='30s',
        uid=std.md5(uid),
        graphTooltip='shared_crosshair'
      )
      .addTemplates(templates)
      .addRows(rows)
    else if platform == 'AIX' then
      dashboard.new(
        '%sAIX' % config.dashboardNamePrefix,
        time_from='now-1h',
        tags=(config.dashboardTags),
        timezone='utc',
        refresh='30s',
        uid=std.md5(uid),
        graphTooltip='shared_crosshair'
      )
      .addTemplates(templates)
      .addRows(rows),

  },
}
