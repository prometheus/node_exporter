local g = import '../../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
{
  new(this):
    {
      local panels = this.grafana.panels,
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
