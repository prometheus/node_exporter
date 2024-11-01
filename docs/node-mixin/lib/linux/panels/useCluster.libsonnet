local g = import '../../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
local utils = commonlib.utils;
{
  new(this):
    {
      local t = this.grafana.targets,
      local table = g.panel.table,
      local fieldOverride = g.panel.table.fieldOverride,
      local instanceLabel = this.config.instanceLabels[0],
      local instancePanels = this.grafana.panels.use,
      //for USE
      cpuUtilization:
        instancePanels.cpuUtilization
        + g.panel.timeSeries.queryOptions.withTargets([t.useCluster.cpuUtilization]),
      cpuSaturation:
        instancePanels.cpuSaturation
        + g.panel.timeSeries.queryOptions.withTargets([t.useCluster.cpuSaturation]),

      memoryUtilization:
        instancePanels.memoryUtilization
        + g.panel.timeSeries.queryOptions.withTargets([t.useCluster.memoryUtilization]),

      memorySaturation:
        instancePanels.memorySaturation
        + g.panel.timeSeries.queryOptions.withTargets([t.useCluster.memorySaturation]),

      networkUtilization:
        instancePanels.networkUtilization
        + g.panel.timeSeries.queryOptions.withTargets([t.useCluster.networkUtilizationReceive, t.useCluster.networkUtilizationTransmit]),
      networkSaturation:
        instancePanels.networkSaturation
        + g.panel.timeSeries.queryOptions.withTargets([t.useCluster.networkSaturationReceive, t.useCluster.networkSaturationTransmit]),


      diskUtilization:
        instancePanels.diskUtilization
        + g.panel.timeSeries.queryOptions.withTargets([t.useCluster.diskUtilization]),

      diskSaturation:
        instancePanels.diskSaturation
        + g.panel.timeSeries.queryOptions.withTargets([t.useCluster.diskSaturation]),

      filesystemUtilization:
        instancePanels.filesystemUtilization
        + g.panel.timeSeries.queryOptions.withTargets([t.useCluster.filesystemUtilization]),

    },
}
