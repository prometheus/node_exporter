local g = import '../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
local utils = commonlib.utils;
{
  new(this):
    {
      local t = this.grafana.targets,
      local table = g.panel.table,
      local fieldOverride = g.panel.table.fieldOverride,
      local instanceLabel = this.config.instanceLabels[0],

      // override description and targets
      memoryUsageTsBytes+:
        g.panel.timeSeries.panelOptions.withDescription(
          |||
            - Physical memory: Total amount of memory installed in this computer;
            - App memory: Physical memory allocated by apps and system processes;
            - Wired memory: Physical memory, containing data that cannot be compressed or swapped to disk;
            - Compressed memory: Physical memory used to store a compressed version of data that has not been used recently;
            - Swap used: Amount of compressed data temporarily moved to disk to make room in memory for more recently used data.
          |||
        )
        + g.panel.timeSeries.queryOptions.withTargets([
          t.memoryUsedBytes,
          t.memoryTotalBytes,
          t.memoryAppBytes,
          t.memoryWiredBytes,
          t.memoryCompressedBytes,
          t.memorySwapUsedBytes,
        ])
        + commonlib.panels.generic.timeSeries.threshold.stylizeByRegexp('Physical memory'),

      //override reduceOption field to version
      osInfo+:
        g.panel.timeSeries.panelOptions.withTitle('OS version')
        + { options+: { reduceOptions: { fields: '/^version$/' } } },
    },
}
