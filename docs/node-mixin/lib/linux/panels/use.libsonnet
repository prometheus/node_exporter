local g = import '../../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
local utils = commonlib.utils;
{
  new(this):
    {
      local t = this.grafana.targets,
      local table = g.panel.table,
      local fieldOverride = g.panel.table.fieldOverride,


      //for USE
      cpuUtilization:
        commonlib.panels.cpu.timeSeries.utilization.new(targets=[t.use.cpuUtilization])
        + g.panel.timeSeries.panelOptions.withTitle('CPU utilization')
        + g.panel.timeSeries.options.legend.withShowLegend(false)
        + g.panel.timeSeries.fieldConfig.defaults.custom.stacking.withMode('standard'),
      cpuSaturation:
        commonlib.panels.cpu.timeSeries.utilization.new(targets=[t.use.cpuSaturation])
        + g.panel.timeSeries.panelOptions.withTitle('CPU saturation (Load1 per CPU)')
        + g.panel.timeSeries.options.legend.withShowLegend(false)
        + g.panel.timeSeries.panelOptions.withDescription(
          |||
            System load average over the last minute. A measurement of how many processes are waiting for CPU cycles. The value is as a percent compared to the number of CPU cores for the node.
          |||
        )
        + { title: 'CPU saturation (Load 1 per CPU)' }
        + g.panel.timeSeries.fieldConfig.defaults.custom.stacking.withMode('standard'),

      memoryUtilization:
        commonlib.panels.memory.timeSeries.usagePercent.new(targets=[t.use.memoryUtilization])
        + g.panel.timeSeries.panelOptions.withTitle('Memory utilization')
        + g.panel.timeSeries.options.legend.withShowLegend(false)
        + g.panel.timeSeries.fieldConfig.defaults.custom.stacking.withMode('standard'),
      memorySaturation:
        commonlib.panels.memory.timeSeries.base.new(
          'Memory saturation (Major page faults)',
          targets=[t.use.memorySaturation],
        )
        + g.panel.timeSeries.panelOptions.withDescription(this.grafana.panels.memory.memoryPagesFaults.description)
        + g.panel.timeSeries.panelOptions.withTitle('Memory saturation (Major page faults)')
        + g.panel.timeSeries.options.legend.withShowLegend(false)
        + g.panel.timeSeries.fieldConfig.defaults.custom.stacking.withMode('standard'),

      networkUtilization:
        this.grafana.panels.network.networkUsagePerSec
        + g.panel.timeSeries.panelOptions.withTitle('Network utilization (Bytes receive/transmit)')
        + g.panel.timeSeries.queryOptions.withTargets([t.use.networkUtilizationReceive, t.use.networkUtilizationTransmit])
        + commonlib.panels.network.timeSeries.base.withNegateOutPackets('/Transmit/')
        + g.panel.timeSeries.options.legend.withShowLegend(false)
        + g.panel.timeSeries.fieldConfig.defaults.custom.stacking.withMode('standard'),
      networkSaturation:
        this.grafana.panels.network.networkDroppedPerSec
        + g.panel.timeSeries.panelOptions.withTitle('Network saturation (Drops receive/transmit)')
        + g.panel.timeSeries.queryOptions.withTargets([t.use.networkSaturationReceive, t.use.networkSaturationTransmit])
        + g.panel.timeSeries.fieldConfig.defaults.custom.stacking.withMode('standard')
        + g.panel.timeSeries.options.legend.withShowLegend(false),


      diskUtilization:
        commonlib.panels.generic.timeSeries.base.new(
          'Disk IO utilization', targets=[t.use.diskUtilization], description='Disk total IO seconds'
        )
        + g.panel.timeSeries.options.legend.withShowLegend(false)
        + g.panel.timeSeries.standardOptions.withUnit('percent')
        + g.panel.timeSeries.fieldConfig.defaults.custom.stacking.withMode('standard'),

      diskSaturation:
        commonlib.panels.generic.timeSeries.base.new(
          'Disk IO saturation', targets=[t.use.diskSaturation], description='Disk saturation (weighted seconds spent, 1 second rate)'
        )
        + g.panel.timeSeries.options.legend.withShowLegend(false)
        + g.panel.timeSeries.standardOptions.withUnit('percent')
        + g.panel.timeSeries.fieldConfig.defaults.custom.stacking.withMode('standard'),

      filesystemUtilization:
        this.grafana.panels.disk.diskFreeTs
        + g.panel.timeSeries.panelOptions.withTitle('Filesytem utilization')
        + g.panel.timeSeries.options.legend.withShowLegend(false)
        + g.panel.timeSeries.queryOptions.withTargets([t.use.filesystemUtilization])
        + g.panel.timeSeries.standardOptions.withUnit('percent')
        + g.panel.timeSeries.standardOptions.withMax(100)
        + g.panel.timeSeries.standardOptions.withMin(0)
        + g.panel.timeSeries.fieldConfig.defaults.custom.stacking.withMode('standard')
        + g.panel.timeSeries.panelOptions.withDescription('Total disk utilization percent.'),

    },
}
