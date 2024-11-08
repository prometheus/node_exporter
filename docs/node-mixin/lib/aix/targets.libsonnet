local g = import '../g.libsonnet';
local prometheusQuery = g.query.prometheus;
local lokiQuery = g.query.loki;

{
  new(this): {
    local variables = this.grafana.variables.main,
    local config = this.config,
    local prometheusDatasource = '${' + variables.datasources.prometheus.name + '}',
    local lokiDatasource = '${' + variables.datasources.loki.name + '}',

    // override memory targets (other metrics in macos)
    memory+: {
      memoryTotalBytes:
        prometheusQuery.new(
          prometheusDatasource,
          'node_memory_total_bytes{%(queriesSelector)s}' % variables
        )
        + prometheusQuery.withLegendFormat('Physical memory'),

      memoryUsedBytes:
        prometheusQuery.new(
          prometheusDatasource,
          |||
            (
              node_memory_total_bytes{%(queriesSelector)s} -
              node_memory_available_bytes{%(queriesSelector)s}
            )
          ||| % variables
        )
        + prometheusQuery.withLegendFormat('Memory used'),

      memoryUsagePercent:
        prometheusQuery.new(
          prometheusDatasource,
          |||
            (
              (
                node_memory_total_bytes{%(queriesSelector)s} -
                node_memory_available_bytes{%(queriesSelector)s}
              ) 
              /avg(node_memory_total_bytes{%(queriesSelector)s})
            ) * 100
          |||
          % variables,
        ),
      memorySwapTotal:
        prometheusQuery.new(
          prometheusDatasource,
          'node_memory_swap_total_bytes{%(queriesSelector)s}' % variables
        ),

      memorySwapUsedBytes:
        prometheusQuery.new(
          prometheusDatasource,
          'node_memory_swap_used_bytes{%(queriesSelector)s}' % variables
        )
        + prometheusQuery.withLegendFormat('Swap used'),
    },
  },
}
