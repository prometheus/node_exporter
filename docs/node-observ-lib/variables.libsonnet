// variables.libsonnet
local g = import './g.libsonnet';
local var = g.dashboard.variable;
local commonlib = import 'common-lib/common/main.libsonnet';
local utils = commonlib.utils;

{
  new(
    this
  ): {
       local filteringSelector = this.config.filteringSelector,
       local groupLabels = this.config.groupLabels,
       local instanceLabels = this.config.instanceLabels,
       local root = self,
       local varMetric = 'node_uname_info',
       local variablesFromLabels(groupLabels, instanceLabels, filteringSelector, multiInstance=true) =
         local chainVarProto(index, chainVar) =
           var.query.new(chainVar.label)
           + var.query.withDatasourceFromVariable(root.datasources.prometheus)
           + var.query.queryTypes.withLabelValues(
             chainVar.label,
             '%s{%s}' % [varMetric, chainVar.chainSelector],
           )
           + var.query.generalOptions.withLabel(utils.toSentenceCase(chainVar.label))
           + var.query.selectionOptions.withIncludeAll(
             value=if (!multiInstance && std.member(instanceLabels, chainVar.label)) then false else true,
             customAllValue='.+'
           )
           + var.query.selectionOptions.withMulti(
             if (!multiInstance && std.member(instanceLabels, chainVar.label)) then false else true,
           )
           + var.query.refresh.onTime()
           + var.query.withSort(
             i=1,
             type='alphabetical',
             asc=true,
             caseInsensitive=false
           );
         std.mapWithIndex(chainVarProto, utils.chainLabels(groupLabels + instanceLabels, [filteringSelector])),
       datasources: {
         prometheus:
           var.datasource.new('datasource', 'prometheus')
           + var.datasource.generalOptions.withLabel('Data source')
           + var.datasource.withRegex(''),
         loki:
           var.datasource.new('loki_datasource', 'loki')
           + var.datasource.generalOptions.withLabel('Loki data source')
           + var.datasource.withRegex('')
           + var.datasource.generalOptions.showOnDashboard.withNothing(),
       },
       // Use on dashboards where multiple entities can be selected, like fleet dashboards
       multiInstance:
         [root.datasources.prometheus]
         + variablesFromLabels(groupLabels, instanceLabels, filteringSelector),
       // Use on dashboards where only single entity can be selected
       singleInstance:
         [root.datasources.prometheus]
         + variablesFromLabels(groupLabels, instanceLabels, filteringSelector, multiInstance=false),

       queriesSelector:
         '%s,%s' % [
           filteringSelector,
           utils.labelsToPromQLSelector(groupLabels + instanceLabels),
         ],
     }
     + if this.config.enableLokiLogs then self.withLokiLogs(this) else {},
  withLokiLogs(this): {
    multiInstance+: [this.grafana.variables.datasources.loki],
    singleInstance+: [this.grafana.variables.datasources.loki],
  },
}
