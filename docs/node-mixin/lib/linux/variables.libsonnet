// variables.libsonnet
local g = import '../g.libsonnet';
local var = g.dashboard.variable;
local commonlib = import 'common-lib/common/main.libsonnet';
local utils = commonlib.utils;
local xtd = import 'github.com/jsonnet-libs/xtd/main.libsonnet';
{
  new(
    this
  ):
    {
      main:
        commonlib.variables.new(
          filteringSelector=this.config.filteringSelector,
          groupLabels=this.config.groupLabels,
          instanceLabels=this.config.instanceLabels,
          varMetric='node_uname_info',
          customAllValue=this.config.customAllValue,
          enableLokiLogs=this.config.enableLokiLogs,
        ),
      // used in USE cluster dashboard
      use:
        commonlib.variables.new(
          filteringSelector=this.config.filteringSelector,
          // drop clusterLabel from groupLabels:
          groupLabels=std.uniq(this.config.groupLabels + [this.config.clusterLabel]),
          instanceLabels=this.config.instanceLabels,
          varMetric='instance:node_cpu_utilisation:rate5m',
          customAllValue=this.config.customAllValue,
          enableLokiLogs=this.config.enableLokiLogs,
        ),
      useCluster:
        commonlib.variables.new(
          filteringSelector=this.config.filteringSelector,
          groupLabels=std.uniq(this.config.groupLabels + [this.config.clusterLabel]),
          instanceLabels=[],
          varMetric='instance:node_cpu_utilisation:rate5m',
          customAllValue=this.config.customAllValue,
          enableLokiLogs=this.config.enableLokiLogs,
        ),
    },

}
