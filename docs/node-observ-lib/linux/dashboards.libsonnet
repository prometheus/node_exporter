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
    local stat = g.panel.stat;
    {
      fleet:
        local title = prefix + 'fleet overview';
        g.dashboard.new(title)
        + g.dashboard.withPanels(
          g.util.grid.wrapPanels(
            [
              // g.panel.row.new("Overview"),
              panels.fleet.fleetOverviewTable { gridPos+: { w: 24, h: 16 } },
              panels.cpu.cpuUsageTopk { gridPos+: { w: 24 } },
              panels.memory.memotyUsageTopKPercent { gridPos+: { w: 24 } },
              panels.disk.diskIOutilPercentTopK { gridPos+: { w: 12 } },
              panels.disk.diskUsagePercentTopK { gridPos+: { w: 12 } },
              panels.network.networkErrorsAndDroppedPerSecTopK { gridPos+: { w: 24 } },
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
              panels.system.uptime,
              panels.system.hostname,
              panels.system.kernelVersion,
              panels.system.osInfo,
              panels.cpu.cpuCount,
              panels.memory.memoryTotalBytes,
              panels.memory.memorySwapTotalBytes,
              panels.disk.diskTotalRoot,
              g.panel.row.new('CPU'),
              panels.cpu.cpuUsageStat { gridPos+: { w: 6, h: 6 } },
              panels.cpu.cpuUsageTsPerCore { gridPos+: { w: 12, h: 6 } },
              panels.system.systemLoad { gridPos+: { w: 6, h: 6 } },
              g.panel.row.new('Memory'),
              panels.memory.memoryUsageStatPercent { gridPos+: { w: 6, h: 6 } },
              panels.memory.memoryUsageTsBytes { gridPos+: { w: 18, h: 6 } },
              g.panel.row.new('Disk'),
              panels.disk.diskIOBytesPerSec { gridPos+: { w: 12, h: 8 } },
              panels.disk.diskUsage { gridPos+: { w: 12, h: 8 } },
              g.panel.row.new('Network'),
              panels.network.networkUsagePerSec { gridPos+: { w: 12, h: 8 } },
              panels.network.networkErrorsAndDroppedPerSec { gridPos+: { w: 12, h: 8 } },
            ]
            +
            if this.config.enableHardware then
              [
                g.panel.row.new('Hardware'),
                panels.hardware.hardwareTemperature { gridPos+: { w: 12, h: 8 } },
              ] else []
            , 6, 2
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
              panels.network.networkOverviewTable { gridPos: { w: 24 } },
              panels.network.networkUsagePerSec,
              panels.network.networkOperStatus,
              panels.network.networkErrorsPerSec,
              panels.network.networkDroppedPerSec,
              panels.network.networkPacketsPerSec,
              panels.network.networkMulticastPerSec,
              panels.network.networkFifo,
              panels.network.networkCompressedPerSec,
              panels.network.networkNFConntrack,
              panels.network.networkSoftnet,
              panels.network.networkSoftnetSqueeze,
              g.panel.row.new('Network sockets'),
              panels.network.networkSockstatAll { gridPos: { w: 24 } },
              panels.network.networkSockstatTCP,
              panels.network.networkSockstatUDP,
              panels.network.networkSockstatMemory,
              panels.network.networkSockstatOther,
              g.panel.row.new('Network netstat'),
              panels.network.networkNetstatIP { gridPos: { w: 24 } },
              panels.network.networkNetstatTCP,
              panels.network.networkNetstatTCPerrors,
              panels.network.networkNetstatUDP,
              panels.network.networkNetstatUDPerrors,
              panels.network.networkNetstatICMP,
              panels.network.networkNetstatICMPerrors,
            ], 12, 8
          )
        )
        + root.applyCommon(vars.singleInstance, uid + '-network', tags, links, annotations, timezone, refresh, period),
      memory:
        g.dashboard.new(prefix + 'memory')
        + g.dashboard.withPanels(
          g.util.grid.wrapPanels(
            [
              panels.memory.memoryUsageStatPercent { gridPos+: { w: 6, h: 6 } },
              panels.memory.memoryUsageTsBytes { gridPos+: { w: 18, h: 6 } },
              g.panel.row.new('Vmstat'),
              panels.memory.memoryPagesInOut,
              panels.memory.memoryPagesSwapInOut,
              panels.memory.memoryPagesFaults,
              panels.memory.memoryOOMkiller,
              g.panel.row.new('Memstat'),
              panels.memory.memoryActiveInactive,
              panels.memory.memoryActiveInactiveDetail,
              panels.memory.memoryCommited,
              panels.memory.memorySharedAndMapped,
              panels.memory.memoryWriteAndDirty,
              panels.memory.memoryVmalloc,
              panels.memory.memorySlab,
              panels.memory.memoryAnonymous,
              panels.memory.memoryHugePagesCounter,
              panels.memory.memoryHugePagesSize,
              panels.memory.memoryDirectMap,
              panels.memory.memoryBounce,
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
              panels.cpu.cpuUsageStat { gridPos+: { w: 6, h: 6 } },
              panels.cpu.cpuUsageTsPerCore { gridPos+: { w: 9, h: 6 } },
              panels.cpu.cpuUsageByMode { gridPos+: { w: 9, h: 6 } },
              panels.system.systemLoad,
              panels.system.systemContextSwitchesAndInterrupts,
              g.panel.row.new('Time'),
              panels.system.osTimezone { gridPos+: { w: 3, h: 4 } },
              panels.system.timeNtpStatus { gridPos+: { x: 0, y: 0, w: 21, h: 4 } },
              panels.system.timeSyncDrift { gridPos+: { w: 24, h: 7 } },
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
              panels.disk.diskFreeTs,
              panels.disk.diskUsage,
              panels.disk.diskInodesFree,
              panels.disk.diskInodesTotal,
              panels.disk.diskErrorsandRO,
              panels.disk.fileDescriptors,
              g.panel.row.new('Disk'),
              panels.disk.diskIOBytesPerSec,
              panels.disk.diskIOps,
              panels.disk.diskIOWaitTime,
              panels.disk.diskQueue,
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
