local g = import '../../../g.libsonnet';
local generic = import '../../generic/stat/main.libsonnet';
local base = import './base.libsonnet';
local stat = g.panel.stat;

base {
  new(
    title='CPU usage',
    targets,
    description=|||
      Total CPU utilization percent is a metric that indicates the overall level of central processing unit (CPU) usage in a computer system.
      It represents the combined load placed on all CPU cores or processors.

      For instance, if the total CPU utilization percent is 50%, it means that,
      on average, half of the CPU's processing capacity is being used to execute tasks. A higher percentage indicates that the CPU is working more intensively, potentially leading to system slowdowns if it remains consistently high.
    |||
  ):
    super.new(title, targets, description),
  stylize(allLayers=true):
    (if allLayers then super.stylize() else {})
    + generic.percentage.stylize(allLayers=false),
}
