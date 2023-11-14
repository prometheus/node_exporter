// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.panel.text', name: 'text' },
  '#withCodeLanguage': { 'function': { args: [{ default: 'plaintext', enums: ['plaintext', 'yaml', 'xml', 'typescript', 'sql', 'go', 'markdown', 'html', 'json'], name: 'value', type: 'string' }], help: '' } },
  withCodeLanguage(value='plaintext'): { CodeLanguage: value },
  '#withCodeOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withCodeOptions(value): { CodeOptions: value },
  '#withCodeOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withCodeOptionsMixin(value): { CodeOptions+: value },
  CodeOptions+:
    {
      '#withLanguage': { 'function': { args: [{ default: 'plaintext', enums: ['plaintext', 'yaml', 'xml', 'typescript', 'sql', 'go', 'markdown', 'html', 'json'], name: 'value', type: 'string' }], help: '' } },
      withLanguage(value='plaintext'): { CodeOptions+: { language: value } },
      '#withShowLineNumbers': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowLineNumbers(value=true): { CodeOptions+: { showLineNumbers: value } },
      '#withShowMiniMap': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowMiniMap(value=true): { CodeOptions+: { showMiniMap: value } },
    },
  '#withTextMode': { 'function': { args: [{ default: null, enums: ['html', 'markdown', 'code'], name: 'value', type: 'string' }], help: '' } },
  withTextMode(value): { TextMode: value },
  '#withOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withOptions(value): { options: value },
  '#withOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withOptionsMixin(value): { options+: value },
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
  '#withType': { 'function': { args: [], help: '' } },
  withType(): { type: 'text' },
}
