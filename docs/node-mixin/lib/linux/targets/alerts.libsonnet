local g = import '../../g.libsonnet';
local prometheusQuery = g.query.prometheus;
local lokiQuery = g.query.loki;

{
  new(this): {
    local variables = this.grafana.variables.main,
    local config = this.config,
    local prometheusDatasource = '${' + variables.datasources.prometheus.name + '}',
    local lokiDatasource = '${' + variables.datasources.loki.name + '}',

    alertsCritical:
      prometheusQuery.new(
        prometheusDatasource,
        'count by (%(instanceLabels)s) (max_over_time(ALERTS{%(queriesSelector)s, alertstate="firing", severity="critical"}[1m])) * group by (%(instanceLabels)s) (node_uname_info{%(queriesSelector)s})' % variables { instanceLabels: std.join(',', this.config.instanceLabels) },
      ),
    alertsWarning:
      prometheusQuery.new(
        prometheusDatasource,
        'count by (%(instanceLabels)s) (max_over_time(ALERTS{%(queriesSelector)s, alertstate="firing", severity="warning"}[1m])) * group by (%(instanceLabels)s) (node_uname_info{%(queriesSelector)s})' % variables { instanceLabels: std.join(',', this.config.instanceLabels) },
      ),


  },
}
