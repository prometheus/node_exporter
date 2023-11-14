// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.panel.row', name: 'row' },
  '#withCollapsed': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
  withCollapsed(value=true): { collapsed: value },
  '#withDatasource': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Name of default datasource.' } },
  withDatasource(value): { datasource: value },
  '#withDatasourceMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Name of default datasource.' } },
  withDatasourceMixin(value): { datasource+: value },
  datasource+:
    {
      '#withType': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withType(value): { datasource+: { type: value } },
      '#withUid': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
      withUid(value): { datasource+: { uid: value } },
    },
  '#withGridPos': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withGridPos(value): { gridPos: value },
  '#withGridPosMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
  withGridPosMixin(value): { gridPos+: value },
  gridPos+:
    {
      '#withH': { 'function': { args: [{ default: 9, enums: null, name: 'value', type: 'integer' }], help: 'Panel' } },
      withH(value=9): { gridPos+: { h: value } },
      '#withStatic': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'true if fixed' } },
      withStatic(value=true): { gridPos+: { static: value } },
      '#withW': { 'function': { args: [{ default: 12, enums: null, name: 'value', type: 'integer' }], help: 'Panel' } },
      withW(value=12): { gridPos+: { w: value } },
      '#withX': { 'function': { args: [{ default: 0, enums: null, name: 'value', type: 'integer' }], help: 'Panel x' } },
      withX(value=0): { gridPos+: { x: value } },
      '#withY': { 'function': { args: [{ default: 0, enums: null, name: 'value', type: 'integer' }], help: 'Panel y' } },
      withY(value=0): { gridPos+: { y: value } },
    },
  '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
  withId(value): { id: value },
  '#withPanels': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withPanels(value): { panels: (if std.isArray(value)
                                then value
                                else [value]) },
  '#withPanelsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
  withPanelsMixin(value): { panels+: (if std.isArray(value)
                                      then value
                                      else [value]) },
  '#withRepeat': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Name of template variable to repeat for.' } },
  withRepeat(value): { repeat: value },
  '#withTitle': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
  withTitle(value): { title: value },
  '#withType': { 'function': { args: [], help: '' } },
  withType(): { type: 'row' },
}
+ (import '../../custom/row.libsonnet')
