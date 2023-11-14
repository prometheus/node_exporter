local g = import '../../../g.libsonnet';
local base = import './base.libsonnet';
local timeSeries = g.panel.timeSeries;
local fieldOverride = g.panel.timeSeries.fieldOverride;
local custom = timeSeries.fieldConfig.defaults.custom;
local defaults = timeSeries.fieldConfig.defaults;
local options = timeSeries.options;
base {
  new(
    title='Network traffic',
    targets,
    description='Network traffic (bits per sec) measures data transmitted and received.',
  ):
    super.new(title, targets, description)
    + self.stylize(),

  stylize(allLayers=true):

    (if allLayers == true then super.stylize() else {})
    + timeSeries.standardOptions.withUnit('bps'),
}
