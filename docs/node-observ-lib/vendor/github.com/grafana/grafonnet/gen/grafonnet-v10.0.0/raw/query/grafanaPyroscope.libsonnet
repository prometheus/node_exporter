// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.query.grafanaPyroscope', name: 'grafanaPyroscope' },
  '#withDatasource': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: "For mixed data sources the selected datasource is on the query level.\nFor non mixed scenarios this is undefined.\nTODO find a better way to do this ^ that's friendly to schema\nTODO this shouldn't be unknown but DataSourceRef | null" } },
  withDatasource(value): { datasource: value },
  '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'true if query is disabled (ie should not be returned to the dashboard)\nNote this does not always imply that the query should not be executed since\nthe results from a hidden query may be used as the input to other queries (SSE etc)' } },
  withHide(value=true): { hide: value },
  '#withQueryType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Specify the query flavor\nTODO make this required and give it a default' } },
  withQueryType(value): { queryType: value },
  '#withRefId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'A unique identifier for the query within the list of targets.\nIn server side expressions, the refId is used as a variable name to identify results.\nBy default, the UI will assign A->Z; however setting meaningful names may be useful.' } },
  withRefId(value): { refId: value },
  '#withGroupBy': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Allows to group the results.' } },
  withGroupBy(value): { groupBy: (if std.isArray(value)
                                  then value
                                  else [value]) },
  '#withGroupByMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Allows to group the results.' } },
  withGroupByMixin(value): { groupBy+: (if std.isArray(value)
                                        then value
                                        else [value]) },
  '#withLabelSelector': { 'function': { args: [{ default: '{}', enums: null, name: 'value', type: 'string' }], help: 'Specifies the query label selectors.' } },
  withLabelSelector(value='{}'): { labelSelector: value },
  '#withMaxNodes': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'Sets the maximum number of nodes in the flamegraph.' } },
  withMaxNodes(value): { maxNodes: value },
  '#withProfileTypeId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Specifies the type of profile to query.' } },
  withProfileTypeId(value): { profileTypeId: value },
}
