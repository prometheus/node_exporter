// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.panel.debug', name: 'debug' },
  '#withDebugMode': { 'function': { args: [{ default: null, enums: ['render', 'events', 'cursor', 'State', 'ThrowError'], name: 'value', type: 'string' }], help: '' } },
  withDebugMode(value): { DebugMode: value },
  '#withUpdateConfig': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withUpdateConfig(value): { UpdateConfig: value },
  '#withUpdateConfigMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withUpdateConfigMixin(value): { UpdateConfig+: value },
  UpdateConfig+:
    {
      '#withDataChanged': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withDataChanged(value=true): { UpdateConfig+: { dataChanged: value } },
      '#withRender': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withRender(value=true): { UpdateConfig+: { render: value } },
      '#withSchemaChanged': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withSchemaChanged(value=true): { UpdateConfig+: { schemaChanged: value } },
    },
  '#withOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withOptions(value): { options: value },
  '#withOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withOptionsMixin(value): { options+: value },
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
  '#withType': { 'function': { args: [], help: '' } },
  withType(): { type: 'debug' },
}
