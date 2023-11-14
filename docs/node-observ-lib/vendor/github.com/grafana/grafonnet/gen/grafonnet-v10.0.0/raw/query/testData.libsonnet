// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.query.testData', name: 'testData' },
  '#withDatasource': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: "For mixed data sources the selected datasource is on the query level.\nFor non mixed scenarios this is undefined.\nTODO find a better way to do this ^ that's friendly to schema\nTODO this shouldn't be unknown but DataSourceRef | null" } },
  withDatasource(value): { datasource: value },
  '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'true if query is disabled (ie should not be returned to the dashboard)\nNote this does not always imply that the query should not be executed since\nthe results from a hidden query may be used as the input to other queries (SSE etc)' } },
  withHide(value=true): { hide: value },
  '#withQueryType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Specify the query flavor\nTODO make this required and give it a default' } },
  withQueryType(value): { queryType: value },
  '#withRefId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'A unique identifier for the query within the list of targets.\nIn server side expressions, the refId is used as a variable name to identify results.\nBy default, the UI will assign A->Z; however setting meaningful names may be useful.' } },
  withRefId(value): { refId: value },
  '#withAlias': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
  withAlias(value): { alias: value },
  '#withChannel': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
  withChannel(value): { channel: value },
  '#withCsvContent': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
  withCsvContent(value): { csvContent: value },
  '#withCsvFileName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
  withCsvFileName(value): { csvFileName: value },
  '#withCsvWave': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withCsvWave(value): { csvWave: (if std.isArray(value)
                                  then value
                                  else [value]) },
  '#withCsvWaveMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withCsvWaveMixin(value): { csvWave+: (if std.isArray(value)
                                        then value
                                        else [value]) },
  csvWave+:
    {
      '#': { help: '', name: 'csvWave' },
      '#withLabels': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withLabels(value): { labels: value },
      '#withName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withName(value): { name: value },
      '#withTimeStep': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
      withTimeStep(value): { timeStep: value },
      '#withValuesCSV': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withValuesCSV(value): { valuesCSV: value },
    },
  '#withErrorType': { 'function': { args: [{ default: null, enums: ['server_panic', 'frontend_exception', 'frontend_observable'], name: 'value', type: 'string' }], help: '' } },
  withErrorType(value): { errorType: value },
  '#withLabels': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
  withLabels(value): { labels: value },
  '#withLevelColumn': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
  withLevelColumn(value=true): { levelColumn: value },
  '#withLines': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
  withLines(value): { lines: value },
  '#withNodes': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withNodes(value): { nodes: value },
  '#withNodesMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withNodesMixin(value): { nodes+: value },
  nodes+:
    {
      '#withCount': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
      withCount(value): { nodes+: { count: value } },
      '#withType': { 'function': { args: [{ default: null, enums: ['random', 'response', 'random edges'], name: 'value', type: 'string' }], help: '' } },
      withType(value): { nodes+: { type: value } },
    },
  '#withPoints': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withPoints(value): { points: (if std.isArray(value)
                                then value
                                else [value]) },
  '#withPointsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withPointsMixin(value): { points+: (if std.isArray(value)
                                      then value
                                      else [value]) },
  '#withPulseWave': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withPulseWave(value): { pulseWave: value },
  '#withPulseWaveMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withPulseWaveMixin(value): { pulseWave+: value },
  pulseWave+:
    {
      '#withOffCount': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
      withOffCount(value): { pulseWave+: { offCount: value } },
      '#withOffValue': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
      withOffValue(value): { pulseWave+: { offValue: value } },
      '#withOnCount': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
      withOnCount(value): { pulseWave+: { onCount: value } },
      '#withOnValue': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
      withOnValue(value): { pulseWave+: { onValue: value } },
      '#withTimeStep': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
      withTimeStep(value): { pulseWave+: { timeStep: value } },
    },
  '#withRawFrameContent': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
  withRawFrameContent(value): { rawFrameContent: value },
  '#withScenarioId': { 'function': { args: [{ default: null, enums: ['random_walk', 'slow_query', 'random_walk_with_error', 'random_walk_table', 'exponential_heatmap_bucket_data', 'linear_heatmap_bucket_data', 'no_data_points', 'datapoints_outside_range', 'csv_metric_values', 'predictable_pulse', 'predictable_csv_wave', 'streaming_client', 'simulation', 'usa', 'live', 'grafana_api', 'arrow', 'annotations', 'table_static', 'server_error_500', 'logs', 'node_graph', 'flame_graph', 'raw_frame', 'csv_file', 'csv_content', 'trace', 'manual_entry', 'variables-query'], name: 'value', type: 'string' }], help: '' } },
  withScenarioId(value): { scenarioId: value },
  '#withSeriesCount': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
  withSeriesCount(value): { seriesCount: value },
  '#withSim': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withSim(value): { sim: value },
  '#withSimMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withSimMixin(value): { sim+: value },
  sim+:
    {
      '#withConfig': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withConfig(value): { sim+: { config: value } },
      '#withConfigMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withConfigMixin(value): { sim+: { config+: value } },
      '#withKey': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withKey(value): { sim+: { key: value } },
      '#withKeyMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withKeyMixin(value): { sim+: { key+: value } },
      key+:
        {
          '#withTick': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
          withTick(value): { sim+: { key+: { tick: value } } },
          '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withType(value): { sim+: { key+: { type: value } } },
          '#withUid': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withUid(value): { sim+: { key+: { uid: value } } },
        },
      '#withLast': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withLast(value=true): { sim+: { last: value } },
      '#withStream': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withStream(value=true): { sim+: { stream: value } },
    },
  '#withSpanCount': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
  withSpanCount(value): { spanCount: value },
  '#withStream': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withStream(value): { stream: value },
  '#withStreamMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withStreamMixin(value): { stream+: value },
  stream+:
    {
      '#withBands': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
      withBands(value): { stream+: { bands: value } },
      '#withNoise': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
      withNoise(value): { stream+: { noise: value } },
      '#withSpeed': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
      withSpeed(value): { stream+: { speed: value } },
      '#withSpread': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
      withSpread(value): { stream+: { spread: value } },
      '#withType': { 'function': { args: [{ default: null, enums: ['signal', 'logs', 'fetch'], name: 'value', type: 'string' }], help: '' } },
      withType(value): { stream+: { type: value } },
      '#withUrl': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withUrl(value): { stream+: { url: value } },
    },
  '#withStringInput': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
  withStringInput(value): { stringInput: value },
  '#withUsa': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withUsa(value): { usa: value },
  '#withUsaMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withUsaMixin(value): { usa+: value },
  usa+:
    {
      '#withFields': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withFields(value): { usa+: { fields: (if std.isArray(value)
                                            then value
                                            else [value]) } },
      '#withFieldsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withFieldsMixin(value): { usa+: { fields+: (if std.isArray(value)
                                                  then value
                                                  else [value]) } },
      '#withMode': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withMode(value): { usa+: { mode: value } },
      '#withPeriod': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withPeriod(value): { usa+: { period: value } },
      '#withStates': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withStates(value): { usa+: { states: (if std.isArray(value)
                                            then value
                                            else [value]) } },
      '#withStatesMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withStatesMixin(value): { usa+: { states+: (if std.isArray(value)
                                                  then value
                                                  else [value]) } },
    },
}
