local g = import '../../../g.libsonnet';

local table = g.panel.table;
local fieldOverride = g.panel.table.fieldOverride;
local custom = table.fieldConfig.defaults.custom;
local defaults = table.fieldConfig.defaults;
local options = table.options;
local base = import '../base.libsonnet';

base {
  new(title, targets, description=''):
    table.new(title)
    + super.new(targets, description),

  stylize(allLayers=true):
    (if allLayers then super.stylize() else {}),

  transformations+: {
    sortBy(field, desc=false):
      {
        id: 'sortBy',
        options: {
          fields: {},
          sort: [
            {
              field: field,
              desc: desc,
            },
          ],
        },
      },
  },
}
