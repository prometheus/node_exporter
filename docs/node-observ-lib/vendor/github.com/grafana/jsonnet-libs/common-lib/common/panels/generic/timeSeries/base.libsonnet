local g = import '../../../g.libsonnet';

local timeSeries = g.panel.timeSeries;
local fieldOverride = g.panel.timeSeries.fieldOverride;
local custom = timeSeries.fieldConfig.defaults.custom;
local defaults = timeSeries.fieldConfig.defaults;
local options = timeSeries.options;
local base = import '../base.libsonnet';
base {
  new(title, targets, description=''):
    timeSeries.new(title)
    + super.new(targets, description),

  stylize(allLayers=true):
    (if allLayers then super.stylize() else {})
    // Style choice: Make lines more thick
    + custom.withLineWidth(2)
    // Style choice: Opacity level
    + custom.withFillOpacity(30)
    // Style choice: Don't show points on lines
    + custom.withShowPoints('never')
    // Style choice: Opacity gradient
    + custom.withGradientMode('opacity')
    // Style choice: Smoother lines
    + custom.withLineInterpolation('smooth')
    // Style choice: Show all values in tooltip, sorted
    + options.tooltip.withMode('multi')
    + options.tooltip.withSort('desc')
    // Style choice: Use simple legend without any values (cleaner look)
    + options.legend.withDisplayMode('list')
    + options.legend.withCalcs([]),

}
