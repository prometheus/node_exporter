local g = import '../../g.libsonnet';
local commonlib = import 'github.com/grafana/jsonnet-libs/common-lib/common/main.libsonnet';
local prometheusQuery = g.query.prometheus;


{
  new(this): {
    local variables = this.grafana.variables.useCluster,
    local config = this.config,
    local baseTargets = this.grafana.targets.useCluster,
    cpuUtilization:
      baseTargets.cpuUtilization
      + prometheusQuery.withExpr('sum by (%s) (%s)' % [this.config.clusterLabel, baseTargets.cpuUtilization.expr])
      + prometheusQuery.withLegendFormat('{{%s}}: Utilization' % this.config.clusterLabel),

    cpuSaturation:
      baseTargets.cpuSaturation
      + prometheusQuery.withExpr('sum by (%s) (%s)' % [this.config.clusterLabel, baseTargets.cpuSaturation.expr])
      + prometheusQuery.withLegendFormat('{{%s}}: Saturation' % this.config.clusterLabel),

    memoryUtilization:
      baseTargets.memoryUtilization
      + prometheusQuery.withExpr('sum by (%s) (%s)' % [this.config.clusterLabel, baseTargets.memoryUtilization.expr])
      + prometheusQuery.withLegendFormat('{{%s}}: Utilization' % this.config.clusterLabel),

    memorySaturation:
      baseTargets.memorySaturation
      + prometheusQuery.withExpr('sum by (%s) (%s)' % [this.config.clusterLabel, baseTargets.memorySaturation.expr])
      + prometheusQuery.withLegendFormat('{{%s}}: Major page fault operations' % this.config.clusterLabel),
    networkUtilizationReceive:
      baseTargets.networkUtilizationReceive
      + prometheusQuery.withExpr('sum by (%s) (%s)' % [this.config.clusterLabel, baseTargets.networkUtilizationReceive.expr])
      + prometheusQuery.withLegendFormat('{{%s}}: Receive' % this.config.clusterLabel),

    networkUtilizationTransmit:
      baseTargets.networkUtilizationTransmit
      + prometheusQuery.withExpr('sum by (%s) (%s)' % [this.config.clusterLabel, baseTargets.networkUtilizationTransmit.expr])
      + prometheusQuery.withLegendFormat('{{%s}}: Transmit' % this.config.clusterLabel),
    networkSaturationReceive:
      baseTargets.networkSaturationReceive
      + prometheusQuery.withExpr('sum by (%s) (%s)' % [this.config.clusterLabel, baseTargets.networkSaturationReceive.expr])
      + prometheusQuery.withLegendFormat('{{%s}}: Receive' % this.config.clusterLabel),
    networkSaturationTransmit:
      baseTargets.networkSaturationReceive
      + prometheusQuery.withExpr('sum by (%s) (%s)' % [this.config.clusterLabel, baseTargets.networkSaturationReceive.expr])
      + prometheusQuery.withLegendFormat('{{%s}}: Transmit' % this.config.clusterLabel),

    diskUtilization:
      baseTargets.diskUtilization
      + prometheusQuery.withExpr('sum by (%s) (%s)' % [this.config.clusterLabel, baseTargets.diskUtilization.expr])
      + prometheusQuery.withLegendFormat('{{%s}}' % this.config.clusterLabel),
    diskSaturation:
      baseTargets.diskSaturation
      + prometheusQuery.withExpr('sum by (%s) (%s)' % [this.config.clusterLabel, baseTargets.diskSaturation.expr])
      + prometheusQuery.withLegendFormat('{{%s}}' % this.config.clusterLabel),

    filesystemUtilization:
      baseTargets.filesystemUtilization
      + prometheusQuery.withExpr('sum by (%s) (%s)' % [this.config.clusterLabel, baseTargets.filesystemUtilization.expr])
      + prometheusQuery.withLegendFormat('{{%s}}' % this.config.clusterLabel),


  },
}
