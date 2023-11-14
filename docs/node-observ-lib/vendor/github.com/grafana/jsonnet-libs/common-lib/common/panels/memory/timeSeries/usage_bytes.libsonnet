local g = import '../../../g.libsonnet';
local generic = import '../../generic/timeSeries/main.libsonnet';
local base = import './base.libsonnet';
local timeSeries = g.panel.timeSeries;
local fieldOverride = g.panel.timeSeries.fieldOverride;
local custom = timeSeries.fieldConfig.defaults.custom;
local defaults = timeSeries.fieldConfig.defaults;
local options = timeSeries.options;
base {
  totalRegexp:: '.*(T|t)otal.*',
  new(
    title='Memory usage',
    targets,
    description=|||
      RAM (random-access memory) currently in use by the operating system and running applications, in bytes.
    |||,
    totalRegexp=self.totalRegexp,
  ):
    super.new(title=title, targets=targets, description=description),

  stylize(allLayers=true, totalRegexp=self.totalRegexp):
    (if allLayers then super.stylize() else {})
    + timeSeries.standardOptions.withUnit('bytes')
    + timeSeries.standardOptions.withMin(0)
    + generic.threshold.stylizeByRegexp(totalRegexp),
}
