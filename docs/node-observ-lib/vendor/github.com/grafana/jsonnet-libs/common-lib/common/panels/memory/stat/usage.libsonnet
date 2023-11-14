local g = import '../../../g.libsonnet';
local generic = import '../../generic/stat/main.libsonnet';
local base = import './base.libsonnet';
local stat = g.panel.stat;
local fieldOverride = g.panel.stat.fieldOverride;
local custom = stat.fieldConfig.defaults.custom;
local defaults = stat.fieldConfig.defaults;
local options = stat.options;
base {
  new(
    title='Memory usage',
    targets,
    description='RAM (random-access memory) currently in use by the operating system and running applications, in percent.'
  ):
    super.new(title=title, targets=targets, description=description),

  stylize(allLayers=true):

    (if allLayers then super.stylize() else {})

    + generic.percentage.stylize(allLayers=false)
    + stat.standardOptions.withUnit('percent'),

}
