local g = import '../../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
local utils = commonlib.utils;
{
  new(this):
    {
      local t = this.grafana.targets,
      local table = g.panel.table,
      local fieldOverride = g.panel.table.fieldOverride,
      local instancePanels = this.grafana.panels.use,
      //for USE
      cpuUtilization: instancePanels.cpuUtilization
                      + g.panel.timeSeries.queryOptions.withTargets([t.useClusterMulti.cpuUtilization]),
      cpuSaturation: instancePanels.cpuSaturation
                     + g.panel.timeSeries.queryOptions.withTargets([t.useClusterMulti.cpuSaturation]),

      memoryUtilization:
        instancePanels.memoryUtilization
        + g.panel.timeSeries.queryOptions.withTargets([t.useClusterMulti.memoryUtilization]),

      memorySaturation:
        instancePanels.memorySaturation
        + g.panel.timeSeries.queryOptions.withTargets([t.useClusterMulti.memorySaturation]),

      networkUtilization:
        instancePanels.networkUtilization
        + g.panel.timeSeries.queryOptions.withTargets([t.useClusterMulti.networkUtilizationReceive, t.useClusterMulti.networkUtilizationTransmit]),
      networkSaturation:
        instancePanels.networkSaturation
        + g.panel.timeSeries.queryOptions.withTargets([t.useClusterMulti.networkSaturationReceive, t.useClusterMulti.networkSaturationTransmit]),


      diskUtilization:
        instancePanels.diskUtilization
        + g.panel.timeSeries.queryOptions.withTargets([t.useClusterMulti.diskUtilization]),

      diskSaturation:
        instancePanels.diskSaturation
        + g.panel.timeSeries.queryOptions.withTargets([t.useClusterMulti.diskSaturation]),

      filesystemUtilization:
        instancePanels.filesystemUtilization
        + g.panel.timeSeries.queryOptions.withTargets([t.useClusterMulti.filesystemUtilization]),

    },
}
