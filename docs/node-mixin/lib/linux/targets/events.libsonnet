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

    reboot:
      prometheusQuery.new(
        prometheusDatasource,
        self.uptimeQuery + '{%(queriesSelector)s}*1000 > $__from < $__to' % variables,
      ),

    serviceFailed:
      lokiQuery.new(
        lokiDatasource,
        '{%(queriesSelector)s, unit="init.scope"} |= "code=exited, status=1/FAILURE"' % variables
      ),
    // those events should be rare, so can be shown as annotations
    criticalEvents:
      lokiQuery.new(
        lokiDatasource,
        '{%(queriesSelector)s, transport="kernel", level="emerg"}' % variables
      ),
    memoryOOMkiller:
      prometheusQuery.new(
        prometheusDatasource,
        'increase(node_vmstat_oom_kill{%(queriesSelector)s}[$__interval:] offset -$__interval)' % variables,
      )
      + prometheusQuery.withLegendFormat('OOM killer invocations'),

    kernelUpdate:
      prometheusQuery.new(
        prometheusDatasource,
        expr=|||
          changes(
          sum by (%(instanceLabels)s) (
              group by (%(instanceLabels)s,release) (node_uname_info{%(queriesSelector)s})
              )
          [$__interval:1m] offset -$__interval) > 1
        ||| % variables { instanceLabels: std.join(',', this.config.instanceLabels) },
      ),

    // new interactive session in logs:
    sessionOpened:
      lokiQuery.new(
        lokiDatasource,
        '{%(queriesSelector)s, unit="systemd-logind.service"}|= "New session"' % variables
      ),
    sessionClosed:
      lokiQuery.new(
        lokiDatasource,
        '{%(queriesSelector)s, unit="systemd-logind.service"} |= "logged out"' % variables
      ),
  },
}
