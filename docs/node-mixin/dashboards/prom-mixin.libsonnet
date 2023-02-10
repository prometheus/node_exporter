local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local statPanel = grafana.statPanel;
local nodePanels = import '../lib/panels/panels.libsonnet';
local commonPanels = import '../lib/panels/common/panels.libsonnet';
local nodeTimeseries = nodePanels.timeseries;
local common = import '../lib/common.libsonnet';
local nodeTemplates = common.templates;

{

  new(config=null, platform=null):: {

    local c = common.new(config=config, platform=platform),
    local commonPromTarget = c.commonPromTarget,
    local templates = c.templates,
    local q = c.queries,

    local uptimePanel =
      commonPanels.uptimeStat.new()
      .addTarget(commonPromTarget(expr=q.uptime)),

    local cpuCountPanel =
      commonPanels.infoStat.new('CPU Count')
      .addTarget(commonPromTarget(expr=q.cpuCount)),

    local memoryTotalPanel =
      commonPanels.infoStat.new('Memory Total')
      .addTarget(commonPromTarget(expr=q.memoryTotal))
      .withUnits('bytes')
      .withDecimals(0),

    local osPanel =
      commonPanels.infoStat.new('OS')
      .addTarget(commonPromTarget(
        expr=q.osInfo, format='table'
      )) { options+: { reduceOptions+: { fields: '/^pretty_name$/' } } },

    local nodeNamePanel =
      commonPanels.infoStat.new('Hostname')
      .addTarget(commonPromTarget(
        expr=q.nodeInfo, format='table'
      ))
      { options+: { reduceOptions+: { fields: '/^nodename$/' } } },

    local kernelVersionPanel =

      commonPanels.infoStat.new('Kernel version')
      .addTarget(commonPromTarget(
        expr=q.nodeInfo, format='table'
      ))
      { options+: { reduceOptions+: { fields: '/^release$/' } } }
    ,

    local totalSwapPanel =
      commonPanels.infoStat.new('Total swap')
      .addTarget(commonPromTarget(
        expr=q.memorySwapTotal
      ))
      .withUnits('bytes')
      .withDecimals(0),

    local totalRootFSPanel =
      commonPanels.infoStat.new('Root mount size')
      .addTarget(commonPromTarget(
        expr=q.fsSizeTotalRoot,
      ))
      .withUnits('bytes')
      .withDecimals(0),

    local cpuStatPanel =
      commonPanels.percentUsageStat.new('CPU Usage')
      .addTarget(commonPromTarget(
        expr=q.cpuUsage
      )),


    local idleCPU =
      nodePanels.timeseries.new('CPU Usage')
      .withUnits('percentunit')
      .withStacking('normal')
      .withMin(0)
      .withMax(1)
      .addTarget(commonPromTarget(
        expr=q.cpuUsagePerCore,
        legendFormat='cpu {{cpu}}',
      )),

    local systemLoad =
      nodePanels.timeseries.new('Load Average')
      .withUnits('short')
      .withMin(0)
      .withFillOpacity(0)
      .addTarget(commonPromTarget(q.systemLoad1, legendFormat='1m load average'))
      .addTarget(commonPromTarget(q.systemLoad5, legendFormat='5m load average'))
      .addTarget(commonPromTarget(q.systemLoad15, legendFormat='15m load average'))
      .addTarget(commonPromTarget(q.cpuCount, legendFormat='logical cores'))
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
    local memoryGraphPanelPrototype = nodePanels.timeseries.new('Memory Usage')
                                      .withMin(0)
                                      .withUnits('bytes'),
    local memoryGraph =
      if platform == 'Linux' then
        memoryGraphPanelPrototype { stack: true }
        .addTarget(commonPromTarget(
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
        .addTarget(commonPromTarget('node_memory_Buffers_bytes{%(nodeExporterSelector)s, instance="$instance"}' % config, legendFormat='memory buffers'))
        .addTarget(commonPromTarget('node_memory_Cached_bytes{%(nodeExporterSelector)s, instance="$instance"}' % config, legendFormat='memory cached'))
        .addTarget(commonPromTarget('node_memory_MemFree_bytes{%(nodeExporterSelector)s, instance="$instance"}' % config, legendFormat='memory free'))
      else if platform == 'Darwin' then
        // not useful to stack
        memoryGraphPanelPrototype { stack: false }
        .addTarget(commonPromTarget('node_memory_total_bytes{%(nodeExporterSelector)s, instance="$instance"}' % config, legendFormat='Physical Memory'))
        .addTarget(commonPromTarget(
          |||
            (
                node_memory_internal_bytes{%(nodeExporterSelector)s, instance="$instance"} -
                node_memory_purgeable_bytes{%(nodeExporterSelector)s, instance="$instance"} +
                node_memory_wired_bytes{%(nodeExporterSelector)s, instance="$instance"} +
                node_memory_compressed_bytes{%(nodeExporterSelector)s, instance="$instance"}
            )
          ||| % config, legendFormat='Memory Used'
        ))
        .addTarget(commonPromTarget(
          |||
            (
                node_memory_internal_bytes{%(nodeExporterSelector)s, instance="$instance"} -
                node_memory_purgeable_bytes{%(nodeExporterSelector)s, instance="$instance"}
            )
          ||| % config, legendFormat='App Memory'
        ))
        .addTarget(commonPromTarget('node_memory_wired_bytes{%(nodeExporterSelector)s, instance="$instance"}' % config, legendFormat='Wired Memory'))
        .addTarget(commonPromTarget('node_memory_compressed_bytes{%(nodeExporterSelector)s, instance="$instance"}' % config, legendFormat='Compressed')),

    // NOTE: avg() is used to circumvent a label change caused by a node_exporter rollout.
    local memoryGaugePanelPrototype =
      commonPanels.percentUsageStat.new('Memory Usage'),

    local memoryGauge =
      if platform == 'Linux' then
        memoryGaugePanelPrototype

        .addTarget(commonPromTarget(q.memoryUsage))

      else if platform == 'Darwin' then
        memoryGaugePanelPrototype
        .addTarget(commonPromTarget(
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
      nodePanels.timeseries.new('Disk I/O')
      .withFillOpacity(0)
      .withMin(0)
      // TODO: Does it make sense to have those three in the same panel?
      .addTarget(commonPromTarget(
        'rate(node_disk_read_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[$__rate_interval])' % config,
        legendFormat='{{device}} read',
      ))
      .addTarget(commonPromTarget(
        'rate(node_disk_written_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[$__rate_interval])' % config,
        legendFormat='{{device}} written',
      ))
      .addTarget(commonPromTarget(
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
      nodePanels.table.new(
        title='Disk Space Usage'
      )
      .setFieldConfig(unit='decbytes')
      //.addThresholdStep(color='light-green', value=null)
      .addThresholdStep(color='light-blue', value=null)
      .addThresholdStep(color='light-yellow', value=0.8)
      .addThresholdStep(color='light-red', value=0.9)
      .addTarget(commonPromTarget(
        |||
          max by (mountpoint) (node_filesystem_size_bytes{%(nodeExporterSelector)s, instance="$instance", %(fsSelector)s, %(fsMountpointSelector)s})
        ||| % config,
        legendFormat='',
        instant=true,
        format='table'
      ))
      .addTarget(commonPromTarget(
        |||
          max by (mountpoint) (node_filesystem_avail_bytes{%(nodeExporterSelector)s, instance="$instance", %(fsSelector)s, %(fsMountpointSelector)s})
        ||| % config,
        legendFormat='',
        instant=true,
        format='table',
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
          // {
          //   id: 'custom.displayMode',
          //   value: 'gradient-gauge',
          // },
          {
            "id": "custom.displayMode",
            "value": "basic"
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
      .sortBy('Mounted on')
      + {
        transformations+: [
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
          }
          
          // {
          //   id: 'sortBy',
          //   options: {
          //     fields: {},
          //     sort: [
          //       {
          //         field: 'Mounted on',
          //       },
          //     ],
          //   },
          // },
        ],
      },

    local networkTrafficPanel =
      commonPanels.networkTrafficGraph.new(
        'Network Traffic', description='Network transmitted and received (bits/s)',
      )
      .addTarget(commonPromTarget(
        expr=q.networkReceiveBitsPerSec,
        legendFormat='{{device}} received',
      ))
      .addTarget(commonPromTarget(
        expr=q.networkTransmitBitsPerSec,
        legendFormat='{{device}} transmitted',
      )),

    local networkErrorsDropsPanel =
      nodePanels.timeseries.new('Network Errors and Dropped Packets',)
      .addTarget(commonPromTarget(
        expr=q.networkReceiveErrorsPerSec,
        legendFormat='{{device}} errors in',
      ))
      .addTarget(commonPromTarget(
        expr=q.networkTransmitErrorsPerSec,
        legendFormat='{{device}} errors out',
      ))
      .addTarget(commonPromTarget(
        expr=q.networkReceiveDropsPerSec,
        legendFormat='{{device}} drop in',
      ))
      .addTarget(commonPromTarget(
        expr=q.networkTransmitDropsPerSec,
        legendFormat='{{device}} drop out',
      ))
      .withDecimals(1)
      .withUnits('pps')
      .withNegativeYByRegex(' out')
      .withAxisLabel('out(-) / in(+)'),


    local infoRow =
      row.new('Overview')
      .addPanel(uptimePanel { span: 3, height: '100px' } )
      .addPanel(nodeNamePanel { span: 3, height: '100px' })
      .addPanel(kernelVersionPanel { span: 3, height: '100px' })
      .addPanel(osPanel { span: 3, height: '100px' })
      .addPanel(cpuCountPanel { span: 3, height: '100px' })
      .addPanel(memoryTotalPanel { span: 3, height: '100px' })
      .addPanel(totalSwapPanel { span: 3, height: '100px' })
      .addPanel(totalRootFSPanel { span: 3, height: '100px' }),

    local cpuRow =
      row.new('CPU')
      .addPanel(idleCPU { span: 6 })
      .addPanel(systemLoad { span: 3 })
      .addPanel(cpuStatPanel { span: 3 }),

    local memoryRow =
      row.new('Memory')
      .addPanel(memoryGraph { span: 9 })
      .addPanel(memoryGauge { span: 3 }),

    local diskRow =
      row.new('Disk')
      .addPanel(diskIO { span: 6 })
      .addPanel(diskSpaceUsage { span: 6 }),

    local networkRow =
      row.new('Network')
      .addPanel(networkTrafficPanel { span: 6 })
      .addPanel(networkErrorsDropsPanel { span: 6 }),

    local rows =
      [
        infoRow,
        cpuRow,
        memoryRow,
        diskRow,
        networkRow,
      ],

    dashboard: if platform == 'Linux' then
      dashboard.new(
        '%sNode Overview ' % config.dashboardNamePrefix,
        time_from=config.dashboardInterval,
        tags=(config.dashboardTags),
        timezone=config.dashboardTimezone,
        refresh=config.dashboardRefresh,
        graphTooltip='shared_crosshair',
        uid='nodes'
      ) { editable: true }
      .addLink(c.links.fleetDash)
      .addLink(c.links.otherDashes)
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
