local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local statPanel = grafana.statPanel;
local grafana70 = import 'github.com/grafana/grafonnet-lib/grafonnet-7.0/grafana.libsonnet';
local gaugePanel = grafana70.panel.gauge;
local table = grafana70.panel.table;
local nodePanels = import 'panels-lib/panels.libsonnet';
local commonPanels = import 'panels-lib/common/panels.libsonnet';
local nodeTimeseries = nodePanels.timeseries;
local nodeTemplates = import 'templates.libsonnet';
{

  new(config=null, platform=null):: {


    local uptimePanel =
      commonPanels.uptime.new(
        statPanel.new(
          'Uptime',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        |||
          time() - node_boot_time_seconds{%(nodeExporterSelector)s, instance="$instance"}
        ||| % config,
      )) { span: 3, height: '100px' },

    local cpuCountPanel =
      commonPanels.info.new(
        statPanel.new(
          'CPU Count',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        expr=
        |||
          count(count by (cpu)(node_cpu_seconds_total{%(nodeExporterSelector)s, instance="$instance"}))
        ||| % config,
      ))
      { span: 3, height: '100px' },
    local memoryTotalPanel =
      commonPanels.info.new(
        statPanel.new(
          'Memory Total',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        |||
          node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance="$instance"}
        ||| % config,
      ))
      .withUnits('bytes')
      .withDecimals(0)
      { span: 3, height: '100px' },

    local osPanel =

      commonPanels.info.new(
        statPanel.new(
          'OS',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        expr=
        |||
          node_os_info{%(nodeExporterSelector)s, instance="$instance"}
        ||| % config,
        format='table'
      ))
      { options+: { reduceOptions+: { fields: '/^pretty_name$/' } } }
      { span: 3, height: '100px' },

    local nodeNamePanel =

      commonPanels.info.new(
        statPanel.new(
          'Hostname',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        expr=
        |||
          node_uname_info{%(nodeExporterSelector)s, instance="$instance"}
        ||| % config,
        format='table'
      ))
      { options+: { reduceOptions+: { fields: '/^nodename$/' } } }
      { span: 3, height: '100px' },

    local kernelVersionPanel =

      commonPanels.info.new(
        statPanel.new(
          'Kernel version',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        expr=
        |||
          node_uname_info{%(nodeExporterSelector)s, instance="$instance"}
        ||| % config,
        format='table'
      ))
      { options+: { reduceOptions+: { fields: '/^release$/' } } }
      { span: 3, height: '100px' },

    local totalSwapPanel =

      commonPanels.info.new(
        statPanel.new(
          'Total swap',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        expr=
        |||
          node_memory_SwapTotal_bytes{%(nodeExporterSelector)s, instance="$instance"}
        ||| % config,
      ))
      .withUnits('bytes')
      .withDecimals(0)
      { span: 3, height: '100px' },

    local totalRootFSPanel =

      commonPanels.info.new(
        statPanel.new(
          'Root mount size',
          datasource='$datasource',
        )
      )
      .addTarget(prometheus.target(
        expr=
        |||
          node_filesystem_size_bytes{%(nodeExporterSelector)s, instance="$instance", mountpoint="/",fstype!="rootfs"}
        ||| % config,
      ))
      .withUnits('bytes')
      .withDecimals(0)
      { span: 3, height: '100px' },

    local idleCPU =
      nodeTimeseries.new(
        graphPanel.new(
          'CPU Usage',
          datasource='$datasource',
          span=6,
        )
      )
      .withUnits('percentunit')
      .withStacking('normal')
      .withMin(0)
      .withMax(1)
      .addTarget(prometheus.target(
        |||
          (
            (1 - sum without (mode) (rate(node_cpu_seconds_total{%(nodeExporterSelector)s, mode=~"idle|iowait|steal", instance="$instance"}[$__rate_interval])))
          / ignoring(cpu) group_left
            count without (cpu, mode) (node_cpu_seconds_total{%(nodeExporterSelector)s, mode="idle", instance="$instance"})
          )
        ||| % config,
        legendFormat='cpu {{cpu}}',
        intervalFactor=5,
      )),

    local systemLoad =
      nodeTimeseries.new(
        graphPanel.new(
          'Load Average',
          datasource='$datasource',
          span=6,
          fill=0,
        )
      )
      .withUnits('short')
      .withMin(0)
      .withFillOpacity(0)
      .addTarget(prometheus.target('node_load1{%(nodeExporterSelector)s, instance="$instance"}' % config, legendFormat='1m load average'))
      .addTarget(prometheus.target('node_load5{%(nodeExporterSelector)s, instance="$instance"}' % config, legendFormat='5m load average'))
      .addTarget(prometheus.target('node_load15{%(nodeExporterSelector)s, instance="$instance"}' % config, legendFormat='15m load average'))
      .addTarget(prometheus.target('count(node_cpu_seconds_total{%(nodeExporterSelector)s, instance="$instance", mode="idle"})' % config, legendFormat='logical cores'))
      .addOverride(
        matcher={
          id: 'byName',
          options: 'logical cores',
        },
        properties=[
          {
            id: 'custom.lineStyle',
            value: {
              fill: 'dash',
              dash: [
                10,
                10,
              ],
            },
          },
        ]
      ),
    local memoryGraphPanelPrototype = nodeTimeseries.new(
      graphPanel.new(
        'Memory Usage',
        datasource='$datasource',
        span=9,
      )
    )
                                      .withMin(0)
                                      .withUnits('bytes'),
    local memoryGraph =
      if platform == 'Linux' then
        memoryGraphPanelPrototype { stack: true }
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
          ||| % config,
          legendFormat='memory used'
        ))
        .addTarget(prometheus.target('node_memory_Buffers_bytes{%(nodeExporterSelector)s, instance="$instance"}' % config, legendFormat='memory buffers'))
        .addTarget(prometheus.target('node_memory_Cached_bytes{%(nodeExporterSelector)s, instance="$instance"}' % config, legendFormat='memory cached'))
        .addTarget(prometheus.target('node_memory_MemFree_bytes{%(nodeExporterSelector)s, instance="$instance"}' % config, legendFormat='memory free'))
      else if platform == 'Darwin' then
        // not useful to stack
        memoryGraphPanelPrototype { stack: false }
        .addTarget(prometheus.target('node_memory_total_bytes{%(nodeExporterSelector)s, instance="$instance"}' % config, legendFormat='Physical Memory'))
        .addTarget(prometheus.target(
          |||
            (
                node_memory_internal_bytes{%(nodeExporterSelector)s, instance="$instance"} -
                node_memory_purgeable_bytes{%(nodeExporterSelector)s, instance="$instance"} +
                node_memory_wired_bytes{%(nodeExporterSelector)s, instance="$instance"} +
                node_memory_compressed_bytes{%(nodeExporterSelector)s, instance="$instance"}
            )
          ||| % config, legendFormat='Memory Used'
        ))
        .addTarget(prometheus.target(
          |||
            (
                node_memory_internal_bytes{%(nodeExporterSelector)s, instance="$instance"} -
                node_memory_purgeable_bytes{%(nodeExporterSelector)s, instance="$instance"}
            )
          ||| % config, legendFormat='App Memory'
        ))
        .addTarget(prometheus.target('node_memory_wired_bytes{%(nodeExporterSelector)s, instance="$instance"}' % config, legendFormat='Wired Memory'))
        .addTarget(prometheus.target('node_memory_compressed_bytes{%(nodeExporterSelector)s, instance="$instance"}' % config, legendFormat='Compressed')),

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
              avg(node_memory_MemAvailable_bytes{%(nodeExporterSelector)s, instance="$instance"}) /
              avg(node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance="$instance"})
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
                  avg(node_memory_internal_bytes{%(nodeExporterSelector)s, instance="$instance"}) -
                  avg(node_memory_purgeable_bytes{%(nodeExporterSelector)s, instance="$instance"}) +
                  avg(node_memory_wired_bytes{%(nodeExporterSelector)s, instance="$instance"}) +
                  avg(node_memory_compressed_bytes{%(nodeExporterSelector)s, instance="$instance"})
                ) /
                avg(node_memory_total_bytes{%(nodeExporterSelector)s, instance="$instance"})
            )
            *
            100
          ||| % config
        )),

    local diskIO =
      nodeTimeseries.new(
        graphPanel.new(
          'Disk I/O',
          datasource='$datasource',
          span=6,
        )
      )
      .withFillOpacity(0)
      .withMin(0)
      // TODO: Does it make sense to have those three in the same panel?
      .addTarget(prometheus.target(
        'rate(node_disk_read_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[$__rate_interval])' % config,
        legendFormat='{{device}} read',
      ))
      .addTarget(prometheus.target(
        'rate(node_disk_written_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[$__rate_interval])' % config,
        legendFormat='{{device}} written',
      ))
      .addTarget(prometheus.target(
        'rate(node_disk_io_time_seconds_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[$__rate_interval])' % config,
        legendFormat='{{device}} io time',
      ))
      .addOverride(
        matcher={
          id: 'byRegexp',
          options: '/ read| written/',
        },
        properties=[
          {
            id: 'unit',
            value: 'bps',
          },
        ]
      )
      .addOverride(
        matcher={
          id: 'byRegexp',
          options: '/ io time/',
        },
        properties=[
          {
            id: 'unit',
            value: 'percentunit',
          },
          {
            id: 'custom.axisSoftMax',
            value: 1,
          },
          {
            id: 'custom.drawStyle',
            value: 'points',
          },
        ]
      ),

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
          max by (mountpoint) (node_filesystem_size_bytes{%(nodeExporterSelector)s, instance="$instance", %(fsSelector)s, %(fsMountpointSelector)s})
        ||| % config,
        legendFormat='',
        instant=true,
        format='table'
      ))
      .addTarget(prometheus.target(
        |||
          max by (mountpoint) (node_filesystem_avail_bytes{%(nodeExporterSelector)s, instance="$instance", %(fsSelector)s, %(fsMountpointSelector)s})
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
      nodeTimeseries.new(
        graphPanel.new(
          'Network Received',
          description='Network received (bits/s)',
          datasource='$datasource',
          span=6,
        )
      )
      .withFillOpacity(10)
      .withMin(0)
      .withUnits('bps')
      .addTarget(prometheus.target(
        'rate(node_network_receive_bytes_total{%(nodeExporterSelector)s, instance="$instance", device!="lo"}[$__rate_interval]) * 8' % config,
        legendFormat='{{device}}',
        intervalFactor=1,
      )),

    local networkTransmitted =
      nodeTimeseries.new(
        graphPanel.new(
          'Network Transmitted',
          description='Network transmitted (bits/s)',
          datasource='$datasource',
          span=6,
        )
      )
      .withFillOpacity(10)
      .withMin(0)
      .withUnits('bps')
      .addTarget(prometheus.target(
        'rate(node_network_transmit_bytes_total{%(nodeExporterSelector)s, instance="$instance", device!="lo"}[$__rate_interval]) * 8' % config,
        legendFormat='{{device}}',
        intervalFactor=1,
      )),

    local infoRow =
      row.new('Overview')
      .addPanel(uptimePanel)
      .addPanel(nodeNamePanel)
      .addPanel(kernelVersionPanel)
      .addPanel(osPanel)
      .addPanel(cpuCountPanel)
      .addPanel(memoryTotalPanel)
      .addPanel(totalSwapPanel)
      .addPanel(totalRootFSPanel),

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
        infoRow,
        cpuRow,
        memoryRow,
        diskRow,
        networkRow,
      ],

    local templates = nodeTemplates.new(config=config, platform=platform).templates,

    dashboard: if platform == 'Linux' then
      dashboard.new(
        '%sNodes' % config.dashboardNamePrefix,
        time_from=config.dashboardInterval,
        tags=(config.dashboardTags),
        timezone=config.dashboardTimezone,
        refresh=config.dashboardRefresh,
        graphTooltip='shared_crosshair',
        uid='nodes'
      ) { editable: true }
      // .addLink(grafana.link.dashboards(
      //   asDropdown=true,
      //   title='Other Node dashboards',
      //   includeVars=true,
      //   keepTime=true,
      //   tags=(config.dashboardTags),
      // ))
      .addTemplates(templates)
      .addRows(rows)
    else if platform == 'Darwin' then
      dashboard.new(
        '%sMacOS' % config.dashboardNamePrefix,
        time_from=config.dashboardInterval,
        tags=(config.dashboardTags),
        timezone=config.dashboardTimezone,
        refresh=config.dashboardRefresh,
        graphTooltip='shared_crosshair',
        uid='nodes-darwin'
      )
      .addTemplates(templates)
      .addRows(rows),

  },
}
