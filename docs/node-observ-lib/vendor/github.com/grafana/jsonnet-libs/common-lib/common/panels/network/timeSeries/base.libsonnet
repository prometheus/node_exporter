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

    + timeSeries.standardOptions.withDecimals(1)
    + timeSeries.standardOptions.withUnit('pps'),

  withNegateOutPackets(regexp='/transmit|tx|out/'):
    defaults.custom.withAxisLabel('out(-) | in(+)')
    + defaults.custom.withAxisCenteredZero(false)
    + timeSeries.standardOptions.withOverrides(
      fieldOverride.byRegexp.new(regexp)
      + fieldOverride.byRegexp.withPropertiesFromOptions(
        defaults.custom.withTransform('negative-Y')
      )
    ),
}
