local g = import '../../../g.libsonnet';
local base = import '../../generic/timeSeries/base.libsonnet';
local timeSeries = g.panel.timeSeries;
local fieldOverride = g.panel.timeSeries.fieldOverride;
local custom = timeSeries.fieldConfig.defaults.custom;
local defaults = timeSeries.fieldConfig.defaults;
local options = timeSeries.options;
base {

  stylize(allLayers=true):

    (if allLayers == true then super.stylize() else {})
    // Decrease opacity (would look better with too many timeseries)
    + defaults.custom.withFillOpacity(1),


  withNegateOutPackets(regexp='/write|written/'):
    defaults.custom.withAxisLabel('write(-) | read(+)')
    + defaults.custom.withAxisCenteredZero(true)
    + timeSeries.standardOptions.withOverrides(
      fieldOverride.byRegexp.new(regexp)
      + fieldOverride.byRegexp.withPropertiesFromOptions(
        defaults.custom.withTransform('negative-Y')
      )
    ),

}
