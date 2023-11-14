// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.preferences', name: 'preferences' },
  '#withHomeDashboardUID': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'UID for the home dashboard' } },
  withHomeDashboardUID(value): { homeDashboardUID: value },
  '#withLanguage': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Selected language (beta)' } },
  withLanguage(value): { language: value },
  '#withQueryHistory': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withQueryHistory(value): { queryHistory: value },
  '#withQueryHistoryMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withQueryHistoryMixin(value): { queryHistory+: value },
  queryHistory+:
    {
      '#withHomeTab': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: "one of: '' | 'query' | 'starred';" } },
      withHomeTab(value): { queryHistory+: { homeTab: value } },
    },
  '#withTheme': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'light, dark, empty is default' } },
  withTheme(value): { theme: value },
  '#withTimezone': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'The timezone selection\nTODO: this should use the timezone defined in common' } },
  withTimezone(value): { timezone: value },
  '#withWeekStart': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'day of the week (sunday, monday, etc)' } },
  withWeekStart(value): { weekStart: value },
}
