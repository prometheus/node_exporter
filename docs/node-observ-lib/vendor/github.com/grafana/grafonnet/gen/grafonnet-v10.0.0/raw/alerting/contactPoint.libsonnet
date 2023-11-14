// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.alerting.contactPoint', name: 'contactPoint' },
  '#withDisableResolveMessage': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
  withDisableResolveMessage(value=true): { disableResolveMessage: value },
  '#withName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Name is used as grouping key in the UI. Contact points with the\nsame name will be grouped in the UI.' } },
  withName(value): { name: value },
  '#withProvenance': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
  withProvenance(value): { provenance: value },
  '#withSettings': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withSettings(value): { settings: value },
  '#withSettingsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withSettingsMixin(value): { settings+: value },
  '#withType': { 'function': { args: [{ default: null, enums: ['alertmanager', ' dingding', ' discord', ' email', ' googlechat', ' kafka', ' line', ' opsgenie', ' pagerduty', ' pushover', ' sensugo', ' slack', ' teams', ' telegram', ' threema', ' victorops', ' webhook', ' wecom'], name: 'value', type: 'string' }], help: '' } },
  withType(value): { type: value },
  '#withUid': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'UID is the unique identifier of the contact point. The UID can be\nset by the user.' } },
  withUid(value): { uid: value },
}
