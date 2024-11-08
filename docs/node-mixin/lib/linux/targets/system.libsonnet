local g = import '../../g.libsonnet';
local prometheusQuery = g.query.prometheus;
local lokiQuery = g.query.loki;

{
  new(this): {
    local variables = this.grafana.variables.main,
    local config = this.config,
    local prometheusDatasource = '${' + variables.datasources.prometheus.name + '}',
    local lokiDatasource = '${' + variables.datasources.loki.name + '}',
    uptimeQuery:: 'node_boot_time_seconds',
    uptime:
      prometheusQuery.new(
        prometheusDatasource,
        'time() - ' + self.uptimeQuery + '{%(queriesSelector)s}' % variables
      ),
    unameInfo:
      prometheusQuery.new(
        prometheusDatasource,
        'node_uname_info{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withFormat('table'),
    osInfo:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          node_os_info{%(queriesSelector)s}
        ||| % variables,
      )
      + prometheusQuery.withFormat('table'),
    osInfoCombined:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          node_uname_info{%(queriesSelector)s} 
          * on (%(groupLabels)s,%(instanceLabels)s)
          group_left(pretty_name)
          node_os_info{%(queriesSelector)s}
        ||| % variables {
          instanceLabels: std.join(',', this.config.instanceLabels),
          groupLabels: std.join(',', this.config.groupLabels),
        },
      )
      + prometheusQuery.withFormat('table'),

    osTimezone:  //timezone label
      prometheusQuery.new(
        prometheusDatasource,
        'node_time_zone_offset_seconds{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withFormat('table'),

    systemLoad1:
      prometheusQuery.new(
        prometheusDatasource,
        'node_load1{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('1m'),
    systemLoad5:
      prometheusQuery.new(
        prometheusDatasource,
        'node_load5{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('5m'),
    systemLoad15:
      prometheusQuery.new(
        prometheusDatasource,
        'node_load15{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('15m'),

    systemContextSwitches:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_context_switches_total{%(queriesSelector)s}[$__rate_interval])' % variables,
      )
      + prometheusQuery.withLegendFormat('Context switches'),

    systemInterrupts:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_intr_total{%(queriesSelector)s}[$__rate_interval])' % variables,
      )
      + prometheusQuery.withLegendFormat('Interrupts'),

    timeNtpStatus:
      prometheusQuery.new(
        prometheusDatasource,
        'node_timex_sync_status{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('NTP status'),

    timeOffset:
      prometheusQuery.new(
        prometheusDatasource,
        'node_timex_offset_seconds{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Time offset'),

    timeEstimatedError:
      prometheusQuery.new(
        prometheusDatasource,
        'node_timex_estimated_error_seconds{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Estimated error in seconds'),
    timeMaxError:
      prometheusQuery.new(
        prometheusDatasource,
        'node_timex_maxerror_seconds{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Maximum error in seconds'),

  },
}
