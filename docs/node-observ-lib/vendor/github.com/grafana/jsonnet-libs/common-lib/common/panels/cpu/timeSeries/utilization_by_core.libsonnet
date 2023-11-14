local g = import '../../../g.libsonnet';
local generic = import '../../generic/timeSeries/main.libsonnet';
local base = import './base.libsonnet';
base {
  new(
    title='CPU usage',
    targets,
    description=|||
      CPU utilization percent by core is a metric that indicates level of central processing unit (CPU) usage in a computer system.
      It represents the load placed on each CPU core or processors.
    |||
  ):
    super.new(title, targets, description)
    + self.stylize(),

  stylize(allLayers=true):
    (if allLayers then super.stylize() else {})
    + generic.percentage.stylize(allLayers=false)
    + g.panel.timeSeries.fieldConfig.defaults.custom.withStacking({ mode: 'normal' }),
}
