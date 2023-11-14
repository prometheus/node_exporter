local g = import '../../../g.libsonnet';
local base = import './base.libsonnet';
local timeSeries = g.panel.timeSeries;
local fieldOverride = g.panel.timeSeries.fieldOverride;
local custom = timeSeries.fieldConfig.defaults.custom;
local defaults = timeSeries.fieldConfig.defaults;
local options = timeSeries.options;
base {
  new(
    title='Disk average wait time',
    targets,
    description=|||
      The average time for requests issued to the device to be served.
      This includes the time spent by the requests in queue and the time spent servicing them.'
    |||
  ):
    super.new(title, targets, description)
    + self.stylize(),

  stylize(allLayers=true):

    (if allLayers == true then super.stylize() else {})
    + timeSeries.standardOptions.withUnit('s')
    + self.withNegateOutPackets(),
}
