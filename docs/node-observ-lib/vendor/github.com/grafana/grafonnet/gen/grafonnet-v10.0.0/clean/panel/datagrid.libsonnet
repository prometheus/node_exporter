// This file is generated, do not manually edit.
(import '../../clean/panel.libsonnet')
+ {
  '#': { help: 'grafonnet.panel.datagrid', name: 'datagrid' },
  panelOptions+:
    {
      '#withType': { 'function': { args: [], help: '' } },
      withType(): { type: 'datagrid' },
    },
  options+:
    {
      '#withSelectedSeries': { 'function': { args: [{ default: 0, enums: null, name: 'value', type: 'integer' }], help: '' } },
      withSelectedSeries(value=0): { options+: { selectedSeries: value } },
    },
}
+ { panelOptions+: { '#withType':: {} } }
