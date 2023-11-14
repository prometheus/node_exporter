// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.panel.news', name: 'news' },
  '#withOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withOptions(value): { options: value },
  '#withOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withOptionsMixin(value): { options+: value },
  options+:
    {
      '#withFeedUrl': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'empty/missing will default to grafana blog' } },
      withFeedUrl(value): { options+: { feedUrl: value } },
      '#withShowImage': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowImage(value=true): { options+: { showImage: value } },
    },
  '#withType': { 'function': { args: [], help: '' } },
  withType(): { type: 'news' },
}
