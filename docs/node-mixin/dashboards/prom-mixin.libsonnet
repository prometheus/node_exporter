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


    local panelsGrid =
      [
        // use negative gravity effect, max w=24, default h=8
        { type: 'row', title: 'Overview' },
        uptimePanel { gridPos: { x: 0, w: 6, h: 2 } },
        nodeNamePanel { gridPos: { x: 6, w: 6, h: 2 } },
        kernelVersionPanel { gridPos: { x: 12, w: 6, h: 2 } },
        osPanel { gridPos: { x: 18, w: 6, h: 2 } },
        cpuCountPanel { gridPos: { x: 0, w: 6, h: 2 } },
        memoryTotalPanel { gridPos: { x: 6, w: 6, h: 2 } },
        totalSwapPanel { gridPos: { x: 12, w: 6, h: 2 } },
        totalRootFSPanel { gridPos: { x: 18, w: 6, h: 2 } },
        { type: 'row', title: 'CPU' } { gridPos: { y: 25 } },
        c.panelsWithTargets.cpuStatPanel { gridPos: { x: 0, w: 6, h: 6, y: 25 } },
        c.panelsWithTargets.idleCPU { gridPos: { x: 6, w: 12, h: 6, y: 25 } },
        c.panelsWithTargets.systemLoad { gridPos: { x: 18, w: 6, h: 6, y: 25 } },
        { type: 'row', title: 'Memory' } { gridPos: { y: 50 } },
        c.panelsWithTargets.memoryGauge { gridPos: { x: 0, w: 6, h: 6, y: 50 } },
        c.panelsWithTargets.memoryGraph { gridPos: { x: 6, w: 18, h: 6, y: 50 } },
        { type: 'row', title: 'Disk' } { gridPos: { y: 75 } },
        c.panelsWithTargets.diskIO { gridPos: { x: 0, w: 12, h: 8, y: 75 } },
        c.panelsWithTargets.fsSpaceUsage { gridPos: { x: 12, w: 12, h: 8, y: 75 } },
        { type: 'row', title: 'Network' } { gridPos: { y: 100 } },
        networkTrafficPanel { gridPos: { x: 0, w: 12, h: 8, y: 100 } },
        networkErrorsDropsPanel { gridPos: { x: 12, w: 12, h: 8, y: 100 } },
      ],
    dashboard: if platform == 'Linux' then
      dashboard.new(
        '%sNode Overview ' % config { nodeQuerySelector: c.nodeQuerySelector }.dashboardNamePrefix,
        time_from=config.dashboardInterval,
        tags=(config.dashboardTags),
        timezone=config.dashboardTimezone,
        refresh=config.dashboardRefresh,
        graphTooltip='shared_crosshair',
        uid='nodes'
      ) { editable: true }
      .addLink(c.links.fleetDash)
      .addLink(c.links.otherDashes)
      .addAnnotations(c.annotations)
      .addTemplates(templates)
      .addPanels(panelsGrid)
    else if platform == 'Darwin' then
      dashboard.new(
        '%sMacOS' % config { nodeQuerySelector: c.nodeQuerySelector }.dashboardNamePrefix,
        time_from=config.dashboardInterval,
        tags=(config.dashboardTags),
        timezone=config.dashboardTimezone,
        refresh=config.dashboardRefresh,
        graphTooltip='shared_crosshair',
        uid='nodes-darwin'
      )
      .addTemplates(templates)
      .addPanels(panelsGrid),

  },
}
