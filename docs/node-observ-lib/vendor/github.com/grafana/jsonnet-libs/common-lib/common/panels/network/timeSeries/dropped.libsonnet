local g = import '../../../g.libsonnet';
local base = import './base.libsonnet';
local timeSeries = g.panel.timeSeries;
local fieldOverride = g.panel.timeSeries.fieldOverride;
local custom = timeSeries.fieldConfig.defaults.custom;
local defaults = timeSeries.fieldConfig.defaults;
local options = timeSeries.options;
base {
  new(
    title='Dropped packets',
    targets,
    description=|||
      Dropped packets occur when data packets traveling through a network are intentionally discarded or lost due to congestion, resource limitations, or network configuration issues. 

      Common causes include network congestion, buffer overflows, QoS settings, and network errors, as corrupted or incomplete packets may be discarded by receiving devices.

      Dropped packets can impact network performance and lead to issues such as degraded voice or video quality in real-time applications.
    |||,
  ):
    super.new(title, targets, description)
    + self.stylize(),

  stylize(allLayers=true):

    (if allLayers == true then super.stylize() else {})

    + timeSeries.standardOptions.withNoValue('No dropped packets'),
}
