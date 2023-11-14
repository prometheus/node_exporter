local g = import '../../../g.libsonnet';
local generic = import '../../generic/stat/main.libsonnet';
local base = import './base.libsonnet';
local stat = g.panel.stat;

base {
  new(
    title='CPU count',
    targets,
    description=|||
      CPU count is the number of processor cores or central processing units (CPUs) in a computer,
      determining its processing capability and ability to handle tasks concurrently.
    |||
  ):
    super.new(title, targets, description),

  stylize(allLayers=true):
    (if allLayers then super.stylize() else {})
    + generic.info.stylize(allLayers=false),

}
