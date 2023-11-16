local g = import './g.libsonnet';
local prometheusQuery = g.query.prometheus;
local lokiQuery = g.query.loki;

{
  new(this): {
    local variables = this.grafana.variables,
    local config = this.config,
    local prometheusDatasource = '${' + variables.datasources.prometheus.name + '}',
    local lokiDatasource = '${' + variables.datasources.loki.name + '}',

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
            node_memory_internal_bytes{%(queriesSelector)s} -
            node_memory_purgeable_bytes{%(queriesSelector)s} +
            node_memory_wired_bytes{%(queriesSelector)s} +
            node_memory_compressed_bytes{%(queriesSelector)s}
          )
        ||| % variables
      )
      + prometheusQuery.withLegendFormat('Memory used'),
    memoryAppBytes:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          (
              node_memory_internal_bytes{%(queriesSelector)s} -
              node_memory_purgeable_bytes{%(queriesSelector)s}
          )
        ||| % variables
      )
      + prometheusQuery.withLegendFormat('App memory'),
    memoryWiredBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_wired_bytes{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Wired memory'),
    memoryCompressedBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_compressed_bytes{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Compressed memory'),

    memoryUsagePercent:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          (
              (
                avg(node_memory_internal_bytes{%(queriesSelector)s}) -
                avg(node_memory_purgeable_bytes{%(queriesSelector)s}) +
                avg(node_memory_wired_bytes{%(queriesSelector)s}) +
                avg(node_memory_compressed_bytes{%(queriesSelector)s})
              ) /
              avg(node_memory_total_bytes{%(queriesSelector)s})
          )
          *
          100
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
}
