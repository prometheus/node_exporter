// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.query.elasticsearch', name: 'elasticsearch' },
  '#withDatasource': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: "For mixed data sources the selected datasource is on the query level.\nFor non mixed scenarios this is undefined.\nTODO find a better way to do this ^ that's friendly to schema\nTODO this shouldn't be unknown but DataSourceRef | null" } },
  withDatasource(value): { datasource: value },
  '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'true if query is disabled (ie should not be returned to the dashboard)\nNote this does not always imply that the query should not be executed since\nthe results from a hidden query may be used as the input to other queries (SSE etc)' } },
  withHide(value=true): { hide: value },
  '#withQueryType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Specify the query flavor\nTODO make this required and give it a default' } },
  withQueryType(value): { queryType: value },
  '#withRefId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'A unique identifier for the query within the list of targets.\nIn server side expressions, the refId is used as a variable name to identify results.\nBy default, the UI will assign A->Z; however setting meaningful names may be useful.' } },
  withRefId(value): { refId: value },
  '#withAlias': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Alias pattern' } },
  withAlias(value): { alias: value },
  '#withBucketAggs': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'List of bucket aggregations' } },
  withBucketAggs(value): { bucketAggs: (if std.isArray(value)
                                        then value
                                        else [value]) },
  '#withBucketAggsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'List of bucket aggregations' } },
  withBucketAggsMixin(value): { bucketAggs+: (if std.isArray(value)
                                              then value
                                              else [value]) },
  bucketAggs+:
    {
      '#': { help: '', name: 'bucketAggs' },
      DateHistogram+:
        {
          '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withId(value): { id: value },
          '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withSettings(value): { settings: value },
          '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withType(value): { type: value },
          '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withField(value): { field: value },
          '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withSettingsMixin(value): { settings+: value },
          settings+:
            {
              '#withInterval': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withInterval(value): { settings+: { interval: value } },
              '#withMinDocCount': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withMinDocCount(value): { settings+: { min_doc_count: value } },
              '#withOffset': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withOffset(value): { settings+: { offset: value } },
              '#withTimeZone': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withTimeZone(value): { settings+: { timeZone: value } },
              '#withTrimEdges': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withTrimEdges(value): { settings+: { trimEdges: value } },
            },
        },
      Histogram+:
        {
          '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withId(value): { id: value },
          '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withSettings(value): { settings: value },
          '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withType(value): { type: value },
          '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withField(value): { field: value },
          '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withSettingsMixin(value): { settings+: value },
          settings+:
            {
              '#withInterval': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withInterval(value): { settings+: { interval: value } },
              '#withMinDocCount': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withMinDocCount(value): { settings+: { min_doc_count: value } },
            },
        },
      Terms+:
        {
          '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withId(value): { id: value },
          '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withSettings(value): { settings: value },
          '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withType(value): { type: value },
          '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withField(value): { field: value },
          '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withSettingsMixin(value): { settings+: value },
          settings+:
            {
              '#withMinDocCount': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withMinDocCount(value): { settings+: { min_doc_count: value } },
              '#withMissing': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withMissing(value): { settings+: { missing: value } },
              '#withOrder': { 'function': { args: [{ default: null, enums: ['desc', 'asc'], name: 'value', type: 'string' }], help: '' } },
              withOrder(value): { settings+: { order: value } },
              '#withOrderBy': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withOrderBy(value): { settings+: { orderBy: value } },
              '#withSize': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withSize(value): { settings+: { size: value } },
            },
        },
      Filters+:
        {
          '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withId(value): { id: value },
          '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withSettings(value): { settings: value },
          '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withType(value): { type: value },
          '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withSettingsMixin(value): { settings+: value },
          settings+:
            {
              '#withFilters': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
              withFilters(value): { settings+: { filters: (if std.isArray(value)
                                                           then value
                                                           else [value]) } },
              '#withFiltersMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
              withFiltersMixin(value): { settings+: { filters+: (if std.isArray(value)
                                                                 then value
                                                                 else [value]) } },
              filters+:
                {
                  '#': { help: '', name: 'filters' },
                  '#withLabel': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withLabel(value): { label: value },
                  '#withQuery': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withQuery(value): { query: value },
                },
            },
        },
      GeoHashGrid+:
        {
          '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withId(value): { id: value },
          '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withSettings(value): { settings: value },
          '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withType(value): { type: value },
          '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withField(value): { field: value },
          '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withSettingsMixin(value): { settings+: value },
          settings+:
            {
              '#withPrecision': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withPrecision(value): { settings+: { precision: value } },
            },
        },
      Nested+:
        {
          '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withId(value): { id: value },
          '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withSettings(value): { settings: value },
          '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withType(value): { type: value },
          '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withField(value): { field: value },
          '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withSettingsMixin(value): { settings+: value },
        },
    },
  '#withMetrics': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'List of metric aggregations' } },
  withMetrics(value): { metrics: (if std.isArray(value)
                                  then value
                                  else [value]) },
  '#withMetricsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'List of metric aggregations' } },
  withMetricsMixin(value): { metrics+: (if std.isArray(value)
                                        then value
                                        else [value]) },
  metrics+:
    {
      '#': { help: '', name: 'metrics' },
      Count+:
        {
          '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withHide(value=true): { hide: value },
          '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withId(value): { id: value },
          '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withType(value): { type: value },
        },
      PipelineMetricAggregation+:
        {
          MovingAverage+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withPipelineAgg': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withPipelineAgg(value): { pipelineAgg: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
            },
          Derivative+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withPipelineAgg': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withPipelineAgg(value): { pipelineAgg: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withUnit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withUnit(value): { settings+: { unit: value } },
                },
            },
          CumulativeSum+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withPipelineAgg': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withPipelineAgg(value): { pipelineAgg: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withFormat': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withFormat(value): { settings+: { format: value } },
                },
            },
          BucketScript+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withPipelineVariables': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
              withPipelineVariables(value): { pipelineVariables: (if std.isArray(value)
                                                                  then value
                                                                  else [value]) },
              '#withPipelineVariablesMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
              withPipelineVariablesMixin(value): { pipelineVariables+: (if std.isArray(value)
                                                                        then value
                                                                        else [value]) },
              pipelineVariables+:
                {
                  '#': { help: '', name: 'pipelineVariables' },
                  '#withName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withName(value): { name: value },
                  '#withPipelineAgg': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withPipelineAgg(value): { pipelineAgg: value },
                },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withScript': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScript(value): { settings+: { script: value } },
                  '#withScriptMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScriptMixin(value): { settings+: { script+: value } },
                  script+:
                    {
                      '#withInline': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                      withInline(value): { settings+: { script+: { inline: value } } },
                    },
                },
            },
        },
      MetricAggregationWithSettings+:
        {
          BucketScript+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withPipelineVariables': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
              withPipelineVariables(value): { pipelineVariables: (if std.isArray(value)
                                                                  then value
                                                                  else [value]) },
              '#withPipelineVariablesMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
              withPipelineVariablesMixin(value): { pipelineVariables+: (if std.isArray(value)
                                                                        then value
                                                                        else [value]) },
              pipelineVariables+:
                {
                  '#': { help: '', name: 'pipelineVariables' },
                  '#withName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withName(value): { name: value },
                  '#withPipelineAgg': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withPipelineAgg(value): { pipelineAgg: value },
                },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withScript': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScript(value): { settings+: { script: value } },
                  '#withScriptMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScriptMixin(value): { settings+: { script+: value } },
                  script+:
                    {
                      '#withInline': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                      withInline(value): { settings+: { script+: { inline: value } } },
                    },
                },
            },
          CumulativeSum+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withPipelineAgg': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withPipelineAgg(value): { pipelineAgg: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withFormat': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withFormat(value): { settings+: { format: value } },
                },
            },
          Derivative+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withPipelineAgg': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withPipelineAgg(value): { pipelineAgg: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withUnit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withUnit(value): { settings+: { unit: value } },
                },
            },
          SerialDiff+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withPipelineAgg': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withPipelineAgg(value): { pipelineAgg: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withLag': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withLag(value): { settings+: { lag: value } },
                },
            },
          RawData+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withSize': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withSize(value): { settings+: { size: value } },
                },
            },
          RawDocument+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withSize': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withSize(value): { settings+: { size: value } },
                },
            },
          UniqueCount+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withMissing': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withMissing(value): { settings+: { missing: value } },
                  '#withPrecisionThreshold': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withPrecisionThreshold(value): { settings+: { precision_threshold: value } },
                },
            },
          Percentiles+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withScript': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScript(value): { settings+: { script: value } },
                  '#withScriptMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScriptMixin(value): { settings+: { script+: value } },
                  script+:
                    {
                      '#withInline': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                      withInline(value): { settings+: { script+: { inline: value } } },
                    },
                  '#withMissing': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withMissing(value): { settings+: { missing: value } },
                  '#withPercents': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
                  withPercents(value): { settings+: { percents: (if std.isArray(value)
                                                                 then value
                                                                 else [value]) } },
                  '#withPercentsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
                  withPercentsMixin(value): { settings+: { percents+: (if std.isArray(value)
                                                                       then value
                                                                       else [value]) } },
                },
            },
          ExtendedStats+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withScript': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScript(value): { settings+: { script: value } },
                  '#withScriptMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScriptMixin(value): { settings+: { script+: value } },
                  script+:
                    {
                      '#withInline': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                      withInline(value): { settings+: { script+: { inline: value } } },
                    },
                  '#withMissing': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withMissing(value): { settings+: { missing: value } },
                  '#withSigma': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withSigma(value): { settings+: { sigma: value } },
                },
              '#withMeta': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withMeta(value): { meta: value },
              '#withMetaMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withMetaMixin(value): { meta+: value },
            },
          Min+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withScript': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScript(value): { settings+: { script: value } },
                  '#withScriptMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScriptMixin(value): { settings+: { script+: value } },
                  script+:
                    {
                      '#withInline': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                      withInline(value): { settings+: { script+: { inline: value } } },
                    },
                  '#withMissing': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withMissing(value): { settings+: { missing: value } },
                },
            },
          Max+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withScript': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScript(value): { settings+: { script: value } },
                  '#withScriptMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScriptMixin(value): { settings+: { script+: value } },
                  script+:
                    {
                      '#withInline': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                      withInline(value): { settings+: { script+: { inline: value } } },
                    },
                  '#withMissing': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withMissing(value): { settings+: { missing: value } },
                },
            },
          Sum+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withScript': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScript(value): { settings+: { script: value } },
                  '#withScriptMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScriptMixin(value): { settings+: { script+: value } },
                  script+:
                    {
                      '#withInline': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                      withInline(value): { settings+: { script+: { inline: value } } },
                    },
                  '#withMissing': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withMissing(value): { settings+: { missing: value } },
                },
            },
          Average+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withMissing': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withMissing(value): { settings+: { missing: value } },
                  '#withScript': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScript(value): { settings+: { script: value } },
                  '#withScriptMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScriptMixin(value): { settings+: { script+: value } },
                  script+:
                    {
                      '#withInline': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                      withInline(value): { settings+: { script+: { inline: value } } },
                    },
                },
            },
          MovingAverage+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withPipelineAgg': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withPipelineAgg(value): { pipelineAgg: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
            },
          MovingFunction+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withPipelineAgg': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withPipelineAgg(value): { pipelineAgg: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withScript': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScript(value): { settings+: { script: value } },
                  '#withScriptMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withScriptMixin(value): { settings+: { script+: value } },
                  script+:
                    {
                      '#withInline': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                      withInline(value): { settings+: { script+: { inline: value } } },
                    },
                  '#withShift': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withShift(value): { settings+: { shift: value } },
                  '#withWindow': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withWindow(value): { settings+: { window: value } },
                },
            },
          Logs+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withLimit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withLimit(value): { settings+: { limit: value } },
                },
            },
          Rate+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withField(value): { field: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withMode': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withMode(value): { settings+: { mode: value } },
                  '#withUnit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withUnit(value): { settings+: { unit: value } },
                },
            },
          TopMetrics+:
            {
              '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withHide(value=true): { hide: value },
              '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value): { id: value },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { type: value },
              '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettings(value): { settings: value },
              '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
              withSettingsMixin(value): { settings+: value },
              settings+:
                {
                  '#withMetrics': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
                  withMetrics(value): { settings+: { metrics: (if std.isArray(value)
                                                               then value
                                                               else [value]) } },
                  '#withMetricsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
                  withMetricsMixin(value): { settings+: { metrics+: (if std.isArray(value)
                                                                     then value
                                                                     else [value]) } },
                  '#withOrder': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withOrder(value): { settings+: { order: value } },
                  '#withOrderBy': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withOrderBy(value): { settings+: { orderBy: value } },
                },
            },
        },
    },
  '#withQuery': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Lucene query' } },
  withQuery(value): { query: value },
  '#withTimeField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Name of time field' } },
  withTimeField(value): { timeField: value },
}
