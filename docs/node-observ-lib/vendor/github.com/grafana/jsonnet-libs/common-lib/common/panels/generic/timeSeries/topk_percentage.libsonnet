local g = import '../../../g.libsonnet';
local base = import './base.libsonnet';
local generic = import './main.libsonnet';
local timeSeries = g.panel.timeSeries;
local fieldOverride = g.panel.timeSeries.fieldOverride;
local fieldConfig = g.panel.timeSeries.fieldConfig;
local standardOptions = g.panel.timeSeries.standardOptions;
// Style to display Top K metrics that can go from 0 to 100%.
// It constructs mean baseline automatically.
base {
  new(
    title,
    target,
    topk=25,
    instanceLabels,
    drillDownDashboardUid,
    description='Top %s' % topk
  ):

    local topTarget = target
                      { expr: 'topk(' + topk + ',' + target.expr + ')' }
                      + g.query.prometheus.withLegendFormat(
                        std.join(': ', std.map(function(l) '{{' + l + '}}', instanceLabels))
                      );
    local meanTarget = target
                       { expr: 'avg(' + target.expr + ')' }
                       + g.query.prometheus.withLegendFormat('Mean');
    super.new(title, targets=[topTarget, meanTarget], description=description)
    + self.withDataLink(instanceLabels, drillDownDashboardUid),
  withDataLink(instanceLabels, drillDownDashboardUid):
    standardOptions.withLinks(
      {
        url: 'd/' + drillDownDashboardUid + '?' + std.join('&', std.map(function(l) 'var-%s=${__field.labels.%s}' % [l, l], instanceLabels)) + '&${__url_time_range}',
        title: 'Drill down to this instance',
      }
    ),
  stylize(allLayers=true):
    (if allLayers then super.stylize() else {})
    + generic.percentage.stylize(allLayers=false)
    + fieldConfig.defaults.custom.withFillOpacity(1)
    + fieldConfig.defaults.custom.withLineWidth(1)
    + timeSeries.options.legend.withDisplayMode('table')
    + timeSeries.options.legend.withPlacement('right')
    + timeSeries.options.legend.withCalcsMixin([
      'mean',
      'max',
      'lastNotNull',
    ])
    + timeSeries.standardOptions.withOverrides(
      fieldOverride.byName.new('Mean')
      + fieldOverride.byName.withPropertiesFromOptions(
        fieldConfig.defaults.custom.withLineStyleMixin(
          {
            fill: 'dash',
            dash: [10, 10],
          }
        )
        + fieldConfig.defaults.custom.withFillOpacity(0)
        + timeSeries.standardOptions.color.withMode('fixed')
        + timeSeries.standardOptions.color.withFixedColor('light-purple'),
      )
    ),
}
