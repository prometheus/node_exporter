// This file is generated, do not manually edit.
(import '../../clean/panel.libsonnet')
+ {
  '#': { help: 'grafonnet.panel.pieChart', name: 'pieChart' },
  panelOptions+:
    {
      '#withType': { 'function': { args: [], help: '' } },
      withType(): { type: 'piechart' },
    },
  fieldConfig+:
    {
      defaults+:
        {
          custom+:
            {
              '#withHideFrom': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
              withHideFrom(value): { fieldConfig+: { defaults+: { custom+: { hideFrom: value } } } },
              '#withHideFromMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
              withHideFromMixin(value): { fieldConfig+: { defaults+: { custom+: { hideFrom+: value } } } },
              hideFrom+:
                {
                  '#withLegend': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
                  withLegend(value=true): { fieldConfig+: { defaults+: { custom+: { hideFrom+: { legend: value } } } } },
                  '#withTooltip': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
                  withTooltip(value=true): { fieldConfig+: { defaults+: { custom+: { hideFrom+: { tooltip: value } } } } },
                  '#withViz': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
                  withViz(value=true): { fieldConfig+: { defaults+: { custom+: { hideFrom+: { viz: value } } } } },
                },
            },
        },
    },
  options+:
    {
      '#withTooltip': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
      withTooltip(value): { options+: { tooltip: value } },
      '#withTooltipMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
      withTooltipMixin(value): { options+: { tooltip+: value } },
      tooltip+:
        {
          '#withMode': { 'function': { args: [{ default: null, enums: ['single', 'multi', 'none'], name: 'value', type: 'string' }], help: 'TODO docs' } },
          withMode(value): { options+: { tooltip+: { mode: value } } },
          '#withSort': { 'function': { args: [{ default: null, enums: ['asc', 'desc', 'none'], name: 'value', type: 'string' }], help: 'TODO docs' } },
          withSort(value): { options+: { tooltip+: { sort: value } } },
        },
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
      '#withDisplayLabels': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withDisplayLabels(value): { options+: { displayLabels: (if std.isArray(value)
                                                              then value
                                                              else [value]) } },
      '#withDisplayLabelsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withDisplayLabelsMixin(value): { options+: { displayLabels+: (if std.isArray(value)
                                                                    then value
                                                                    else [value]) } },
      '#withLegend': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withLegend(value): { options+: { legend: value } },
      '#withLegendMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withLegendMixin(value): { options+: { legend+: value } },
      legend+:
        {
          '#withAsTable': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withAsTable(value=true): { options+: { legend+: { asTable: value } } },
          '#withCalcs': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
          withCalcs(value): { options+: { legend+: { calcs: (if std.isArray(value)
                                                             then value
                                                             else [value]) } } },
          '#withCalcsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
          withCalcsMixin(value): { options+: { legend+: { calcs+: (if std.isArray(value)
                                                                   then value
                                                                   else [value]) } } },
          '#withDisplayMode': { 'function': { args: [{ default: null, enums: ['list', 'table', 'hidden'], name: 'value', type: 'string' }], help: 'TODO docs\nNote: "hidden" needs to remain as an option for plugins compatibility' } },
          withDisplayMode(value): { options+: { legend+: { displayMode: value } } },
          '#withIsVisible': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withIsVisible(value=true): { options+: { legend+: { isVisible: value } } },
          '#withPlacement': { 'function': { args: [{ default: null, enums: ['bottom', 'right'], name: 'value', type: 'string' }], help: 'TODO docs' } },
          withPlacement(value): { options+: { legend+: { placement: value } } },
          '#withShowLegend': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withShowLegend(value=true): { options+: { legend+: { showLegend: value } } },
          '#withSortBy': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withSortBy(value): { options+: { legend+: { sortBy: value } } },
          '#withSortDesc': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withSortDesc(value=true): { options+: { legend+: { sortDesc: value } } },
          '#withWidth': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
          withWidth(value): { options+: { legend+: { width: value } } },
          '#withValues': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
          withValues(value): { options+: { legend+: { values: (if std.isArray(value)
                                                               then value
                                                               else [value]) } } },
          '#withValuesMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
          withValuesMixin(value): { options+: { legend+: { values+: (if std.isArray(value)
                                                                     then value
                                                                     else [value]) } } },
        },
      '#withPieType': { 'function': { args: [{ default: null, enums: ['pie', 'donut'], name: 'value', type: 'string' }], help: 'Select the pie chart display style.' } },
      withPieType(value): { options+: { pieType: value } },
    },
}
+ { panelOptions+: { '#withType':: {} } }
