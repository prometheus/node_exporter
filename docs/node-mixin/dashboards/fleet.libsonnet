local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local nodePanels = import '../lib/panels/panels.libsonnet';
local commonPanels = import '../lib/panels/common/panels.libsonnet';
local nodeTimeseries = nodePanels.timeseries;
local common = import '../lib/common.libsonnet';

{

  new(config=null, platform=null):: {
    local c = common.new(config=config, platform=platform),
    local commonPromTarget = c.commonPromTarget,


    local templates = [
      if std.member(std.split(config.instanceLabels, ','), template.name)
      then
        template
        {
          allValue: '.+',
          includeAll: true,
          multi: true,
        }
      else template
      for template in c.templates
    ],

    local q = c.queries,

    local fleetTable =
      nodePanels.table.new(
        title='Linux Nodes Overview'
      )
      .addTarget(commonPromTarget(expr=q.osInfo, format='table', instant=true) { refId: 'INFO' })
      .addTarget(commonPromTarget(expr=q.nodeInfo, format='table', instant=true) { refId: 'OS' })
      .addTarget(commonPromTarget(expr=q.uptime, format='table', instant=true) { refId: 'UPTIME' })
      .addTarget(commonPromTarget(expr=q.systemLoad1, format='table', instant=true) { refId: 'LOAD1' })
      .addTarget(commonPromTarget(expr=q.systemLoad5, format='table', instant=true) { refId: 'LOAD5' })
      .addTarget(commonPromTarget(expr=q.systemLoad15, format='table', instant=true) { refId: 'LOAD15' })
      .addTarget(commonPromTarget(
        expr=q.cpuCount,
        format='table',
        instant=true,
      ) { refId: 'CPUCOUNT' })
      .addTarget(commonPromTarget(
        expr=q.cpuUsage, format='table', instant=true,
      ) { refId: 'CPUUSAGE' })
      .addTarget(commonPromTarget(expr=q.memoryTotal, format='table', instant=true) { refId: 'MEMTOTAL' })
      .addTarget(commonPromTarget(expr=q.memoryUsage, format='table', instant=true) { refId: 'MEMUSAGE' })
      .addTarget(commonPromTarget(expr=q.fsSizeTotalRoot, format='table', instant=true) { refId: 'FSTOTAL' })
      .addTarget(commonPromTarget(
        expr=
        |||
          100-(max by (%(instanceLabels)s) (node_filesystem_avail_bytes{%(nodeQuerySelector)s, fstype!="", mountpoint="/"})
          /
          max by (%(instanceLabels)s) (node_filesystem_size_bytes{%(nodeQuerySelector)s, fstype!="", mountpoint="/"}) * 100)
        ||| % config { nodeQuerySelector: c.nodeQuerySelector },
        format='table',
        instant=true,
      ) { refId: 'FSUSAGE' })
      .addTarget(commonPromTarget(
        expr='count by (%(instanceLabels)s) (max_over_time(ALERTS{%(nodeQuerySelector)s, alertstate="firing", severity="critical"}[1m])) * group by (%(instanceLabels)s) (node_uname_info{})' % config { nodeQuerySelector: c.nodeQuerySelector },
        format='table',
        instant=true
      ) { refId: 'CRITICAL' })
      .addTarget(commonPromTarget(
        expr='count by (%(instanceLabels)s) (max_over_time(ALERTS{%(nodeQuerySelector)s, alertstate="firing", severity="warning"}[1m])) * group by (%(instanceLabels)s) (node_uname_info{})' % config { nodeQuerySelector: c.nodeQuerySelector },
        format='table',
        instant=true
      ) { refId: 'WARNING' })
      .withTransform()
      .joinByField(field=std.split(config.instanceLabels, ',')[0])
      //disable kernel and os:
      //.filterFieldsByName('instance|pretty_name|nodename|release|Value.+')
      .filterFieldsByName(std.split(config.instanceLabels, ',')[0] + '|nodename|Value.+')
      .organize(
        excludeByName={
          'Value #OS': true,
          'Value #INFO': true,
          'Value #LOAD5': true,
          'Value #LOAD15': true,
        },
        renameByName={
          instance: 'Instance',
          pretty_name: 'OS',
          nodename: 'Hostname',
          release: 'Kernel version',
          'Value #LOAD1': 'Load 1m',
          'Value #LOAD5': 'Load 5m',
          'Value #LOAD15': 'Load 15m',
          'Value #CPUCOUNT': 'Cores',
          'Value #CPUUSAGE': 'CPU usage',
          'Value #MEMTOTAL': 'Memory total',
          'Value #MEMUSAGE': 'Memory usage',
          'Value #FSTOTAL': 'Root disk size',
          'Value #FSUSAGE': 'Root disk usage',
          'Value #UPTIME': 'Uptime',
          'Value #CRITICAL': 'Crit Alerts',
          'Value #WARNING': 'Warnings',
        }
      )
      .withFooter(reducer=['mean'], fields=[
        'Value #LOAD1',
        'Value #MEMUSAGE',
        'Value #CPUUSAGE',
      ])
      //.addThresholdStep(color='light-green', value=null)
      .addThresholdStep(color='light-blue', value=null)
      .addThresholdStep(color='light-yellow', value=80)
      .addThresholdStep(color='light-red', value=90)
      .addOverride(
        matcher={
          id: 'byName',
          options: 'Instance',
        },
        properties=[
          {
            id: 'links',
            value: [
              {
                targetBlank: true,
                title: c.links.instanceDataLinkForTable.title,
                url: c.links.instanceDataLinkForTable.url,
              },
            ],
          },
          {
            id: 'custom.filterable',
            value: true,
          },
        ]
      )
      .addOverride(
        matcher={
          id: 'byRegexp',
          options: 'OS|Kernel version|Hostname',
        },
        properties=[
          {
            id: 'custom.filterable',
            value: true,
          },
        ]
      )
      .addOverride(
        matcher={
          id: 'byRegexp',
          options: 'Memory total|Root disk size',
        },
        properties=[
          {
            id: 'unit',
            value: 'bytes',
          },
          {
            id: 'decimals',
            value: 0,
          },
        ]
      )
      .addOverride(
        matcher={
          id: 'byName',
          options: 'Cores',
        },
        properties=[
          {
            id: 'custom.width',
            value: 60,
          },
        ]
      )
      .addOverride(
        matcher={
          id: 'byRegexp',
          options: 'Load.+',
        },
        properties=[
          {
            id: 'custom.width',
            value: 60,
          },
        ]
      )
      .addOverride(
        matcher={
          id: 'byName',
          options: 'Uptime',
        },
        properties=[
          {
            id: 'unit',
            value: 'dtdurations',
          },
          {
            id: 'custom.displayMode',
            value: 'color-text',
          },
          {
            id: 'thresholds',
            value: {
              mode: 'absolute',
              steps: [
                {
                  color: 'light-orange',
                  value: null,
                },
                {
                  color: 'text',
                  value: 300,
                },
              ],
            },
          },
        ]
      )
      .addOverride(
        matcher={
          id: 'byRegexp',
          options: 'CPU usage|Memory usage|Root disk usage',
        },
        properties=[
          {
            id: 'unit',
            value: 'percent',
          },
          // {
          //   id: 'custom.displayMode',
          //   value: 'gradient-gauge',
          // },
          {
            id: 'custom.displayMode',
            value: 'basic',
          },
          {
            id: 'max',
            value: 100,
          },
          {
            id: 'min',
            value: 0,
          },
        ]
      )
      .sortBy('Instance')
    ,

    local memoryUsagePanel =
      nodePanels.timeseries.new('Memory Usage', description='Top 25')
      .withUnits('percent')
      .withMin(0)
      .withMax(100)
      .withColor(mode='continuous-BlYlRd')
      .withFillOpacity(1)
      .withGradientMode('scheme')
      .withLegend(mode='table', calcs=['mean', 'max', 'lastNotNull'], placement='right')
      .addDataLink(
        title=c.links.instanceDataLink.title,
        url=c.links.instanceDataLink.url,
      )
      .addTarget(commonPromTarget(
        expr='topk(25, ' + q.memoryUsage + ')',
        legendFormat=c.labelsToLegend(std.split(config.instanceLabels, ','))
      ))
      .addTarget(commonPromTarget(
        expr='avg(' + q.memoryUsage + ')',
        legendFormat='Mean',
      ))
      .addOverride(
        matcher={
          id: 'byName',
          options: 'Mean',

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
          {
            id: 'custom.fillOpacity',
            value: 0,
          },
          {
            id: 'color',
            value: {
              mode: 'fixed',
              fixedColor: 'light-purple',
            },
          },
          {
            id: 'custom.lineWidth',
            value: 2,
          },
        ]
      ),

    local cpuUsagePanel =
      nodePanels.timeseries.new('CPU Usage', description='Top 25')
      .withUnits('percent')
      .withMin(0)
      .withMax(100)
      .withFillOpacity(1)
      .withColor(mode='continuous-BlYlRd')
      .withGradientMode('scheme')
      .withLegend(mode='table', calcs=['mean', 'max', 'lastNotNull'], placement='right')
      .addDataLink(
        title=c.links.instanceDataLink.title,
        url=c.links.instanceDataLink.url,
      )
      .addTarget(commonPromTarget(
        expr='topk(25, ' + q.cpuUsage + ')',
        legendFormat=c.labelsToLegend(std.split(config.instanceLabels, ',')),
      ))
      .addTarget(commonPromTarget(
        expr='avg(' + q.cpuUsage + ')',
        legendFormat='Mean',
      ))
      .addOverride(
        matcher={
          id: 'byName',
          options: 'Mean',

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
          {
            id: 'custom.fillOpacity',
            value: 0,
          },
          {
            id: 'color',
            value: {
              mode: 'fixed',
              fixedColor: 'light-purple',
            },
          },
          {
            id: 'custom.lineWidth',
            value: 2,
          },
        ]
      ),

    local diskIOPanel =
      nodePanels.timeseries.new('Disks I/O', description='Top 25')
      .withUnits('percentunit')
      .withMin(0)
      .withMax(1)
      .withFillOpacity(1)
      .withColor(mode='continuous-BlYlRd')
      .withGradientMode('scheme')
      .withLegend(mode='table', calcs=['mean', 'max', 'lastNotNull'], placement='right')
      .addDataLink(
        title=c.links.instanceDataLink.title,
        url=c.links.instanceDataLink.url,
      )
      .addTarget(commonPromTarget(
        expr='topk(25, ' + q.diskIoTime + ')',
        legendFormat=c.labelsToLegend(std.split(config.instanceLabels, ',')) + ': {{device}}',
      ))
      .addOverride(
        matcher={
          id: 'byName',
          options: 'Mean',

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
          {
            id: 'custom.fillOpacity',
            value: 0,
          },
          {
            id: 'color',
            value: {
              mode: 'fixed',
              fixedColor: 'light-purple',
            },
          },
          {
            id: 'custom.lineWidth',
            value: 2,
          },
        ]
      ),
    local diskSpacePanel =
      nodePanels.timeseries.new('Disks Space Usage', description='Top 25')
      .withUnits('percentunit')
      .withMin(0)
      .withMax(1)
      .withFillOpacity(1)
      .withColor(mode='continuous-BlYlRd')
      .withGradientMode('scheme')
      .withLegend(mode='table', calcs=['mean', 'max', 'lastNotNull'], placement='right')
      .addDataLink(
        title=c.links.instanceDataLink.title,
        url=c.links.instanceDataLink.url,
      )
      .addTarget(commonPromTarget(
        expr='topk(25, ' + q.diskSpaceUsage + ')',
        legendFormat=c.labelsToLegend(std.split(config.instanceLabels, ',')) + ': {{mountpoint}}',
      ))
      .addOverride(
        matcher={
          id: 'byName',
          options: 'Mean',

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
          {
            id: 'custom.fillOpacity',
            value: 0,
          },
          {
            id: 'color',
            value: {
              mode: 'fixed',
              fixedColor: 'light-purple',
            },
          },
          {
            id: 'custom.lineWidth',
            value: 2,
          },
        ]
      ),
    local networkErrorsDropsPanel =
      nodePanels.timeseries.new('Network Errors and Dropped Packets', description='Top 25')
      .withLegend(mode='table', calcs=['mean', 'max', 'lastNotNull'], placement='right')
      .addTarget(commonPromTarget(
        expr='topk(25, ' + q.networkReceiveErrorsPerSec + ' + ' + q.networkTransmitErrorsPerSec + ' + ' + q.networkReceiveDropsPerSec + ' + ' + q.networkTransmitDropsPerSec + ') > 0.5',
        legendFormat=c.labelsToLegend(std.split(config.instanceLabels, ',')) + ': {{device}}',
      ))
      .withDecimals(1)
      .withUnits('pps')
      .withDrawStyle('points')
      .withPointsSize(5)
      .addDataLink(
        title=c.links.instanceDataLink.title,
        url=c.links.instanceDataLink.url,
      ),

    local rows =
      [
        row.new('Overview')
        .addPanel(fleetTable { span: 12, height: '800px' })
        .addPanel(cpuUsagePanel { span: 12 })
        .addPanel(memoryUsagePanel { span: 12 })
        .addPanel(diskIOPanel { span: 6 }).addPanel(diskSpacePanel { span: 6 })
        .addPanel(networkErrorsDropsPanel { span: 12 }),
      ],

    dashboard: if platform == 'Linux' then
      dashboard.new(
        '%sNode Fleet Overview' % config.dashboardNamePrefix,
        time_from=config.dashboardInterval,
        tags=(config.dashboardTags),
        timezone=config.dashboardTimezone,
        refresh=config.dashboardRefresh,
        graphTooltip='shared_crosshair',
        uid='node-fleet'
      ) { editable: true }
      .addLink(c.links.otherDashes { includeVars: false })
      .addTemplates(templates)
      .addRows(rows)
    else if platform == 'Darwin' then {},
  },
}
