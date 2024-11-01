local g = import '../../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
{
  new(this):
    {
      local panels = this.grafana.panels,

      fleet:
        g.panel.row.new('Fleet overview')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels(
          [
            // g.panel.row.new("Overview"),
            panels.fleet.fleetOverviewTable { gridPos+: { w: 24, h: 16 } },
            panels.cpu.cpuUsageTopk { gridPos+: { w: 24 } },
            panels.memory.memotyUsageTopKPercent { gridPos+: { w: 24 } },
            panels.disk.diskIOutilPercentTopK { gridPos+: { w: 12 } },
            panels.disk.diskUsagePercentTopK { gridPos+: { w: 12 } },
            panels.network.networkErrorsAndDroppedPerSecTopK { gridPos+: { w: 24 } },
          ]
        ),
      overview:
        g.panel.row.new('Overview')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels(
          [
            panels.system.uptime { gridPos+: { w: 6, h: 2 } },
            panels.system.hostname { gridPos+: { w: 6, h: 2 } },
            panels.system.kernelVersion { gridPos+: { w: 6, h: 2 } },
            panels.system.osInfo { gridPos+: { w: 6, h: 2 } },
            panels.cpu.cpuCount { gridPos+: { w: 6, h: 2 } },
            panels.memory.memoryTotalBytes { gridPos+: { w: 6, h: 2 } },
            panels.memory.memorySwapTotalBytes { gridPos+: { w: 6, h: 2 } },
            panels.disk.diskTotalRoot { gridPos+: { w: 6, h: 2 } },
          ]
        ),
      time:
        g.panel.row.new('Time')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels([
          panels.system.osTimezone { gridPos+: { w: 3, h: 4 } },
          panels.system.timeNtpStatus { gridPos+: { x: 0, y: 0, w: 21, h: 4 } },
          panels.system.timeSyncDrift { gridPos+: { w: 24, h: 7 } },
        ]),
      cpuOverview:
        g.panel.row.new('CPU')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels(
          [
            panels.cpu.cpuUsageStat { gridPos+: { w: 6, h: 6 } },
            panels.cpu.cpuUsageTsPerCore { gridPos+: { w: 12, h: 6 } },
            panels.system.systemLoad { gridPos+: { w: 6, h: 6 } },
          ]
        ),
      cpuAndSystem:
        g.panel.row.new('System')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels([
          panels.cpu.cpuUsageStat { gridPos+: { w: 6, h: 6 } },
          panels.cpu.cpuUsageTsPerCore { gridPos+: { w: 9, h: 6 } },
          panels.cpu.cpuUsageByMode { gridPos+: { w: 9, h: 6 } },
          panels.system.systemLoad,
          panels.system.systemContextSwitchesAndInterrupts,
        ]),
      memoryOverview:
        g.panel.row.new('Memory')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels(
          [
            panels.memory.memoryUsageStatPercent { gridPos+: { w: 6, h: 6 } },
            panels.memory.memoryUsageTsBytes { gridPos+: { w: 18, h: 6 } },
          ]
        ),
      memoryVmstat:
        g.panel.row.new('Vmstat')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels(
          [
            panels.memory.memoryPagesInOut,
            panels.memory.memoryPagesSwapInOut,
            panels.memory.memoryPagesFaults,
            panels.memory.memoryOOMkiller,
          ]
        ),
      memoryMemstat:
        g.panel.row.new('Memstat')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels(
          [
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
          ]
        ),
      diskOverview:
        g.panel.row.new('Disk')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels([
          panels.disk.diskIOBytesPerSec { gridPos+: { w: 12, h: 8 } },
          panels.disk.diskUsage { gridPos+: { w: 12, h: 8 } },
        ]),
      disk:
        g.panel.row.new('Disk')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels([
          panels.disk.diskIOBytesPerSec,
          panels.disk.diskIOps,
          panels.disk.diskIOWaitTime,
          panels.disk.diskQueue,
        ]),
      filesystem:
        g.panel.row.new('Filesystem')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels([
          panels.disk.diskFreeTs,
          panels.disk.diskUsage,
          panels.disk.diskInodesFree,
          panels.disk.diskInodesTotal,
          panels.disk.diskErrorsandRO,
          panels.disk.fileDescriptors,
        ]),
      networkOverview:
        g.panel.row.new('Network')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels([
          panels.network.networkUsagePerSec { gridPos+: { w: 12, h: 8 } },
          panels.network.networkErrorsAndDroppedPerSec { gridPos+: { w: 12, h: 8 } },
        ]),
      network:
        g.panel.row.new('Network')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels([
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
        ]),
      networkSockets:
        g.panel.row.new('Network sockets')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels([
          panels.network.networkSockstatAll { gridPos: { w: 24 } },
          panels.network.networkSockstatTCP,
          panels.network.networkSockstatUDP,
          panels.network.networkSockstatMemory,
          panels.network.networkSockstatOther,
        ]),
      networkNetstat:
        g.panel.row.new('Network netstat')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels([
          panels.network.networkNetstatIP { gridPos: { w: 24 } },
          panels.network.networkNetstatTCP,
          panels.network.networkNetstatTCPerrors,
          panels.network.networkNetstatUDP,
          panels.network.networkNetstatUDPerrors,
          panels.network.networkNetstatICMP,
          panels.network.networkNetstatICMPerrors,
        ]),

      hardware:
        g.panel.row.new('Hardware')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels([
          panels.hardware.hardwareTemperature { gridPos+: { w: 12, h: 8 } },
        ]),


      //use
      cpuUseMethod:
        g.panel.row.new('CPU')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels([
          panels.use.cpuUtilization { gridPos+: { w: 12, h: 7 } },
          panels.use.cpuSaturation { gridPos+: { w: 12, h: 7 } },
        ]),
      cpuUseClusterMethod:
        self.cpuUseMethod
        + g.panel.row.withPanels([
          panels.useCluster.cpuUtilization { gridPos+: { w: 12, h: 7 } },
          panels.useCluster.cpuSaturation { gridPos+: { w: 12, h: 7 } },
        ]),
      cpuUseClusterMethodMulti:
        self.cpuUseClusterMethod
        + g.panel.row.withPanels([
          panels.useClusterMulti.cpuUtilization { gridPos+: { w: 12, h: 7 } },
          panels.useClusterMulti.cpuSaturation { gridPos+: { w: 12, h: 7 } },
        ]),


      memoryUseMethod:
        g.panel.row.new('Memory')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels([
          panels.use.memoryUtilization { gridPos+: { w: 12, h: 7 } },
          panels.use.memorySaturation { gridPos+: { w: 12, h: 7 } },
        ]),
      memoryUseClusterMethod:
        self.memoryUseMethod
        + g.panel.row.withPanels([
          panels.useCluster.memoryUtilization { gridPos+: { w: 12, h: 7 } },
          panels.useCluster.memorySaturation { gridPos+: { w: 12, h: 7 } },
        ]),
      memoryUseClusterMethodMulti:
        self.memoryUseMethod
        + g.panel.row.withPanels([
          panels.useClusterMulti.memoryUtilization { gridPos+: { w: 12, h: 7 } },
          panels.useClusterMulti.memorySaturation { gridPos+: { w: 12, h: 7 } },
        ]),


      diskUseMethod:
        g.panel.row.new('Disk')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels([
          panels.use.diskUtilization { gridPos+: { w: 12, h: 7 } },
          panels.use.diskSaturation { gridPos+: { w: 12, h: 7 } },
        ]),
      diskUseClusterMethod:
        self.diskUseMethod
        + g.panel.row.withPanels([
          panels.useCluster.diskUtilization { gridPos+: { w: 12, h: 7 } },
          panels.useCluster.diskSaturation { gridPos+: { w: 12, h: 7 } },
        ]),
      diskUseClusterMethodMulti:
        self.diskUseMethod
        + g.panel.row.withPanels([
          panels.useClusterMulti.diskUtilization { gridPos+: { w: 12, h: 7 } },
          panels.useClusterMulti.diskSaturation { gridPos+: { w: 12, h: 7 } },
        ]),

      filesystemUseMethod:
        g.panel.row.new('Filesystem')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels([
          panels.use.filesystemUtilization { gridPos+: { w: 24, h: 7 } },
        ]),
      filesystemUseClusterMethod:
        self.filesystemUseMethod
        + g.panel.row.withPanels([
          panels.useCluster.filesystemUtilization { gridPos+: { w: 24, h: 7 } },
        ]),
      filesystemUseClusterMethodMulti:
        self.filesystemUseMethod
        + g.panel.row.withPanels([
          panels.useClusterMulti.filesystemUtilization { gridPos+: { w: 24, h: 7 } },
        ]),

      networkUseMethod:
        g.panel.row.new('Network')
        + g.panel.row.withCollapsed(false)
        + g.panel.row.withPanels([
          panels.use.networkUtilization { gridPos+: { w: 12, h: 7 } },
          panels.use.networkSaturation { gridPos+: { w: 12, h: 7 } },
        ]),
      networkUseClusterMethod:
        self.networkUseMethod
        + g.panel.row.withPanels([
          panels.useCluster.networkUtilization { gridPos+: { w: 12, h: 7 } },
          panels.useCluster.networkSaturation { gridPos+: { w: 12, h: 7 } },
        ]),
      networkUseClusterMethodMulti:
        self.networkUseMethod
        + g.panel.row.withPanels([
          panels.useClusterMulti.networkUtilization { gridPos+: { w: 12, h: 7 } },
          panels.useClusterMulti.networkSaturation { gridPos+: { w: 12, h: 7 } },
        ]),
    },
}
