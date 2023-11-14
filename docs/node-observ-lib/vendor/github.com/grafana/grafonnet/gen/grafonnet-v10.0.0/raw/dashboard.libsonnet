// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.dashboard', name: 'dashboard' },
  '#withAnnotations': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO -- should not be a public interface on its own, but required for Veneer' } },
  withAnnotations(value): { annotations: value },
  '#withAnnotationsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO -- should not be a public interface on its own, but required for Veneer' } },
  withAnnotationsMixin(value): { annotations+: value },
  annotations+:
    {
      '#withList': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withList(value): { annotations+: { list: (if std.isArray(value)
                                                then value
                                                else [value]) } },
      '#withListMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withListMixin(value): { annotations+: { list+: (if std.isArray(value)
                                                      then value
                                                      else [value]) } },
      list+:
        {
          '#': { help: '', name: 'list' },
          '#withDatasource': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO: Should be DataSourceRef' } },
          withDatasource(value): { datasource: value },
          '#withDatasourceMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO: Should be DataSourceRef' } },
          withDatasourceMixin(value): { datasource+: value },
          datasource+:
            {
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withType(value): { datasource+: { type: value } },
              '#withUid': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withUid(value): { datasource+: { uid: value } },
            },
          '#withEnable': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'When enabled the annotation query is issued with every dashboard refresh' } },
          withEnable(value=true): { enable: value },
          '#withFilter': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withFilter(value): { filter: value },
          '#withFilterMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withFilterMixin(value): { filter+: value },
          filter+:
            {
              '#withExclude': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Should the specified panels be included or excluded' } },
              withExclude(value=true): { filter+: { exclude: value } },
              '#withIds': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Panel IDs that should be included or excluded' } },
              withIds(value): { filter+: { ids: (if std.isArray(value)
                                                 then value
                                                 else [value]) } },
              '#withIdsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Panel IDs that should be included or excluded' } },
              withIdsMixin(value): { filter+: { ids+: (if std.isArray(value)
                                                       then value
                                                       else [value]) } },
            },
          '#withHide': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Annotation queries can be toggled on or off at the top of the dashboard.\nWhen hide is true, the toggle is not shown in the dashboard.' } },
          withHide(value=true): { hide: value },
          '#withIconColor': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Color to use for the annotation event markers' } },
          withIconColor(value): { iconColor: value },
          '#withName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Name of annotation.' } },
          withName(value): { name: value },
          '#withTarget': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO: this should be a regular DataQuery that depends on the selected dashboard\nthese match the properties of the "grafana" datasouce that is default in most dashboards' } },
          withTarget(value): { target: value },
          '#withTargetMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO: this should be a regular DataQuery that depends on the selected dashboard\nthese match the properties of the "grafana" datasouce that is default in most dashboards' } },
          withTargetMixin(value): { target+: value },
          target+:
            {
              '#withLimit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'Only required/valid for the grafana datasource...\nbut code+tests is already depending on it so hard to change' } },
              withLimit(value): { target+: { limit: value } },
              '#withMatchAny': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Only required/valid for the grafana datasource...\nbut code+tests is already depending on it so hard to change' } },
              withMatchAny(value=true): { target+: { matchAny: value } },
              '#withTags': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Only required/valid for the grafana datasource...\nbut code+tests is already depending on it so hard to change' } },
              withTags(value): { target+: { tags: (if std.isArray(value)
                                                   then value
                                                   else [value]) } },
              '#withTagsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Only required/valid for the grafana datasource...\nbut code+tests is already depending on it so hard to change' } },
              withTagsMixin(value): { target+: { tags+: (if std.isArray(value)
                                                         then value
                                                         else [value]) } },
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Only required/valid for the grafana datasource...\nbut code+tests is already depending on it so hard to change' } },
              withType(value): { target+: { type: value } },
            },
          '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO -- this should not exist here, it is based on the --grafana-- datasource' } },
          withType(value): { type: value },
        },
    },
  '#withDescription': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Description of dashboard.' } },
  withDescription(value): { description: value },
  '#withEditable': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Whether a dashboard is editable or not.' } },
  withEditable(value=true): { editable: value },
  '#withFiscalYearStartMonth': { 'function': { args: [{ default: 0, enums: null, name: 'value', type: 'integer' }], help: 'The month that the fiscal year starts on.  0 = January, 11 = December' } },
  withFiscalYearStartMonth(value=0): { fiscalYearStartMonth: value },
  '#withGnetId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'For dashboards imported from the https://grafana.com/grafana/dashboards/ portal' } },
  withGnetId(value): { gnetId: value },
  '#withGraphTooltip': { 'function': { args: [{ default: 0, enums: [0, 1, 2], name: 'value', type: 'integer' }], help: '0 for no shared crosshair or tooltip (default).\n1 for shared crosshair.\n2 for shared crosshair AND shared tooltip.' } },
  withGraphTooltip(value=0): { graphTooltip: value },
  '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'Unique numeric identifier for the dashboard.\nTODO must isolate or remove identifiers local to a Grafana instance...?' } },
  withId(value): { id: value },
  '#withLinks': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'TODO docs' } },
  withLinks(value): { links: (if std.isArray(value)
                              then value
                              else [value]) },
  '#withLinksMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'TODO docs' } },
  withLinksMixin(value): { links+: (if std.isArray(value)
                                    then value
                                    else [value]) },
  links+:
    {
      '#': { help: '', name: 'links' },
      '#withAsDropdown': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withAsDropdown(value=true): { asDropdown: value },
      '#withIcon': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withIcon(value): { icon: value },
      '#withIncludeVars': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withIncludeVars(value=true): { includeVars: value },
      '#withKeepTime': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withKeepTime(value=true): { keepTime: value },
      '#withTags': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withTags(value): { tags: (if std.isArray(value)
                                then value
                                else [value]) },
      '#withTagsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withTagsMixin(value): { tags+: (if std.isArray(value)
                                      then value
                                      else [value]) },
      '#withTargetBlank': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withTargetBlank(value=true): { targetBlank: value },
      '#withTitle': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withTitle(value): { title: value },
      '#withTooltip': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withTooltip(value): { tooltip: value },
      '#withType': { 'function': { args: [{ default: null, enums: ['link', 'dashboards'], name: 'value', type: 'string' }], help: 'TODO docs' } },
      withType(value): { type: value },
      '#withUrl': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withUrl(value): { url: value },
    },
  '#withLiveNow': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'When set to true, the dashboard will redraw panels at an interval matching the pixel width.\nThis will keep data "moving left" regardless of the query refresh rate.  This setting helps\navoid dashboards presenting stale live data' } },
  withLiveNow(value=true): { liveNow: value },
  '#withPanels': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withPanels(value): { panels: (if std.isArray(value)
                                then value
                                else [value]) },
  '#withPanelsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withPanelsMixin(value): { panels+: (if std.isArray(value)
                                      then value
                                      else [value]) },
  '#withRefresh': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Refresh rate of dashboard. Represented via interval string, e.g. "5s", "1m", "1h", "1d".' } },
  withRefresh(value): { refresh: value },
  '#withRefreshMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Refresh rate of dashboard. Represented via interval string, e.g. "5s", "1m", "1h", "1d".' } },
  withRefreshMixin(value): { refresh+: value },
  '#withRevision': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'This property should only be used in dashboards defined by plugins.  It is a quick check\nto see if the version has changed since the last time.  Unclear why using the version property\nis insufficient.' } },
  withRevision(value): { revision: value },
  '#withSchemaVersion': { 'function': { args: [{ default: 36, enums: null, name: 'value', type: 'integer' }], help: "Version of the JSON schema, incremented each time a Grafana update brings\nchanges to said schema.\nTODO this is the existing schema numbering system. It will be replaced by Thema's themaVersion" } },
  withSchemaVersion(value=36): { schemaVersion: value },
  '#withSnapshot': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
  withSnapshot(value): { snapshot: value },
  '#withSnapshotMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
  withSnapshotMixin(value): { snapshot+: value },
  snapshot+:
    {
      '#withCreated': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO docs' } },
      withCreated(value): { snapshot+: { created: value } },
      '#withExpires': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO docs' } },
      withExpires(value): { snapshot+: { expires: value } },
      '#withExternal': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'TODO docs' } },
      withExternal(value=true): { snapshot+: { external: value } },
      '#withExternalUrl': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO docs' } },
      withExternalUrl(value): { snapshot+: { externalUrl: value } },
      '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'TODO docs' } },
      withId(value): { snapshot+: { id: value } },
      '#withKey': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO docs' } },
      withKey(value): { snapshot+: { key: value } },
      '#withName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO docs' } },
      withName(value): { snapshot+: { name: value } },
      '#withOrgId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'TODO docs' } },
      withOrgId(value): { snapshot+: { orgId: value } },
      '#withUpdated': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO docs' } },
      withUpdated(value): { snapshot+: { updated: value } },
      '#withUrl': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO docs' } },
      withUrl(value): { snapshot+: { url: value } },
      '#withUserId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'TODO docs' } },
      withUserId(value): { snapshot+: { userId: value } },
    },
  '#withStyle': { 'function': { args: [{ default: 'dark', enums: ['dark', 'light'], name: 'value', type: 'string' }], help: 'Theme of dashboard.' } },
  withStyle(value='dark'): { style: value },
  '#withTags': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Tags associated with dashboard.' } },
  withTags(value): { tags: (if std.isArray(value)
                            then value
                            else [value]) },
  '#withTagsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Tags associated with dashboard.' } },
  withTagsMixin(value): { tags+: (if std.isArray(value)
                                  then value
                                  else [value]) },
  '#withTemplating': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
  withTemplating(value): { templating: value },
  '#withTemplatingMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
  withTemplatingMixin(value): { templating+: value },
  templating+:
    {
      '#withList': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withList(value): { templating+: { list: (if std.isArray(value)
                                               then value
                                               else [value]) } },
      '#withListMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withListMixin(value): { templating+: { list+: (if std.isArray(value)
                                                     then value
                                                     else [value]) } },
      list+:
        {
          '#': { help: '', name: 'list' },
          '#withDatasource': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Ref to a DataSource instance' } },
          withDatasource(value): { datasource: value },
          '#withDatasourceMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Ref to a DataSource instance' } },
          withDatasourceMixin(value): { datasource+: value },
          datasource+:
            {
              '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'The plugin type-id' } },
              withType(value): { datasource+: { type: value } },
              '#withUid': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Specific datasource instance' } },
              withUid(value): { datasource+: { uid: value } },
            },
          '#withDescription': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withDescription(value): { description: value },
          '#withError': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withError(value): { 'error': value },
          '#withErrorMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withErrorMixin(value): { 'error'+: value },
          '#withGlobal': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withGlobal(value=true): { global: value },
          '#withHide': { 'function': { args: [{ default: null, enums: [0, 1, 2], name: 'value', type: 'integer' }], help: '' } },
          withHide(value): { hide: value },
          '#withId': { 'function': { args: [{ default: '00000000-0000-0000-0000-000000000000', enums: null, name: 'value', type: 'string' }], help: '' } },
          withId(value='00000000-0000-0000-0000-000000000000'): { id: value },
          '#withIndex': { 'function': { args: [{ default: -1, enums: null, name: 'value', type: 'integer' }], help: '' } },
          withIndex(value=-1): { index: value },
          '#withLabel': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withLabel(value): { label: value },
          '#withName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withName(value): { name: value },
          '#withQuery': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO: Move this into a separated QueryVariableModel type' } },
          withQuery(value): { query: value },
          '#withQueryMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO: Move this into a separated QueryVariableModel type' } },
          withQueryMixin(value): { query+: value },
          '#withRootStateKey': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withRootStateKey(value): { rootStateKey: value },
          '#withSkipUrlSync': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withSkipUrlSync(value=true): { skipUrlSync: value },
          '#withState': { 'function': { args: [{ default: null, enums: ['NotStarted', 'Loading', 'Streaming', 'Done', 'Error'], name: 'value', type: 'string' }], help: '' } },
          withState(value): { state: value },
          '#withType': { 'function': { args: [{ default: null, enums: ['query', 'adhoc', 'constant', 'datasource', 'interval', 'textbox', 'custom', 'system'], name: 'value', type: 'string' }], help: 'FROM: packages/grafana-data/src/types/templateVars.ts\nTODO docs\nTODO this implies some wider pattern/discriminated union, probably?' } },
          withType(value): { type: value },
        },
    },
  '#withTime': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Time range for dashboard, e.g. last 6 hours, last 7 days, etc' } },
  withTime(value): { time: value },
  '#withTimeMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Time range for dashboard, e.g. last 6 hours, last 7 days, etc' } },
  withTimeMixin(value): { time+: value },
  time+:
    {
      '#withFrom': { 'function': { args: [{ default: 'now-6h', enums: null, name: 'value', type: 'string' }], help: '' } },
      withFrom(value='now-6h'): { time+: { from: value } },
      '#withTo': { 'function': { args: [{ default: 'now', enums: null, name: 'value', type: 'string' }], help: '' } },
      withTo(value='now'): { time+: { to: value } },
    },
  '#withTimepicker': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs\nTODO this appears to be spread all over in the frontend. Concepts will likely need tidying in tandem with schema changes' } },
  withTimepicker(value): { timepicker: value },
  '#withTimepickerMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs\nTODO this appears to be spread all over in the frontend. Concepts will likely need tidying in tandem with schema changes' } },
  withTimepickerMixin(value): { timepicker+: value },
  timepicker+:
    {
      '#withCollapse': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Whether timepicker is collapsed or not.' } },
      withCollapse(value=true): { timepicker+: { collapse: value } },
      '#withEnable': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Whether timepicker is enabled or not.' } },
      withEnable(value=true): { timepicker+: { enable: value } },
      '#withHidden': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Whether timepicker is visible or not.' } },
      withHidden(value=true): { timepicker+: { hidden: value } },
      '#withRefreshIntervals': { 'function': { args: [{ default: ['5s', '10s', '30s', '1m', '5m', '15m', '30m', '1h', '2h', '1d'], enums: null, name: 'value', type: 'array' }], help: 'Selectable intervals for auto-refresh.' } },
      withRefreshIntervals(value): { timepicker+: { refresh_intervals: (if std.isArray(value)
                                                                        then value
                                                                        else [value]) } },
      '#withRefreshIntervalsMixin': { 'function': { args: [{ default: ['5s', '10s', '30s', '1m', '5m', '15m', '30m', '1h', '2h', '1d'], enums: null, name: 'value', type: 'array' }], help: 'Selectable intervals for auto-refresh.' } },
      withRefreshIntervalsMixin(value): { timepicker+: { refresh_intervals+: (if std.isArray(value)
                                                                              then value
                                                                              else [value]) } },
      '#withTimeOptions': { 'function': { args: [{ default: ['5m', '15m', '1h', '6h', '12h', '24h', '2d', '7d', '30d'], enums: null, name: 'value', type: 'array' }], help: 'TODO docs' } },
      withTimeOptions(value): { timepicker+: { time_options: (if std.isArray(value)
                                                              then value
                                                              else [value]) } },
      '#withTimeOptionsMixin': { 'function': { args: [{ default: ['5m', '15m', '1h', '6h', '12h', '24h', '2d', '7d', '30d'], enums: null, name: 'value', type: 'array' }], help: 'TODO docs' } },
      withTimeOptionsMixin(value): { timepicker+: { time_options+: (if std.isArray(value)
                                                                    then value
                                                                    else [value]) } },
    },
  '#withTimezone': { 'function': { args: [{ default: 'browser', enums: null, name: 'value', type: 'string' }], help: 'Timezone of dashboard. Accepts IANA TZDB zone ID or "browser" or "utc".' } },
  withTimezone(value='browser'): { timezone: value },
  '#withTitle': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Title of dashboard.' } },
  withTitle(value): { title: value },
  '#withUid': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Unique dashboard identifier that can be generated by anyone. string (8-40)' } },
  withUid(value): { uid: value },
  '#withVersion': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'Version of the dashboard, incremented each time the dashboard is updated.' } },
  withVersion(value): { version: value },
  '#withWeekStart': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO docs' } },
  withWeekStart(value): { weekStart: value },
}
