local g = import '../../../g.libsonnet';
local base = import '../base.libsonnet';
local statusHistory = g.panel.statusHistory;
local fieldOverride = g.panel.statusHistory.fieldOverride;
local custom = statusHistory.fieldConfig.defaults.custom;
local defaults = statusHistory.fieldConfig.defaults;
local options = statusHistory.options;
base {

  new(title, targets, description=''):
    statusHistory.new(title)
    + super.new(targets, description)
    // Minimize number of points to avoid 'Too many data points' error on large time intervals
    + statusHistory.queryOptions.withMaxDataPoints(50),

  stylize(allLayers=true):
    (if allLayers then super.stylize() else {}),
}
