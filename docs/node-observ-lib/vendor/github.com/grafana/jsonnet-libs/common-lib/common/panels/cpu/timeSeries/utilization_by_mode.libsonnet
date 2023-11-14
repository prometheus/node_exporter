local g = import '../../../g.libsonnet';
local base = import './base.libsonnet';
base {
  new(
    title='CPU usage by modes',
    targets,
    description='CPU usage by different modes.'
  ):
    super.new(title, targets, description)
    + self.stylize(),

  stylize(allLayers=true):
    local timeSeries = g.panel.timeSeries;
    local fieldOverride = g.panel.timeSeries.fieldOverride;

    (if allLayers then super.stylize() else {})

    + timeSeries.standardOptions.withUnit('percent')
    + timeSeries.fieldConfig.defaults.custom.withFillOpacity(80)
    + timeSeries.fieldConfig.defaults.custom.withStacking({ mode: 'normal' })
    + timeSeries.standardOptions.withOverrides(
      [
        fieldOverride.byName.new('idle')
        + fieldOverride.byName.withPropertiesFromOptions(
          timeSeries.standardOptions.color.withMode('fixed')
          + timeSeries.standardOptions.color.withFixedColor('light-blue'),
        ),
        fieldOverride.byName.new('interrupt')
        + fieldOverride.byName.withPropertiesFromOptions(
          timeSeries.standardOptions.color.withMode('fixed')
          + timeSeries.standardOptions.color.withFixedColor('light-purple'),
        ),
        fieldOverride.byName.new('user')
        + fieldOverride.byName.withPropertiesFromOptions(
          timeSeries.standardOptions.color.withMode('fixed')
          + timeSeries.standardOptions.color.withFixedColor('light-orange'),
        ),
        fieldOverride.byRegexp.new('system|privileged')
        + fieldOverride.byRegexp.withPropertiesFromOptions(
          timeSeries.standardOptions.color.withMode('fixed')
          + timeSeries.standardOptions.color.withFixedColor('light-red'),
        ),
      ]
    ),
}
