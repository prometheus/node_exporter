local d = import 'github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet';

// match name/title to reduce diff in docs
local panelNames = {
  alertgroups: 'alertGroups',
  annolist: 'annotationsList',
  barchart: 'barChart',
  bargauge: 'barGauge',
  dashlist: 'dashboardList',
  nodeGraph: 'nodeGraph',
  piechart: 'pieChart',
  'state-timeline': 'stateTimeline',
  'status-history': 'statusHistory',
  timeseries: 'timeSeries',
  xychart: 'xyChart',
};

local getPanelName(type) =
  std.get(panelNames, type, type);

{
  '#new':: d.func.new(
    'Creates a new %s panel with a title.' % getPanelName(self.panelOptions.withType().type),
    args=[d.arg('title', d.T.string)]
  ),
  new(title):
    self.panelOptions.withTitle(title)
    + self.panelOptions.withType()
    + self.panelOptions.withPluginVersion()
    // Default to Mixed datasource so panels can be datasource agnostic, this
    // requires query targets to explicitly set datasource, which is a lot more
    // interesting from a reusability standpoint.
    + self.queryOptions.withDatasource('datasource', '-- Mixed --'),

  // Backwards compatible entries, ignored in docs
  link+: self.panelOptions.link + { '#':: { ignore: true } },
  thresholdStep+: self.standardOptions.threshold.step + { '#':: { ignore: true } },
  transformation+: self.queryOptions.transformation + { '#':: { ignore: true } },
  valueMapping+: self.standardOptions.mapping + { '#':: { ignore: true } },
  fieldOverride+: self.standardOptions.override + { '#':: { ignore: true } },

  '#gridPos': {},  // use withGridPos instead, a bit more concise.
  local gridPos = self.gridPos,
  panelOptions+: {
    '#withPluginVersion': {},

    '#withGridPos': d.func.new(
      |||
        `withGridPos` configures the height, width and xy coordinates of the panel. Also see `grafonnet.util.grid` for helper functions to calculate these fields.

        All arguments default to `null`, which means they will remain unchanged or unset.
      |||,
      args=[
        d.arg('h', d.T.number, default='null'),
        d.arg('w', d.T.number, default='null'),
        d.arg('x', d.T.number, default='null'),
        d.arg('y', d.T.number, default='null'),
      ]
    ),
    withGridPos(h=null, w=null, x=null, y=null):
      (if h != null then gridPos.withH(h) else {})
      + (if w != null then gridPos.withW(w) else {})
      + (if x != null then gridPos.withX(x) else {})
      + (if y != null then gridPos.withY(y) else {}),
  },

  '#datasource':: {},  // use withDatasource instead, bit more concise
  local datasource = self.datasource,
  queryOptions+: {
    '#withDatasource':: d.func.new(
      |||
        `withDatasource` sets the datasource for all queries in a panel.

        The default datasource for a panel is set to 'Mixed datasource' so panels can be datasource agnostic, which is a lot more interesting from a reusability standpoint. Note that this requires query targets to explicitly set datasource for the same reason.
      |||,
      args=[
        d.arg('type', d.T.string),
        d.arg('uid', d.T.string),
      ]
    ),
    withDatasource(type, uid):
      datasource.withType(type)
      + datasource.withUid(uid),
  },

  standardOptions+: {
    threshold+: { step+: { '#':: d.package.newSub('threshold.step', '') } },

    local overrides = super.override,
    local commonOverrideFunctions = {
      '#new':: d.fn(
        '`new` creates a new override of type `%s`.' % self.type,
        args=[
          d.arg('value', d.T.string),
        ]
      ),
      new(value):
        overrides.matcher.withId(self.type)
        + overrides.matcher.withOptions(value),

      '#withProperty':: d.fn(
        |||
          `withProperty` adds a property that needs to be overridden. This function can
          be called multiple time, adding more properties.
        |||,
        args=[
          d.arg('id', d.T.string),
          d.arg('value', d.T.any),
        ]
      ),
      withProperty(id, value):
        overrides.withPropertiesMixin([
          overrides.properties.withId(id)
          + overrides.properties.withValue(value),
        ]),

      '#withPropertiesFromOptions':: d.fn(
        |||
          `withPropertiesFromOptions` takes an object with properties that need to be
          overridden. See example code above.
        |||,
        args=[
          d.arg('options', d.T.object),
        ]
      ),
      withPropertiesFromOptions(options):
        local infunc(input, path=[]) =
          std.foldl(
            function(acc, p)
              acc + (
                if p == 'custom'
                then infunc(input[p], path=path + [p])
                else
                  overrides.withPropertiesMixin([
                    overrides.properties.withId(std.join('.', path + [p]))
                    + overrides.properties.withValue(input[p]),
                  ])
              ),
            std.objectFields(input),
            {}
          );
        infunc(options.fieldConfig.defaults),
    },

    override:
      {
        '#':: d.package.newSub(
          'override',
          |||
            Overrides allow you to customize visualization settings for specific fields or
            series. This is accomplished by adding an override rule that targets
            a particular set of fields and that can each define multiple options.

            ```jsonnet
            override.byType.new('number')
            + override.byType.withPropertiesFromOptions(
              panel.standardOptions.withDecimals(2)
              + panel.standardOptions.withUnit('s')
            )
            ```
          |||
        ),
        byName: commonOverrideFunctions + { type:: 'byName' },
        byRegexp: commonOverrideFunctions + { type:: 'byRegexp' },
        byType: commonOverrideFunctions + { type:: 'byType' },
        byQuery: commonOverrideFunctions + { type:: 'byQuery' },
        // TODO: byValue takes more complex `options` than string
        byValue: commonOverrideFunctions + { type:: 'byValue' },
      },
  },
}
