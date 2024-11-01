local g = import '../g.libsonnet';
local prometheusQuery = g.query.prometheus;
local lokiQuery = g.query.loki;

{
  new(this): {
    local variables = this.grafana.variables.main,
    local config = this.config,
    local prometheusDatasource = '${' + variables.datasources.prometheus.name + '}',
    local lokiDatasource = '${' + variables.datasources.loki.name + '}',

    hardwareTemperature:
      prometheusQuery.new(
        prometheusDatasource,
        'node_hwmon_temp_celsius{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('{{chip}}/{{sensor}}'),

  },
}
