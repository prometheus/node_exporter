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

      cpuCount: commonlib.panels.cpu.stat.count.new(targets=[t.cpuCount]),
      cpuUsageTsPerCore: commonlib.panels.cpu.timeSeries.utilization.new(targets=[t.cpuUsagePerCore])
                         + g.panel.timeSeries.fieldConfig.defaults.custom.withStacking({ mode: 'normal' }),

      cpuUsageTopk: commonlib.panels.generic.timeSeries.topkPercentage.new(
        title='CPU usage',
        target=t.cpuUsage,
        topk=25,
        instanceLabels=this.config.instanceLabels,
        drillDownDashboardUid=this.grafana.dashboards['overview.json'].uid,
      ),
      cpuUsageStat: commonlib.panels.cpu.stat.usage.new(targets=[t.cpuUsage]),
      cpuUsageByMode: commonlib.panels.cpu.timeSeries.utilizationByMode.new(
        targets=[t.cpuUsageByMode],
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
