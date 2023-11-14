local g = import '../../../g.libsonnet';
local timeSeries = g.panel.timeSeries;
local fieldOverride = g.panel.timeSeries.fieldOverride;
local fieldConfig = g.panel.timeSeries.fieldConfig;

// Turns any series to threshold line: dashed red line without gradient fill
{
  local this = self,

  stylize():
    fieldConfig.defaults.custom.withLineStyleMixin(
      {
        fill: 'dash',
        dash: [10, 10],
      }
    )
    + fieldConfig.defaults.custom.withFillOpacity(0)
    + timeSeries.standardOptions.color.withMode('fixed')
    + timeSeries.standardOptions.color.withFixedColor('light-orange'),
  stylizeByRegexp(regexp):
    timeSeries.standardOptions.withOverrides(
      fieldOverride.byRegexp.new(regexp)
      + fieldOverride.byRegexp.withPropertiesFromOptions(this.stylize())
    ),
}
