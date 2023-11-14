// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.panel', name: 'panel' },
  '#withDatasource': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'The datasource used in all targets.' } },
  withDatasource(value): { datasource: value },
  '#withDatasourceMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'The datasource used in all targets.' } },
  withDatasourceMixin(value): { datasource+: value },
  datasource+:
    {
      '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withType(value): { datasource+: { type: value } },
      '#withUid': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withUid(value): { datasource+: { uid: value } },
    },
  '#withDescription': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Description.' } },
  withDescription(value): { description: value },
  '#withFieldConfig': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withFieldConfig(value): { fieldConfig: value },
  '#withFieldConfigMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withFieldConfigMixin(value): { fieldConfig+: value },
  fieldConfig+:
    {
      '#withDefaults': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withDefaults(value): { fieldConfig+: { defaults: value } },
      '#withDefaultsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withDefaultsMixin(value): { fieldConfig+: { defaults+: value } },
      defaults+:
        {
          '#withColor': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
          withColor(value): { fieldConfig+: { defaults+: { color: value } } },
          '#withColorMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
          withColorMixin(value): { fieldConfig+: { defaults+: { color+: value } } },
          color+:
            {
              '#withFixedColor': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Stores the fixed color value if mode is fixed' } },
              withFixedColor(value): { fieldConfig+: { defaults+: { color+: { fixedColor: value } } } },
              '#withMode': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'The main color scheme mode' } },
              withMode(value): { fieldConfig+: { defaults+: { color+: { mode: value } } } },
              '#withSeriesBy': { 'function': { args: [{ default: null, enums: ['min', 'max', 'last'], name: 'value', type: 'string' }], help: 'TODO docs' } },
              withSeriesBy(value): { fieldConfig+: { defaults+: { color+: { seriesBy: value } } } },
            },
          '#withCustom': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'custom is specified by the PanelFieldConfig field\nin panel plugin schemas.' } },
          withCustom(value): { fieldConfig+: { defaults+: { custom: value } } },
          '#withCustomMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'custom is specified by the PanelFieldConfig field\nin panel plugin schemas.' } },
          withCustomMixin(value): { fieldConfig+: { defaults+: { custom+: value } } },
          '#withDecimals': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'Significant digits (for display)' } },
          withDecimals(value): { fieldConfig+: { defaults+: { decimals: value } } },
          '#withDescription': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Human readable field metadata' } },
          withDescription(value): { fieldConfig+: { defaults+: { description: value } } },
          '#withDisplayName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'The display value for this field.  This supports template variables blank is auto' } },
          withDisplayName(value): { fieldConfig+: { defaults+: { displayName: value } } },
          '#withDisplayNameFromDS': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'This can be used by data sources that return and explicit naming structure for values and labels\nWhen this property is configured, this value is used rather than the default naming strategy.' } },
          withDisplayNameFromDS(value): { fieldConfig+: { defaults+: { displayNameFromDS: value } } },
          '#withFilterable': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'True if data source field supports ad-hoc filters' } },
          withFilterable(value=true): { fieldConfig+: { defaults+: { filterable: value } } },
          '#withLinks': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'The behavior when clicking on a result' } },
          withLinks(value): { fieldConfig+: { defaults+: { links: (if std.isArray(value)
                                                                   then value
                                                                   else [value]) } } },
          '#withLinksMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'The behavior when clicking on a result' } },
          withLinksMixin(value): { fieldConfig+: { defaults+: { links+: (if std.isArray(value)
                                                                         then value
                                                                         else [value]) } } },
          '#withMappings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Convert input values into a display string' } },
          withMappings(value): { fieldConfig+: { defaults+: { mappings: (if std.isArray(value)
                                                                         then value
                                                                         else [value]) } } },
          '#withMappingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Convert input values into a display string' } },
          withMappingsMixin(value): { fieldConfig+: { defaults+: { mappings+: (if std.isArray(value)
                                                                               then value
                                                                               else [value]) } } },
          mappings+:
            {
              '#': { help: '', name: 'mappings' },
              ValueMap+:
                {
                  '#withOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
                  withOptions(value): { options: value },
                  '#withOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
                  withOptionsMixin(value): { options+: value },
                  '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withType(value): { type: value },
                },
              RangeMap+:
                {
                  '#withOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
                  withOptions(value): { options: value },
                  '#withOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
                  withOptionsMixin(value): { options+: value },
                  options+:
                    {
                      '#withFrom': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'to and from are `number | null` in current ts, really not sure what to do' } },
                      withFrom(value): { options+: { from: value } },
                      '#withResult': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
                      withResult(value): { options+: { result: value } },
                      '#withResultMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
                      withResultMixin(value): { options+: { result+: value } },
                      result+:
                        {
                          '#withColor': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                          withColor(value): { options+: { result+: { color: value } } },
                          '#withIcon': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                          withIcon(value): { options+: { result+: { icon: value } } },
                          '#withIndex': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
                          withIndex(value): { options+: { result+: { index: value } } },
                          '#withText': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                          withText(value): { options+: { result+: { text: value } } },
                        },
                      '#withTo': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
                      withTo(value): { options+: { to: value } },
                    },
                  '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withType(value): { type: value },
                },
              RegexMap+:
                {
                  '#withOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
                  withOptions(value): { options: value },
                  '#withOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
                  withOptionsMixin(value): { options+: value },
                  options+:
                    {
                      '#withPattern': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                      withPattern(value): { options+: { pattern: value } },
                      '#withResult': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
                      withResult(value): { options+: { result: value } },
                      '#withResultMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
                      withResultMixin(value): { options+: { result+: value } },
                      result+:
                        {
                          '#withColor': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                          withColor(value): { options+: { result+: { color: value } } },
                          '#withIcon': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                          withIcon(value): { options+: { result+: { icon: value } } },
                          '#withIndex': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
                          withIndex(value): { options+: { result+: { index: value } } },
                          '#withText': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                          withText(value): { options+: { result+: { text: value } } },
                        },
                    },
                  '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withType(value): { type: value },
                },
              SpecialValueMap+:
                {
                  '#withOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
                  withOptions(value): { options: value },
                  '#withOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
                  withOptionsMixin(value): { options+: value },
                  options+:
                    {
                      '#withMatch': { 'function': { args: [{ default: null, enums: ['true', 'false'], name: 'value', type: 'string' }], help: '' } },
                      withMatch(value): { options+: { match: value } },
                      '#withPattern': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                      withPattern(value): { options+: { pattern: value } },
                      '#withResult': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
                      withResult(value): { options+: { result: value } },
                      '#withResultMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
                      withResultMixin(value): { options+: { result+: value } },
                      result+:
                        {
                          '#withColor': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                          withColor(value): { options+: { result+: { color: value } } },
                          '#withIcon': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                          withIcon(value): { options+: { result+: { icon: value } } },
                          '#withIndex': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
                          withIndex(value): { options+: { result+: { index: value } } },
                          '#withText': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                          withText(value): { options+: { result+: { text: value } } },
                        },
                    },
                  '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
                  withType(value): { type: value },
                },
            },
          '#withMax': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
          withMax(value): { fieldConfig+: { defaults+: { max: value } } },
          '#withMin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
          withMin(value): { fieldConfig+: { defaults+: { min: value } } },
          '#withNoValue': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Alternative to empty string' } },
          withNoValue(value): { fieldConfig+: { defaults+: { noValue: value } } },
          '#withPath': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'An explicit path to the field in the datasource.  When the frame meta includes a path,\nThis will default to `${frame.meta.path}/${field.name}\n\nWhen defined, this value can be used as an identifier within the datasource scope, and\nmay be used to update the results' } },
          withPath(value): { fieldConfig+: { defaults+: { path: value } } },
          '#withThresholds': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withThresholds(value): { fieldConfig+: { defaults+: { thresholds: value } } },
          '#withThresholdsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withThresholdsMixin(value): { fieldConfig+: { defaults+: { thresholds+: value } } },
          thresholds+:
            {
              '#withMode': { 'function': { args: [{ default: null, enums: ['absolute', 'percentage'], name: 'value', type: 'string' }], help: '' } },
              withMode(value): { fieldConfig+: { defaults+: { thresholds+: { mode: value } } } },
              '#withSteps': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: "Must be sorted by 'value', first value is always -Infinity" } },
              withSteps(value): { fieldConfig+: { defaults+: { thresholds+: { steps: (if std.isArray(value)
                                                                                      then value
                                                                                      else [value]) } } } },
              '#withStepsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: "Must be sorted by 'value', first value is always -Infinity" } },
              withStepsMixin(value): { fieldConfig+: { defaults+: { thresholds+: { steps+: (if std.isArray(value)
                                                                                            then value
                                                                                            else [value]) } } } },
              steps+:
                {
                  '#': { help: '', name: 'steps' },
                  '#withColor': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO docs' } },
                  withColor(value): { color: value },
                  '#withIndex': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'Threshold index, an old property that is not needed an should only appear in older dashboards' } },
                  withIndex(value): { index: value },
                  '#withState': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO docs\nTODO are the values here enumerable into a disjunction?\nSome seem to be listed in typescript comment' } },
                  withState(value): { state: value },
                  '#withValue': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'TODO docs\nFIXME the corresponding typescript field is required/non-optional, but nulls currently appear here when serializing -Infinity to JSON' } },
                  withValue(value): { value: value },
                },
            },
          '#withUnit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Numeric Options' } },
          withUnit(value): { fieldConfig+: { defaults+: { unit: value } } },
          '#withWriteable': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'True if data source can write a value to the path.  Auth/authz are supported separately' } },
          withWriteable(value=true): { fieldConfig+: { defaults+: { writeable: value } } },
        },
      '#withOverrides': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withOverrides(value): { fieldConfig+: { overrides: (if std.isArray(value)
                                                          then value
                                                          else [value]) } },
      '#withOverridesMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withOverridesMixin(value): { fieldConfig+: { overrides+: (if std.isArray(value)
                                                                then value
                                                                else [value]) } },
      overrides+:
        {
          '#': { help: '', name: 'overrides' },
          '#withMatcher': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withMatcher(value): { matcher: value },
          '#withMatcherMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withMatcherMixin(value): { matcher+: value },
          matcher+:
            {
              '#withId': { 'function': { args: [{ default: '', enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value=''): { matcher+: { id: value } },
              '#withOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withOptions(value): { matcher+: { options: value } },
            },
          '#withProperties': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
          withProperties(value): { properties: (if std.isArray(value)
                                                then value
                                                else [value]) },
          '#withPropertiesMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
          withPropertiesMixin(value): { properties+: (if std.isArray(value)
                                                      then value
                                                      else [value]) },
          properties+:
            {
              '#': { help: '', name: 'properties' },
              '#withId': { 'function': { args: [{ default: '', enums: null, name: 'value', type: 'string' }], help: '' } },
              withId(value=''): { id: value },
              '#withValue': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withValue(value): { value: value },
            },
        },
    },
  '#withGridPos': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withGridPos(value): { gridPos: value },
  '#withGridPosMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withGridPosMixin(value): { gridPos+: value },
  gridPos+:
    {
      '#withH': { 'function': { args: [{ default: 9, enums: null, name: 'value', type: 'integer' }], help: 'Panel' } },
      withH(value=9): { gridPos+: { h: value } },
      '#withStatic': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'true if fixed' } },
      withStatic(value=true): { gridPos+: { static: value } },
      '#withW': { 'function': { args: [{ default: 12, enums: null, name: 'value', type: 'integer' }], help: 'Panel' } },
      withW(value=12): { gridPos+: { w: value } },
      '#withX': { 'function': { args: [{ default: 0, enums: null, name: 'value', type: 'integer' }], help: 'Panel x' } },
      withX(value=0): { gridPos+: { x: value } },
      '#withY': { 'function': { args: [{ default: 0, enums: null, name: 'value', type: 'integer' }], help: 'Panel y' } },
      withY(value=0): { gridPos+: { y: value } },
    },
  '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'TODO docs' } },
  withId(value): { id: value },
  '#withInterval': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO docs\nTODO tighter constraint' } },
  withInterval(value): { interval: value },
  '#withLibraryPanel': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withLibraryPanel(value): { libraryPanel: value },
  '#withLibraryPanelMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withLibraryPanelMixin(value): { libraryPanel+: value },
  libraryPanel+:
    {
      '#withName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withName(value): { libraryPanel+: { name: value } },
      '#withUid': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withUid(value): { libraryPanel+: { uid: value } },
    },
  '#withLinks': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Panel links.\nTODO fill this out - seems there are a couple variants?' } },
  withLinks(value): { links: (if std.isArray(value)
                              then value
                              else [value]) },
  '#withLinksMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Panel links.\nTODO fill this out - seems there are a couple variants?' } },
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
  '#withMaxDataPoints': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'TODO docs' } },
  withMaxDataPoints(value): { maxDataPoints: value },
  '#withOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'options is specified by the PanelOptions field in panel\nplugin schemas.' } },
  withOptions(value): { options: value },
  '#withOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'options is specified by the PanelOptions field in panel\nplugin schemas.' } },
  withOptionsMixin(value): { options+: value },
  '#withPluginVersion': { 'function': { args: [], help: '' } },
  withPluginVersion(): { pluginVersion: 'v10.0.0' },
  '#withRepeat': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Name of template variable to repeat for.' } },
  withRepeat(value): { repeat: value },
  '#withRepeatDirection': { 'function': { args: [{ default: 'h', enums: ['h', 'v'], name: 'value', type: 'string' }], help: "Direction to repeat in if 'repeat' is set.\n\"h\" for horizontal, \"v\" for vertical.\nTODO this is probably optional" } },
  withRepeatDirection(value='h'): { repeatDirection: value },
  '#withRepeatPanelId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'Id of the repeating panel.' } },
  withRepeatPanelId(value): { repeatPanelId: value },
  '#withTags': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'TODO docs' } },
  withTags(value): { tags: (if std.isArray(value)
                            then value
                            else [value]) },
  '#withTagsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'TODO docs' } },
  withTagsMixin(value): { tags+: (if std.isArray(value)
                                  then value
                                  else [value]) },
  '#withTargets': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'TODO docs' } },
  withTargets(value): { targets: (if std.isArray(value)
                                  then value
                                  else [value]) },
  '#withTargetsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'TODO docs' } },
  withTargetsMixin(value): { targets+: (if std.isArray(value)
                                        then value
                                        else [value]) },
  '#withThresholds': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'TODO docs - seems to be an old field from old dashboard alerts?' } },
  withThresholds(value): { thresholds: (if std.isArray(value)
                                        then value
                                        else [value]) },
  '#withThresholdsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'TODO docs - seems to be an old field from old dashboard alerts?' } },
  withThresholdsMixin(value): { thresholds+: (if std.isArray(value)
                                              then value
                                              else [value]) },
  '#withTimeFrom': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO docs\nTODO tighter constraint' } },
  withTimeFrom(value): { timeFrom: value },
  '#withTimeRegions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'TODO docs' } },
  withTimeRegions(value): { timeRegions: (if std.isArray(value)
                                          then value
                                          else [value]) },
  '#withTimeRegionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'TODO docs' } },
  withTimeRegionsMixin(value): { timeRegions+: (if std.isArray(value)
                                                then value
                                                else [value]) },
  '#withTimeShift': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'TODO docs\nTODO tighter constraint' } },
  withTimeShift(value): { timeShift: value },
  '#withTitle': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Panel title.' } },
  withTitle(value): { title: value },
  '#withTransformations': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withTransformations(value): { transformations: (if std.isArray(value)
                                                  then value
                                                  else [value]) },
  '#withTransformationsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withTransformationsMixin(value): { transformations+: (if std.isArray(value)
                                                        then value
                                                        else [value]) },
  transformations+:
    {
      '#': { help: '', name: 'transformations' },
      '#withDisabled': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Disabled transformations are skipped' } },
      withDisabled(value=true): { disabled: value },
      '#withFilter': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withFilter(value): { filter: value },
      '#withFilterMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withFilterMixin(value): { filter+: value },
      filter+:
        {
          '#withId': { 'function': { args: [{ default: '', enums: null, name: 'value', type: 'string' }], help: '' } },
          withId(value=''): { filter+: { id: value } },
          '#withOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withOptions(value): { filter+: { options: value } },
        },
      '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Unique identifier of transformer' } },
      withId(value): { id: value },
      '#withOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Options to be passed to the transformer\nValid options depend on the transformer id' } },
      withOptions(value): { options: value },
    },
  '#withTransparent': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Whether to display the panel without a background.' } },
  withTransparent(value=true): { transparent: value },
  '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'The panel plugin type id. May not be empty.' } },
  withType(value): { type: value },
}
