local g = import '../../../g.libsonnet';
local base = import './base.libsonnet';
local table = g.panel.table;
local fieldOverride = g.panel.table.fieldOverride;
local custom = table.fieldConfig.defaults.custom;
local defaults = table.fieldConfig.defaults;
local options = table.options;

base {
  new(
    title='Disk space usage',
    totalTarget,
    usageTarget=null,
    freeTarget=null,
    groupLabel,
    description=|||
      This table provides information about total disk space, used space, available space, and usage percentages for each mounted file system on the system.
    |||,
  ):
    // validate inputs
    std.prune(
      {
        checks: [
          if (usageTarget == null && freeTarget == null) then error 'Must provide at leason one of "usageTarget" or "freeTraget"',
          if !(std.objectHas(totalTarget, 'format') && std.assertEqual(totalTarget.format, 'table')) then error 'totalTarget format must be "table"',
          if !(std.objectHas(totalTarget, 'instant') && std.assertEqual(totalTarget.instant, true)) then error 'totalTarget must have type "instant"',
          if usageTarget != null && !(std.objectHas(usageTarget, 'format') && std.assertEqual(usageTarget.format, 'table')) then error 'usageTarget format must be "table"',
          if usageTarget != null && !(std.objectHas(usageTarget, 'instant') && std.assertEqual(usageTarget.instant, true)) then error 'usageTarget must have  type "instant"',
          if freeTarget != null && !(std.objectHas(freeTarget, 'format') && std.assertEqual(freeTarget.format, 'table')) then error 'freeTarget format must be "table"',
          if freeTarget != null && !(std.objectHas(freeTarget, 'instant') && std.assertEqual(freeTarget.instant, true)) then error 'freeTarget must have  type "instant"',
        ],
      }
    ) +
    if usageTarget != null
    then
      (
        super.new(
          title=title,
          targets=[
            totalTarget { refId: 'TOTAL' },
            usageTarget { refId: 'USAGE' },
          ],
          description=description,
        )
        + $.withUsageTableCommonMixin()
        + table.queryOptions.withTransformationsMixin(
          [
            {
              id: 'groupBy',
              options: {
                fields: {
                  'Value #TOTAL': {
                    aggregations: [
                      'lastNotNull',
                    ],
                    operation: 'aggregate',
                  },
                  'Value #USAGE': {
                    aggregations: [
                      'lastNotNull',
                    ],
                    operation: 'aggregate',
                  },
                  [groupLabel]: {
                    aggregations: [],
                    operation: 'groupby',
                  },
                },
              },
            },
            {
              id: 'merge',
              options: {},
            },
            {
              id: 'calculateField',
              options: {
                alias: 'Available',
                binary: {
                  left: 'Value #TOTAL (lastNotNull)',
                  operator: '-',
                  reducer: 'sum',
                  right: 'Value #USAGE (lastNotNull)',
                },
                mode: 'binary',
                reduce: {
                  reducer: 'sum',
                },
              },
            },
            {
              id: 'calculateField',
              options: {
                alias: 'Used, %',
                binary: {
                  left: 'Value #USAGE (lastNotNull)',
                  operator: '/',
                  reducer: 'sum',
                  right: 'Value #TOTAL (lastNotNull)',
                },
                mode: 'binary',
                reduce: {
                  reducer: 'sum',
                },
              },
            },
            {
              id: 'organize',
              options: {
                excludeByName: {},
                indexByName: {
                  [groupLabel]: 0,
                  'Value #TOTAL (lastNotNull)': 1,
                  Available: 2,
                  'Value #USAGE (lastNotNull)': 3,
                  'Used, %': 4,
                },
                renameByName: {
                  'Value #TOTAL (lastNotNull)': 'Size',
                  'Value #USAGE (lastNotNull)': 'Used',
                  [groupLabel]: 'Mounted on',
                },
              },
            },
            self.transformations.sortBy('Mounted on'),
          ]
        )
      )
    else if freeTarget != null && usageTarget == null
    then
      (
        super.new(
          title=title,
          targets=[
            totalTarget { refId: 'TOTAL' },
            freeTarget { refId: 'FREE' },
          ],
          description=description,
        )
        + $.withUsageTableCommonMixin()
        + table.queryOptions.withTransformationsMixin(
          [
            {
              id: 'groupBy',
              options: {
                fields: {
                  'Value #TOTAL': {
                    aggregations: [
                      'lastNotNull',
                    ],
                    operation: 'aggregate',
                  },
                  'Value #FREE': {
                    aggregations: [
                      'lastNotNull',
                    ],
                    operation: 'aggregate',
                  },
                  [groupLabel]: {
                    aggregations: [],
                    operation: 'groupby',
                  },
                },
              },
            },
            {
              id: 'merge',
              options: {},
            },
            {
              id: 'calculateField',
              options: {
                alias: 'Used',
                binary: {
                  left: 'Value #TOTAL (lastNotNull)',
                  operator: '-',
                  reducer: 'sum',
                  right: 'Value #FREE (lastNotNull)',
                },
                mode: 'binary',
                reduce: {
                  reducer: 'sum',
                },
              },
            },
            {
              id: 'calculateField',
              options: {
                alias: 'Used, %',
                binary: {
                  left: 'Used',
                  operator: '/',
                  reducer: 'sum',
                  right: 'Value #TOTAL (lastNotNull)',
                },
                mode: 'binary',
                reduce: {
                  reducer: 'sum',
                },
              },
            },
            {
              id: 'organize',
              options: {
                excludeByName: {},
                indexByName: {
                  [groupLabel]: 0,
                  'Value #TOTAL (lastNotNull)': 1,
                  'Value #FREE (lastNotNull)': 2,
                  Used: 3,
                  'Used, %': 4,
                },
                renameByName: {
                  'Value #TOTAL (lastNotNull)': 'Size',
                  'Value #FREE (lastNotNull)': 'Available',
                  [groupLabel]: 'Mounted on',
                },
              },
            },
            self.transformations.sortBy('Mounted on'),
          ]
        )
      )
    else {},

  withUsageTableCommonMixin():
    table.standardOptions.thresholds.withSteps(
      [
        table.thresholdStep.withColor('light-blue')
        + table.thresholdStep.withValue(null),
        table.thresholdStep.withColor('light-yellow')
        + table.thresholdStep.withValue(0.8),
        table.thresholdStep.withColor('light-red')
        + table.thresholdStep.withValue(0.9),
      ]
    )

    + table.standardOptions.withOverrides([
      fieldOverride.byName.new('Mounted on')
      + fieldOverride.byName.withProperty('custom.width', '260'),
      fieldOverride.byName.new('Size')
      + fieldOverride.byName.withProperty('custom.width', '80'),
      fieldOverride.byName.new('Used')
      + fieldOverride.byName.withProperty('custom.width', '80'),
      fieldOverride.byName.new('Available')
      + fieldOverride.byName.withProperty('custom.width', '80'),
      fieldOverride.byName.new('Used, %')
      + fieldOverride.byName.withProperty('custom.displayMode', 'basic')
      + fieldOverride.byName.withPropertiesFromOptions(
        table.standardOptions.withMax(1)
        + table.standardOptions.withMin(0)
        + table.standardOptions.withUnit('percentunit')
      ),
    ])
    + table.standardOptions.withUnit('bytes'),
}
