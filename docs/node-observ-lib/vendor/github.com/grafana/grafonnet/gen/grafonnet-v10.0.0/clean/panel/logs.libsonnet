// This file is generated, do not manually edit.
(import '../../clean/panel.libsonnet')
+ {
  '#': { help: 'grafonnet.panel.logs', name: 'logs' },
  panelOptions+:
    {
      '#withType': { 'function': { args: [], help: '' } },
      withType(): { type: 'logs' },
    },
  options+:
    {
      '#withDedupStrategy': { 'function': { args: [{ default: null, enums: ['none', 'exact', 'numbers', 'signature'], name: 'value', type: 'string' }], help: '' } },
      withDedupStrategy(value): { options+: { dedupStrategy: value } },
      '#withEnableLogDetails': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withEnableLogDetails(value=true): { options+: { enableLogDetails: value } },
      '#withPrettifyLogMessage': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withPrettifyLogMessage(value=true): { options+: { prettifyLogMessage: value } },
      '#withShowCommonLabels': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowCommonLabels(value=true): { options+: { showCommonLabels: value } },
      '#withShowLabels': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowLabels(value=true): { options+: { showLabels: value } },
      '#withShowTime': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowTime(value=true): { options+: { showTime: value } },
      '#withSortOrder': { 'function': { args: [{ default: null, enums: ['Descending', 'Ascending'], name: 'value', type: 'string' }], help: '' } },
      withSortOrder(value): { options+: { sortOrder: value } },
      '#withWrapLogMessage': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withWrapLogMessage(value=true): { options+: { wrapLogMessage: value } },
    },
}
+ { panelOptions+: { '#withType':: {} } }
