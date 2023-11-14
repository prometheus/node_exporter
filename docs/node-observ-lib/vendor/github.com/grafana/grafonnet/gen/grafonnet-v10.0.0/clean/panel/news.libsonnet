// This file is generated, do not manually edit.
(import '../../clean/panel.libsonnet')
+ {
  '#': { help: 'grafonnet.panel.news', name: 'news' },
  panelOptions+:
    {
      '#withType': { 'function': { args: [], help: '' } },
      withType(): { type: 'news' },
    },
  options+:
    {
      '#withFeedUrl': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'empty/missing will default to grafana blog' } },
      withFeedUrl(value): { options+: { feedUrl: value } },
      '#withShowImage': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowImage(value=true): { options+: { showImage: value } },
    },
}
+ { panelOptions+: { '#withType':: {} } }
