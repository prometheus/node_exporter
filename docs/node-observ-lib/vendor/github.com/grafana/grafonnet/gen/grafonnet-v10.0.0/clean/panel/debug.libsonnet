// This file is generated, do not manually edit.
(import '../../clean/panel.libsonnet')
+ {
  '#': { help: 'grafonnet.panel.debug', name: 'debug' },
  panelOptions+:
    {
      '#withType': { 'function': { args: [], help: '' } },
      withType(): { type: 'debug' },
    },
  options+:
    {
      '#withCounters': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withCounters(value): { options+: { counters: value } },
      '#withCountersMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withCountersMixin(value): { options+: { counters+: value } },
      counters+:
        {
          '#withDataChanged': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withDataChanged(value=true): { options+: { counters+: { dataChanged: value } } },
          '#withRender': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withRender(value=true): { options+: { counters+: { render: value } } },
          '#withSchemaChanged': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withSchemaChanged(value=true): { options+: { counters+: { schemaChanged: value } } },
        },
      '#withMode': { 'function': { args: [{ default: null, enums: ['render', 'events', 'cursor', 'State', 'ThrowError'], name: 'value', type: 'string' }], help: '' } },
      withMode(value): { options+: { mode: value } },
    },
}
+ { panelOptions+: { '#withType':: {} } }
