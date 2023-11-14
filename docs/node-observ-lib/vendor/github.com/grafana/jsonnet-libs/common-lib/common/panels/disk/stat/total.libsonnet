local g = import '../../../g.libsonnet';
local generic = import '../../generic/stat/main.libsonnet';
local base = import './base.libsonnet';
local stat = g.panel.stat;

base {
  new(
    title,
    targets,
    description=''
  ):
    super.new(title=title, targets=targets, description=description),

  stylize(allLayers=true):

    (if allLayers then super.stylize() else {})

    + generic.info.stylize(allLayers=false)
    + stat.standardOptions.withUnit('bytes'),
}
