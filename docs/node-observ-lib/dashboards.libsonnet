local g = import './g.libsonnet';
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
    local stat = g.panel.stat;
    {
      fleet:
        local title = prefix + 'fleet overview';
        g.dashboard.new(title)
        + g.dashboard.withPanels(
          g.util.grid.wrapPanels(
            [
              // g.panel.row.new("Overview"),
              panels.fleetOverviewTable { gridPos+: { w: 24, h: 16 } },
              panels.cpuUsageTopk { gridPos+: { w: 24 } },
              panels.memotyUsageTopKPercent { gridPos+: { w: 24 } },
              panels.diskIOutilPercentTopK { gridPos+: { w: 12 } },
              panels.diskUsagePercentTopK { gridPos+: { w: 12 } },
              panels.networkErrorsAndDroppedPerSecTopK { gridPos+: { w: 24 } },
            ], 12, 7
          )
        )
        // hide link to self
        + root.applyCommon(vars.multiInstance, uid + '-fleet', tags, links { backToFleet+:: {}, backToOverview+:: {} }, annotations, timezone, refresh, period),
      overview:
        g.dashboard.new(prefix + 'overview')
        + g.dashboard.withPanels(
          g.util.grid.wrapPanels(
            [
              g.panel.row.new('Overview'),
              panels.uptime,
              panels.hostname,
              panels.kernelVersion,
              panels.osInfo,
              panels.cpuCount,
              panels.memoryTotalBytes,
              panels.memorySwapTotalBytes,
              panels.diskTotalRoot,
              g.panel.row.new('CPU'),
              panels.cpuUsageStat { gridPos+: { w: 6, h: 6 } },
              panels.cpuUsageTsPerCore { gridPos+: { w: 12, h: 6 } },
              panels.systemLoad { gridPos+: { w: 6, h: 6 } },
              g.panel.row.new('Memory'),
              panels.memoryUsageStatPercent { gridPos+: { w: 6, h: 6 } },
              panels.memoryUsageTsBytes { gridPos+: { w: 18, h: 6 } },
              g.panel.row.new('Disk'),
              panels.diskIOBytesPerSec { gridPos+: { w: 12, h: 8 } },
              panels.diskUsage { gridPos+: { w: 12, h: 8 } },
              g.panel.row.new('Network'),
              panels.networkUsagePerSec { gridPos+: { w: 12, h: 8 } },
              panels.networkErrorsAndDroppedPerSec { gridPos+: { w: 12, h: 8 } },
            ], 6, 2
          )
        )
        // defaults to uid=nodes for backward compatibility with old node-mixins
        + root.applyCommon(vars.singleInstance, (if uid == 'node' then 'nodes' else uid + '-overview'), tags, links { backToOverview+:: {} }, annotations, timezone, refresh, period),
      network:
        g.dashboard.new(prefix + 'network')
        + g.dashboard.withPanels(
          g.util.grid.wrapPanels(
            [
              g.panel.row.new('Network'),
              panels.networkOverviewTable { gridPos: { w: 24 } },
              panels.networkUsagePerSec,
              panels.networkOperStatus,
              panels.networkErrorsPerSec,
              panels.networkDroppedPerSec,
              panels.networkPacketsPerSec,
              panels.networkMulticastPerSec,
              panels.networkFifo,
              panels.networkCompressedPerSec,
              panels.networkNFConntrack,
              panels.networkSoftnet,
              panels.networkSoftnetSqueeze,
              g.panel.row.new('Network sockets'),
              panels.networkSockstatAll { gridPos: { w: 24 } },
              panels.networkSockstatTCP,
              panels.networkSockstatUDP,
              panels.networkSockstatMemory,
              panels.networkSockstatOther,
              g.panel.row.new('Network netstat'),
              panels.networkNetstatIP { gridPos: { w: 24 } },
              panels.networkNetstatTCP,
              panels.networkNetstatTCPerrors,
              panels.networkNetstatUDP,
              panels.networkNetstatUDPerrors,
              panels.networkNetstatICMP,
              panels.networkNetstatICMPerrors,
            ], 12, 8
          )
        )
        + root.applyCommon(vars.singleInstance, uid + '-network', tags, links, annotations, timezone, refresh, period),
      memory:
        g.dashboard.new(prefix + 'memory')
        + g.dashboard.withPanels(
          g.util.grid.wrapPanels(
            [
              panels.memoryUsageStatPercent { gridPos+: { w: 6, h: 6 } },
              panels.memoryUsageTsBytes { gridPos+: { w: 18, h: 6 } },
              g.panel.row.new('Vmstat'),
              panels.memoryPagesInOut,
              panels.memoryPagesSwapInOut,
              panels.memoryPagesFaults,
              panels.memoryOOMkiller,
              g.panel.row.new('Memstat'),
              panels.memoryActiveInactive,
              panels.memoryActiveInactiveDetail,
              panels.memoryCommited,
              panels.memorySharedAndMapped,
              panels.memoryWriteAndDirty,
              panels.memoryVmalloc,
              panels.memorySlab,
              panels.memoryAnonymous,
              panels.memoryHugePagesCounter,
              panels.memoryHugePagesSize,
              panels.memoryDirectMap,
              panels.memoryBounce,
            ], 12, 8
          )
        )
        + root.applyCommon(vars.singleInstance, uid + '-memory', tags, links, annotations, timezone, refresh, period),

      system:
        g.dashboard.new(prefix + 'CPU and system')
        + g.dashboard.withPanels(
          g.util.grid.wrapPanels(
            [
              g.panel.row.new('System'),
              panels.cpuUsageStat { gridPos+: { w: 6, h: 6 } },
              panels.cpuUsageTsPerCore { gridPos+: { w: 9, h: 6 } },
              panels.cpuUsageByMode { gridPos+: { w: 9, h: 6 } },
              panels.systemLoad,
              panels.systemContextSwitchesAndInterrupts,
              g.panel.row.new('Time'),
              panels.osTimezone { gridPos+: { w: 3, h: 4 } },
              panels.timeNtpStatus { gridPos+: { x: 0, y: 0, w: 21, h: 4 } },
              panels.timeSyncDrift { gridPos+: { w: 24, h: 7 } },
            ], 12, 7
          )
        )
        + root.applyCommon(vars.singleInstance, uid + '-system', tags, links, annotations, timezone, refresh, period),

      disks:
        g.dashboard.new(prefix + 'filesystem and disks')
        + g.dashboard.withPanels(
          g.util.grid.wrapPanels(
            [
              g.panel.row.new('Filesystem'),
              panels.diskFreeTs,
              panels.diskUsage,
              panels.diskInodesFree,
              panels.diskInodesTotal,
              panels.diskErrorsandRO,
              panels.fileDescriptors,
              g.panel.row.new('Disk'),
              panels.diskIOBytesPerSec,
              panels.diskIOps,
              panels.diskIOWaitTime,
              panels.diskQueue,
            ], 12, 8
          )
        )
        + root.applyCommon(vars.singleInstance, uid + '-disk', tags, links, annotations, timezone, refresh, period),
    }
    +
    if this.config.enableLokiLogs
    then
      {
        logs:
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
