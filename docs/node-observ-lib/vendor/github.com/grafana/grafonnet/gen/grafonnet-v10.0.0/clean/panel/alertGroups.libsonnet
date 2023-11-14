// This file is generated, do not manually edit.
(import '../../clean/panel.libsonnet')
+ {
  '#': { help: 'grafonnet.panel.alertGroups', name: 'alertGroups' },
  panelOptions+:
    {
      '#withType': { 'function': { args: [], help: '' } },
      withType(): { type: 'alertGroups' },
    },
  options+:
    {
      '#withAlertmanager': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Name of the alertmanager used as a source for alerts' } },
      withAlertmanager(value): { options+: { alertmanager: value } },
      '#withExpandAll': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Expand all alert groups by default' } },
      withExpandAll(value=true): { options+: { expandAll: value } },
      '#withLabels': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Comma-separated list of values used to filter alert results' } },
      withLabels(value): { options+: { labels: value } },
    },
}
+ { panelOptions+: { '#withType':: {} } }
