local g = import '../../../g.libsonnet';
local base = import './base.libsonnet';
local timeSeries = g.panel.timeSeries;
local fieldOverride = g.panel.timeSeries.fieldOverride;
local custom = timeSeries.fieldConfig.defaults.custom;
local defaults = timeSeries.fieldConfig.defaults;
local options = timeSeries.options;
base {
  new(
    title='Network packets',
    targets,
    description='Network packet count tracks the number of data packets transmitted and received over a network connection, providing insight into network activity and performance.',
  ):
    super.new(title, targets, description),
}
