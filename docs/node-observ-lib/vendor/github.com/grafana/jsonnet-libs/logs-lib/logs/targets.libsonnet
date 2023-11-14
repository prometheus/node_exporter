local utils = import '../utils.libsonnet';
local g = import './g.libsonnet';
local lokiQuery = g.query.loki;
function(
  variables,
  formatParser,
  logsVolumeGroupBy,
  extraFilters,
) {
  formatParser:: if formatParser != null then '| %s | __error__=``' % formatParser else '',
  logsTarget::
    lokiQuery.new(
      datasource='${' + variables.datasource.name + '}',
      expr=|||
        {%s} 
        |~ "$regex_search"
        %s
        %s
      ||| % [
        variables.queriesSelector,
        self.formatParser,
        extraFilters,
      ]
    ),

  logsVolumeTarget::
    lokiQuery.new(
      datasource='${' + variables.datasource.name + '}',
      expr=|||
        sum by (%s) (count_over_time({%s}
        |~ "$regex_search"
        %s
        [$__interval]))
      ||| % [
        logsVolumeGroupBy,
        variables.queriesSelector,
        self.formatParser,
      ]
    )
    + lokiQuery.withLegendFormat('{{ %s }}' % logsVolumeGroupBy),
}
