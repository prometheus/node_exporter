// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.panel.dashboardList', name: 'dashboardList' },
  '#withOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withOptions(value): { options: value },
  '#withOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withOptionsMixin(value): { options+: value },
  options+:
    {
      '#withFolderId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
      withFolderId(value): { options+: { folderId: value } },
      '#withIncludeVars': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withIncludeVars(value=true): { options+: { includeVars: value } },
      '#withKeepTime': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withKeepTime(value=true): { options+: { keepTime: value } },
      '#withMaxItems': { 'function': { args: [{ default: 10, enums: null, name: 'value', type: 'integer' }], help: '' } },
      withMaxItems(value=10): { options+: { maxItems: value } },
      '#withQuery': { 'function': { args: [{ default: '', enums: null, name: 'value', type: 'string' }], help: '' } },
      withQuery(value=''): { options+: { query: value } },
      '#withShowHeadings': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowHeadings(value=true): { options+: { showHeadings: value } },
      '#withShowRecentlyViewed': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowRecentlyViewed(value=true): { options+: { showRecentlyViewed: value } },
      '#withShowSearch': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowSearch(value=true): { options+: { showSearch: value } },
      '#withShowStarred': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowStarred(value=true): { options+: { showStarred: value } },
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
  withType(): { type: 'dashlist' },
}
