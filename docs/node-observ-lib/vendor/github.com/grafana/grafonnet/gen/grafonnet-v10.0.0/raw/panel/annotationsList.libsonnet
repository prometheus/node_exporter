// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.panel.annotationsList', name: 'annotationsList' },
  '#withOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withOptions(value): { options: value },
  '#withOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withOptionsMixin(value): { options+: value },
  options+:
    {
      '#withLimit': { 'function': { args: [{ default: 10, enums: null, name: 'value', type: 'integer' }], help: '' } },
      withLimit(value=10): { options+: { limit: value } },
      '#withNavigateAfter': { 'function': { args: [{ default: '10m', enums: null, name: 'value', type: 'string' }], help: '' } },
      withNavigateAfter(value='10m'): { options+: { navigateAfter: value } },
      '#withNavigateBefore': { 'function': { args: [{ default: '10m', enums: null, name: 'value', type: 'string' }], help: '' } },
      withNavigateBefore(value='10m'): { options+: { navigateBefore: value } },
      '#withNavigateToPanel': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withNavigateToPanel(value=true): { options+: { navigateToPanel: value } },
      '#withOnlyFromThisDashboard': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withOnlyFromThisDashboard(value=true): { options+: { onlyFromThisDashboard: value } },
      '#withOnlyInTimeRange': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withOnlyInTimeRange(value=true): { options+: { onlyInTimeRange: value } },
      '#withShowTags': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowTags(value=true): { options+: { showTags: value } },
      '#withShowTime': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowTime(value=true): { options+: { showTime: value } },
      '#withShowUser': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowUser(value=true): { options+: { showUser: value } },
      '#withTags': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withTags(value): { options+: { tags: (if std.isArray(value)
                                            then value
                                            else [value]) } },
      '#withTagsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withTagsMixin(value): { options+: { tags+: (if std.isArray(value)
                                                  then value
                                                  else [value]) } },
    },
  '#withType': { 'function': { args: [], help: '' } },
  withType(): { type: 'annolist' },
}
