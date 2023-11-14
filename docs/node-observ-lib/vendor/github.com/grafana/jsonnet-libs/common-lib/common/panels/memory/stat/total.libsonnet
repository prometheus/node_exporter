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
    title='Memory total',
    targets,
    description=|||
      Amount of random-access memory (RAM) installed.
      It represents the system's available working memory that applications and the operating system use to perform tasks.
      A higher memory total generally leads to better system performance and the ability to run more demanding applications and processes simultaneously.
    |||
  ):
    super.new(title=title, targets=targets, description=description),

  stylize(allLayers=true):

    (if allLayers then super.stylize() else {})

    + generic.info.stylize(allLayers=false)
    + stat.standardOptions.withUnit('bytes'),
}
