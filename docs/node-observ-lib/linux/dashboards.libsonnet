local g = import '../g.libsonnet';
local logslib = import 'github.com/grafana/jsonnet-libs/logs-lib/logs/main.libsonnet';
{
  local root = self,
  new(this):
    local prefix = this.config.dashboardNamePrefix;
    local links = this.grafana.links;
    local tags = this.config.dashboardTags;
    local uid = g.util.string.slugify(this.config.uid);
    local vars = this.grafana.variables;
    local annotations = this.grafana.annotations;
    local refresh = this.config.dashboardRefresh;
    local period = this.config.dashboardPeriod;
    local timezone = this.config.dashboardTimezone;
    local panels = this.grafana.panels;
    local rows = this.grafana.rows;
    local stat = g.panel.stat;
    {
      'fleet.json':
        local title = prefix + 'fleet overview';
        g.dashboard.new(title)
        + g.dashboard.withPanels(
          g.util.grid.wrapPanels(rows.fleet.panels, 12, 7)
        )
        // hide link to self
        + root.applyCommon(vars.multiInstance, uid + '-fleet', tags, links { backToFleet+:: {}, backToOverview+:: {} }, annotations, timezone, refresh, period),
      'overview.json':
        g.dashboard.new(prefix + 'overview')
        + g.dashboard.withPanels(
          g.util.panel.resolveCollapsedFlagOnRows(
            g.util.grid.wrapPanels(
              [
                rows.overview,
                rows.cpuOverview,
                rows.memoryOverview,
                rows.diskOverview,
                rows.networkOverview,
              ]
              +
              if this.config.enableHardware then
                [
                  rows.hardware,
                ] else []
              , 6, 2
            )
          )
        )
        // defaults to uid=nodes for backward compatibility with old node-mixins
        + root.applyCommon(vars.singleInstance, (if uid == 'node' then 'nodes' else uid + '-overview'), tags, links { backToOverview+:: {} }, annotations, timezone, refresh, period),
      'network.json':
        g.dashboard.new(prefix + 'network')
        + g.dashboard.withPanels(
          g.util.panel.resolveCollapsedFlagOnRows(
            g.util.grid.wrapPanels(
              [
                rows.network,
                rows.networkSockets
                + g.panel.row.withCollapsed(true),
                rows.networkNetstat
                + g.panel.row.withCollapsed(true),
              ], 12, 8
            )
          )
        )
        + root.applyCommon(vars.singleInstance, uid + '-network', tags, links, annotations, timezone, refresh, period),
      'memory.json':
        g.dashboard.new(prefix + 'memory')
        + g.dashboard.withPanels(
          g.util.panel.resolveCollapsedFlagOnRows(
            g.util.grid.wrapPanels(
              [
                rows.memoryOverview,
                rows.memoryVmstat,
                rows.memoryMemstat,
              ], 12, 8
            )
          )
        )
        + root.applyCommon(vars.singleInstance, uid + '-memory', tags, links, annotations, timezone, refresh, period),

      'system.json':
        g.dashboard.new(prefix + 'CPU and system')
        + g.dashboard.withPanels(
          g.util.panel.resolveCollapsedFlagOnRows(
            g.util.grid.wrapPanels(
              [
                rows.cpuAndSystem,
                rows.time,
              ], 12, 7
            )
          )
        )
        + root.applyCommon(vars.singleInstance, uid + '-system', tags, links, annotations, timezone, refresh, period),

      'disks.json':
        g.dashboard.new(prefix + 'filesystem and disks')
        + g.dashboard.withPanels(
          g.util.panel.resolveCollapsedFlagOnRows(
            g.util.grid.wrapPanels(
              [
                rows.filesystem,
                rows.disk,
              ], 12, 8
            )
          )
        )
        + root.applyCommon(vars.singleInstance, uid + '-disk', tags, links, annotations, timezone, refresh, period),
    }
    +
    if this.config.enableLokiLogs
    then
      {
        'logs.json':
          logslib.new(
            prefix + 'logs',
            datasourceName=this.grafana.variables.datasources.loki.name,
            datasourceRegex=this.grafana.variables.datasources.loki.regex,
            filterSelector=this.config.logsFilteringSelector,
            labels=this.config.groupLabels + this.config.instanceLabels + this.config.extraLogLabels,
            formatParser=null,
            showLogsVolume=this.config.showLogsVolume,
            logsVolumeGroupBy=this.config.logsVolumeGroupBy,
            extraFilters=this.config.logsExtraFilters
          )
          {
            dashboards+:
              {
                logs+:
                  // reference to self, already generated variables, to keep them, but apply other common data in applyCommon
                  root.applyCommon(super.logs.templating.list, uid=uid + '-logs', tags=tags, links=links, annotations=annotations, timezone=timezone, refresh=refresh, period=period),
              },
            panels+:
              {
                // modify log panel
                logs+:
                  g.panel.logs.options.withEnableLogDetails(true)
                  + g.panel.logs.options.withShowTime(false)
                  + g.panel.logs.options.withWrapLogMessage(false),
              },
            variables+: {
              // add prometheus datasource for annotations processing
              toArray+: [
                this.grafana.variables.datasources.prometheus { hide: 2 },
              ],
            },
          }.dashboards.logs,
      }
    else {},
  applyCommon(vars, uid, tags, links, annotations, timezone, refresh, period):
    g.dashboard.withTags(tags)
    + g.dashboard.withUid(uid)
    + g.dashboard.withLinks(std.objectValues(links))
    + g.dashboard.withTimezone(timezone)
    + g.dashboard.withRefresh(refresh)
    + g.dashboard.time.withFrom(period)
    + g.dashboard.withVariables(vars)
    + g.dashboard.withAnnotations(std.objectValues(annotations)),
}
