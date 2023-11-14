local g = import '../../../g.libsonnet';
local base = import './base.libsonnet';
local timeSeries = g.panel.timeSeries;
local fieldOverride = g.panel.timeSeries.fieldOverride;
local custom = timeSeries.fieldConfig.defaults.custom;
local defaults = timeSeries.fieldConfig.defaults;
local options = timeSeries.options;
base {
  new(
    title='Disk reads/writes',
    targets,
    description='Disk read/writes in bytes per second.',
  ):
    super.new(title, targets, description)
    + self.stylize(),

  stylize(allLayers=true):

    (if allLayers == true then super.stylize() else {})

    + timeSeries.standardOptions.withUnit('Bps')
    // move 'IO busy time' to second axis if found
    + timeSeries.standardOptions.withOverrides(
      fieldOverride.byRegexp.new('/time|used|busy|util/')
      + fieldOverride.byRegexp.withPropertiesFromOptions(
        timeSeries.standardOptions.withUnit('percent')
        + timeSeries.fieldConfig.defaults.custom.withDrawStyle('points')
        + timeSeries.fieldConfig.defaults.custom.withAxisSoftMax(100)
      )
    ),
}
