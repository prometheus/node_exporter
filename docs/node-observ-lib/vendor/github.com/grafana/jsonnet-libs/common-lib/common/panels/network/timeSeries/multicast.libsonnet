local g = import '../../../g.libsonnet';
local base = import './base.libsonnet';
local timeSeries = g.panel.timeSeries;
local fieldOverride = g.panel.timeSeries.fieldOverride;
local custom = timeSeries.fieldConfig.defaults.custom;
local defaults = timeSeries.fieldConfig.defaults;
local options = timeSeries.options;
base {
  new(
    title='Multicast packets',
    targets,
    description='Packets sent from one source to multiple recipients simultaneously, allowing efficient one-to-many communication in a network.',
  ):
    super.new(title, targets, description),
}
