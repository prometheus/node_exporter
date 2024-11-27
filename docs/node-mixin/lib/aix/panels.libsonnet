local g = import '../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';

{
  new(this):
    {
      local t = this.grafana.targets,
      local table = g.panel.table,
      local fieldOverride = g.panel.table.fieldOverride,
      local instanceLabel = this.config.instanceLabels[0],

      // override description and targets
      memory+: {
        memoryUsageTsBytes+:
          g.panel.timeSeries.queryOptions.withTargets([
            t.memory.memoryUsedBytes,
            t.memory.memoryTotalBytes,
            t.memory.memorySwapUsedBytes,
          ])
          + commonlib.panels.generic.timeSeries.threshold.stylizeByRegexp('Physical memory'),
      },

    },
}
