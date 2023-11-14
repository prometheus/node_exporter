local g = import '../../../g.libsonnet';
local base = import './base.libsonnet';
local timeSeries = g.panel.timeSeries;
local fieldOverride = g.panel.timeSeries.fieldOverride;
local custom = timeSeries.fieldConfig.defaults.custom;
local defaults = timeSeries.fieldConfig.defaults;
local options = timeSeries.options;
base {
  new(
    title='Network errors',
    targets,
    description=|||
      Network errors refer to issues that occur during the transmission of data across a network. 

      These errors can result from various factors, including physical issues, jitter, collisions, noise and interference.

      Monitoring network errors is essential for diagnosing and resolving issues, as they can indicate problems with network hardware or environmental factors affecting network quality.
    |||,
  ):
    super.new(title, targets, description)
    + self.stylize(),
  stylize(allLayers=true):

    (if allLayers == true then super.stylize() else {})
    + timeSeries.standardOptions.withNoValue('No errors'),
}
