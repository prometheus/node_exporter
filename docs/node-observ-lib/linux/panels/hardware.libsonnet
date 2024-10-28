local g = import '../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
local utils = commonlib.utils;
{
  new(this):
    {
      local t = this.grafana.targets,
      local table = g.panel.table,
      local fieldOverride = g.panel.table.fieldOverride,
      local instanceLabel = this.config.instanceLabels[0],
      hardwareTemperature:
        commonlib.panels.hardware.timeSeries.temperature.new(
          'Temperature',
          targets=[t.hardwareTemperature]
        ),
    },
}
