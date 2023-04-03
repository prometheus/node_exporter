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
    local templates = c.templates,
    local q = c.queries,

    local cpuUsageModes =
      nodeTimeseries.new('CPU Usage')
      .withStacking('normal')
      .withUnits('percent')
      .withFillOpacity(100)
      .withMax(100)
      .withMin(0)
      .addTarget(commonPromTarget(
        expr=q.cpuUsageModes,
        legendFormat='{{mode}}',
      )),

    local timeSyncDrift =
      nodeTimeseries.new('Time Syncronized Drift')
      .withUnits('seconds')
      .addTarget(commonPromTarget(
        expr=q.node_timex_estimated_error_seconds,
        legendFormat='Estimated error in seconds',
      ))
      .addTarget(commonPromTarget(
        expr=q.node_timex_offset_seconds,
        legendFormat='Time offset in between local system and reference clock',
      ))
      .addTarget(commonPromTarget(
        expr=q.node_timex_maxerror_seconds,
        legendFormat='Maximum error in seconds'
      )),

    local timeSyncronizedStatus =
      nodeTimeseries.new('Time Syncronized Status')
      .withColor(mode='palette-classic')
      .withFillOpacity(75)
      .withLegend(show=false)
      {
        maxDataPoints: 100,
        type: 'status-history',
        fieldConfig+: {
          defaults+: {
            mappings+: [
              {
                type: 'value',
                options: {
                  '1': {
                    text: 'In sync',
                    color: 'light-green',
                    index: 1,
                  },
                },
              },
              {
                type: 'value',
                options: {
                  '0': {
                    text: 'Not in sync',
                    color: 'light-yellow',
                    index: 0,
                  },
                },
              },

            ],
          },
        },
      }
      .addTarget(commonPromTarget(
        expr=q.node_timex_sync_status,
        legendFormat='Sync status',
      )),

    local systemdStates =
      nodeTimeseries.new('Systemd Units States')
      .addTarget(commonPromTarget(
        expr=q.node_systemd_units,
        legendFormat='{{ state }}'
      )),

    local panelsGrid =
      [
        //use negative gravity(skip y), max w=24, default h should be '6'.
        c.panelsWithTargets.cpuStatPanel { gridPos: { x: 0, w: 6, h: 6 } },
        c.panelsWithTargets.idleCPU { gridPos: { x: 6, h: 6, w: 9 } },
        cpuUsageModes { gridPos: { x: 15, h: 6, w: 9 } },
        //pseudorow y:25
        c.panelsWithTargets.systemLoad { gridPos: { x: 0, h: 6, w: 12, y: 25 } },
        c.panelsWithTargets.systemContextSwitches { gridPos: { x: 12, h: 6, w: 12, y: 25 } },
        { type: 'row', title: 'Systemd', gridPos: { x: 0, w: 24, y: 50 } },
        systemdStates { gridPos: { x: 0, h: 6, w: 12, y: 50 } },
        { type: 'row', title: 'Time', gridPos: { x: 0, w: 24, y: 75 } },
        timeSyncronizedStatus { gridPos: { x: 0, h: 3, w: 24, y: 75 } },
        timeSyncDrift { gridPos: { x: 0, h: 6, w: 24, y: 75 } },
      ],

    dashboard: if platform == 'Linux' then
      dashboard.new(
        '%sNode CPU and System' % config { nodeQuerySelector: c.nodeQuerySelector }.dashboardNamePrefix,
        time_from=config.dashboardInterval,
        tags=(config.dashboardTags),
        timezone=config.dashboardTimezone,
        refresh=config.dashboardRefresh,
        graphTooltip='shared_crosshair',
        uid=config.grafanaDashboardIDs['nodes-system.json'],
      )
      .addLink(c.links.fleetDash)
      .addLink(c.links.nodeDash)
      .addLink(c.links.otherDashes)
      .addAnnotations(c.annotations)
      .addTemplates(templates)
      .addPanels(panelsGrid)
    else if platform == 'Darwin' then {},
  },
}
