// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.panel.datagrid', name: 'datagrid' },
  '#withOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withOptions(value): { options: value },
  '#withOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withOptionsMixin(value): { options+: value },
  options+:
    {
      '#withSelectedSeries': { 'function': { args: [{ default: 0, enums: null, name: 'value', type: 'integer' }], help: '' } },
      withSelectedSeries(value=0): { options+: { selectedSeries: value } },
    },
  '#withType': { 'function': { args: [], help: '' } },
  withType(): { type: 'datagrid' },
}
