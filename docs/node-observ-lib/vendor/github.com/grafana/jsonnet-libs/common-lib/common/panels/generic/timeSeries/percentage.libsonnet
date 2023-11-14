local g = import '../../../g.libsonnet';
local timeSeries = g.panel.timeSeries;
local base = import './base.libsonnet';
// This panel can be used to display gauge metrics with possible values range 0-100%.
// Examples: cpu utilization, memory utilization etc.
base {
  stylize(allLayers=true):
    (if allLayers then super.stylize() else {})
    + timeSeries.standardOptions.withDecimals(1)
    + timeSeries.standardOptions.withUnit('percent')
    // Change color from blue(cold) to red(hot)
    + timeSeries.standardOptions.color.withMode('continuous-BlYlRd')
    + timeSeries.fieldConfig.defaults.custom.withGradientMode('scheme')
    + timeSeries.standardOptions.withMax(100)
    + timeSeries.standardOptions.withMin(0),
}
