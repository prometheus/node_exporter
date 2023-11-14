local g = import '../../g.libsonnet';

local timeSeries = g.panel.timeSeries;
local fieldOverride = g.panel.timeSeries.fieldOverride;
local custom = timeSeries.fieldConfig.defaults.custom;
local defaults = timeSeries.fieldConfig.defaults;
local options = timeSeries.options;

// This is base of ALL panels in the common lib
{
  new(targets, description=''):
    // hidden field to hold styles modifiers

    timeSeries.queryOptions.withTargets(targets)
    + timeSeries.panelOptions.withDescription(description)
    // set first target's datasource
    // to panel's datasource if only single type of
    // datasoures are used accross all targets:
    + (if std.length(std.set(targets, function(t) t.datasource.type)) == 1 then
         timeSeries.queryOptions.withDatasource(
           targets[0].datasource.type, targets[0].datasource.uid
         ) else {})
    + self.stylize(),

  stylize(): {},

}
