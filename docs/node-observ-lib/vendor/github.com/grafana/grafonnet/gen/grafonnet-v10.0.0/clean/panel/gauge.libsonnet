// This file is generated, do not manually edit.
(import '../../clean/panel.libsonnet')
+ {
  '#': { help: 'grafonnet.panel.gauge', name: 'gauge' },
  panelOptions+:
    {
      '#withType': { 'function': { args: [], help: '' } },
      withType(): { type: 'gauge' },
    },
  options+:
    {
      '#withText': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
      withText(value): { options+: { text: value } },
      '#withTextMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
      withTextMixin(value): { options+: { text+: value } },
      text+:
        {
          '#withTitleSize': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'Explicit title text size' } },
          withTitleSize(value): { options+: { text+: { titleSize: value } } },
          '#withValueSize': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'Explicit value text size' } },
          withValueSize(value): { options+: { text+: { valueSize: value } } },
        },
      '#withOrientation': { 'function': { args: [{ default: null, enums: ['auto', 'vertical', 'horizontal'], name: 'value', type: 'string' }], help: 'TODO docs' } },
      withOrientation(value): { options+: { orientation: value } },
      '#withReduceOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
      withReduceOptions(value): { options+: { reduceOptions: value } },
      '#withReduceOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
      withReduceOptionsMixin(value): { options+: { reduceOptions+: value } },
      reduceOptions+:
        {
          '#withCalcs': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'When !values, pick one value for the whole field' } },
          withCalcs(value): { options+: { reduceOptions+: { calcs: (if std.isArray(value)
                                                                    then value
                                                                    else [value]) } } },
          '#withCalcsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'When !values, pick one value for the whole field' } },
          withCalcsMixin(value): { options+: { reduceOptions+: { calcs+: (if std.isArray(value)
                                                                          then value
                                                                          else [value]) } } },
          '#withFields': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Which fields to show.  By default this is only numeric fields' } },
          withFields(value): { options+: { reduceOptions+: { fields: value } } },
          '#withLimit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'if showing all values limit' } },
          withLimit(value): { options+: { reduceOptions+: { limit: value } } },
          '#withValues': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'If true show each row value' } },
          withValues(value=true): { options+: { reduceOptions+: { values: value } } },
        },
      '#withShowThresholdLabels': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowThresholdLabels(value=true): { options+: { showThresholdLabels: value } },
      '#withShowThresholdMarkers': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
      withShowThresholdMarkers(value=true): { options+: { showThresholdMarkers: value } },
    },
}
+ { panelOptions+: { '#withType':: {} } }
