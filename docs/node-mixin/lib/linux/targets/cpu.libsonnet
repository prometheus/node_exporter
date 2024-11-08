local g = import '../../g.libsonnet';
local prometheusQuery = g.query.prometheus;
local lokiQuery = g.query.loki;

{
  new(this): {
    local variables = this.grafana.variables.main,
    local config = this.config,
    local prometheusDatasource = '${' + variables.datasources.prometheus.name + '}',
    local lokiDatasource = '${' + variables.datasources.loki.name + '}',

    cpuCount:
      prometheusQuery.new(
        prometheusDatasource,
        'count without (cpu) (node_cpu_seconds_total{%(queriesSelector)s, mode="idle"})' % variables
      )
      + prometheusQuery.withLegendFormat('Cores'),
    cpuUsage:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          (((count by (%(instanceLabels)s) (count(node_cpu_seconds_total{%(queriesSelector)s}) by (cpu, %(instanceLabels)s))) 
          - 
          avg by (%(instanceLabels)s) (sum by (%(instanceLabels)s, mode)(irate(node_cpu_seconds_total{mode='idle',%(queriesSelector)s}[$__rate_interval])))) * 100) 
          / 
          count by(%(instanceLabels)s) (count(node_cpu_seconds_total{%(queriesSelector)s}) by (cpu, %(instanceLabels)s))
        ||| % variables { instanceLabels: std.join(',', this.config.instanceLabels) },
      )
      + prometheusQuery.withLegendFormat('CPU usage'),
    cpuUsagePerCore:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          (
            (1 - sum without (mode) (rate(node_cpu_seconds_total{%(queriesSelector)s, mode=~"idle|iowait|steal"}[$__rate_interval])))
          / ignoring(cpu) group_left
            count without (cpu, mode) (node_cpu_seconds_total{%(queriesSelector)s, mode="idle"})
          ) * 100
        ||| % variables,
      )
      + prometheusQuery.withLegendFormat('CPU {{cpu}}'),
    cpuUsageByMode:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          sum by(%(instanceLabels)s, mode) (irate(node_cpu_seconds_total{%(queriesSelector)s}[$__rate_interval])) 
          / on(%(instanceLabels)s) 
          group_left sum by (%(instanceLabels)s)((irate(node_cpu_seconds_total{%(queriesSelector)s}[$__rate_interval]))) * 100
        ||| % variables { instanceLabels: std.join(',', this.config.instanceLabels) },
      )
      + prometheusQuery.withLegendFormat('{{ mode }}'),
  },
}
