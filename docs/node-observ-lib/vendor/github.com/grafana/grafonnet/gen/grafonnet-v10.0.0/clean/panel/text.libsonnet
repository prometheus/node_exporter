// This file is generated, do not manually edit.
(import '../../clean/panel.libsonnet')
+ {
  '#': { help: 'grafonnet.panel.text', name: 'text' },
  panelOptions+:
    {
      '#withType': { 'function': { args: [], help: '' } },
      withType(): { type: 'text' },
    },
  options+:
    {
      '#withCode': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withCode(value): { options+: { code: value } },
      '#withCodeMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withCodeMixin(value): { options+: { code+: value } },
      code+:
        {
          '#withLanguage': { 'function': { args: [{ default: 'plaintext', enums: ['plaintext', 'yaml', 'xml', 'typescript', 'sql', 'go', 'markdown', 'html', 'json'], name: 'value', type: 'string' }], help: '' } },
          withLanguage(value='plaintext'): { options+: { code+: { language: value } } },
          '#withShowLineNumbers': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withShowLineNumbers(value=true): { options+: { code+: { showLineNumbers: value } } },
          '#withShowMiniMap': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withShowMiniMap(value=true): { options+: { code+: { showMiniMap: value } } },
        },
      '#withContent': { 'function': { args: [{ default: '# Title\n\nFor markdown syntax help: [commonmark.org/help](https://commonmark.org/help/)', enums: null, name: 'value', type: 'string' }], help: '' } },
      withContent(value='# Title\n\nFor markdown syntax help: [commonmark.org/help](https://commonmark.org/help/)'): { options+: { content: value } },
      '#withMode': { 'function': { args: [{ default: null, enums: ['html', 'markdown', 'code'], name: 'value', type: 'string' }], help: '' } },
      withMode(value): { options+: { mode: value } },
    },
}
+ { panelOptions+: { '#withType':: {} } }
