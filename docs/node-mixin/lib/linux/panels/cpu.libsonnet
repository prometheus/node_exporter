local g = import '../../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
local utils = commonlib.utils;
local xtd = import 'github.com/jsonnet-libs/xtd/main.libsonnet';
{
  new(this):
    {
      local t = this.grafana.targets,
      local table = g.panel.table,
      local fieldOverride = g.panel.table.fieldOverride,
      local instanceLabel = xtd.array.slice(this.config.instanceLabels, -1)[0],

      cpuCount: commonlib.panels.cpu.stat.count.new(targets=[t.cpu.cpuCount]),
      cpuUsageTsPerCore: commonlib.panels.cpu.timeSeries.utilization.new(targets=[t.cpu.cpuUsagePerCore])
                         + g.panel.timeSeries.fieldConfig.defaults.custom.withStacking({ mode: 'normal' }),

      cpuUsageTopk: commonlib.panels.generic.timeSeries.topkPercentage.new(
        title='CPU usage',
        target=t.cpu.cpuUsage,
        topk=25,
        instanceLabels=this.config.instanceLabels,
        drillDownDashboardUid=this.grafana.dashboards['nodes.json'].uid,
      ),
      cpuUsageStat: commonlib.panels.cpu.stat.usage.new(targets=[t.cpu.cpuUsage]),
      cpuUsageByMode: commonlib.panels.cpu.timeSeries.utilizationByMode.new(
        targets=[t.cpu.cpuUsageByMode],
        description=|||
          - System: Processes executing in kernel mode.
          - User: Normal processes executing in user mode.
          - Nice: Niced processes executing in user mode.
          - Idle: Waiting for something to happen.
          - Iowait: Waiting for I/O to complete.
          - Irq: Servicing interrupts.
          - Softirq: Servicing softirqs.
          - Steal: Time spent in other operating systems when running in a virtualized environment.
        |||
      ),
    },
}
