local g = import '../../../g.libsonnet';
local uptime = import '../stat/uptime.libsonnet';
local base = import './base.libsonnet';
local table = g.panel.table;
local fieldOverride = table.fieldOverride;

base {
  local this = self,

  new(): error 'not supported',
  stylize(): error 'not supported',

  // when attached to table, this function applies style to row named 'name="Uptime"'
  stylizeByName(name='Uptime'):
    table.standardOptions.withOverrides(
      fieldOverride.byName.new(name)
      + fieldOverride.byName.withProperty('custom.cellOptions', { type: 'color-text' })
      + fieldOverride.byName.withPropertiesFromOptions(uptime.stylize(allLayers=false),)
    ),
}
