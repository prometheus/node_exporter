// This file is generated, do not manually edit.
(import '../../clean/panel.libsonnet')
+ {
  '#': { help: 'grafonnet.panel.barChart', name: 'barChart' },
  panelOptions+:
    {
      '#withType': { 'function': { args: [], help: '' } },
      withType(): { type: 'barchart' },
    },
  fieldConfig+:
    {
      defaults+:
        {
          custom+:
            {
              '#withAxisCenteredZero': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withAxisCenteredZero(value=true): { fieldConfig+: { defaults+: { custom+: { axisCenteredZero: value } } } },
              '#withAxisColorMode': { 'function': { args: [{ default: null, enums: ['text', 'series'], name: 'value', type: 'string' }], help: 'TODO docs' } },
              withAxisColorMode(value): { fieldConfig+: { defaults+: { custom+: { axisColorMode: value } } } },
              '#withAxisGridShow': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withAxisGridShow(value=true): { fieldConfig+: { defaults+: { custom+: { axisGridShow: value } } } },
              '#withAxisLabel': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withAxisLabel(value): { fieldConfig+: { defaults+: { custom+: { axisLabel: value } } } },
              '#withAxisPlacement': { 'function': { args: [{ default: null, enums: ['auto', 'top', 'right', 'bottom', 'left', 'hidden'], name: 'value', type: 'string' }], help: 'TODO docs' } },
              withAxisPlacement(value): { fieldConfig+: { defaults+: { custom+: { axisPlacement: value } } } },
              '#withAxisSoftMax': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
              withAxisSoftMax(value): { fieldConfig+: { defaults+: { custom+: { axisSoftMax: value } } } },
              '#withAxisSoftMin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
              withAxisSoftMin(value): { fieldConfig+: { defaults+: { custom+: { axisSoftMin: value } } } },
              '#withAxisWidth': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
              withAxisWidth(value): { fieldConfig+: { defaults+: { custom+: { axisWidth: value } } } },
              '#withScaleDistribution': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
              withScaleDistribution(value): { fieldConfig+: { defaults+: { custom+: { scaleDistribution: value } } } },
              '#withScaleDistributionMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
              withScaleDistributionMixin(value): { fieldConfig+: { defaults+: { custom+: { scaleDistribution+: value } } } },
              scaleDistribution+:
                {
                  '#withLinearThreshold': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
                  withLinearThreshold(value): { fieldConfig+: { defaults+: { custom+: { scaleDistribution+: { linearThreshold: value } } } } },
                  '#withLog': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
                  withLog(value): { fieldConfig+: { defaults+: { custom+: { scaleDistribution+: { log: value } } } } },
                  '#withType': { 'function': { args: [{ default: null, enums: ['linear', 'log', 'ordinal', 'symlog'], name: 'value', type: 'string' }], help: 'TODO docs' } },
                  withType(value): { fieldConfig+: { defaults+: { custom+: { scaleDistribution+: { type: value } } } } },
                },
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
              '#withFillOpacity': { 'function': { args: [{ default: 80, enums: null, name: 'value', type: 'integer' }], help: 'Controls the fill opacity of the bars.' } },
              withFillOpacity(value=80): { fieldConfig+: { defaults+: { custom+: { fillOpacity: value } } } },
              '#withGradientMode': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Set the mode of the gradient fill. Fill gradient is based on the line color. To change the color, use the standard color scheme field option.\nGradient appearance is influenced by the Fill opacity setting.' } },
              withGradientMode(value): { fieldConfig+: { defaults+: { custom+: { gradientMode: value } } } },
              '#withLineWidth': { 'function': { args: [{ default: 1, enums: null, name: 'value', type: 'integer' }], help: 'Controls line width of the bars.' } },
              withLineWidth(value=1): { fieldConfig+: { defaults+: { custom+: { lineWidth: value } } } },
              '#withThresholdsStyle': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
              withThresholdsStyle(value): { fieldConfig+: { defaults+: { custom+: { thresholdsStyle: value } } } },
              '#withThresholdsStyleMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
              withThresholdsStyleMixin(value): { fieldConfig+: { defaults+: { custom+: { thresholdsStyle+: value } } } },
              thresholdsStyle+:
                {
                  '#withMode': { 'function': { args: [{ default: null, enums: ['off', 'line', 'dashed', 'area', 'line+area', 'dashed+area', 'series'], name: 'value', type: 'string' }], help: 'TODO docs' } },
                  withMode(value): { fieldConfig+: { defaults+: { custom+: { thresholdsStyle+: { mode: value } } } } },
                },
            },
        },
    },
  options+:
    {
      '#withLegend': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
      withLegend(value): { options+: { legend: value } },
      '#withLegendMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
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
        },
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
      '#withBarRadius': { 'function': { args: [{ default: 0, enums: null, name: 'value', type: 'number' }], help: 'Controls the radius of each bar.' } },
      withBarRadius(value=0): { options+: { barRadius: value } },
      '#withBarWidth': { 'function': { args: [{ default: 0.96999999999999997, enums: null, name: 'value', type: 'number' }], help: 'Controls the width of bars. 1 = Max width, 0 = Min width.' } },
      withBarWidth(value=0.96999999999999997): { options+: { barWidth: value } },
      '#withColorByField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Use the color value for a sibling field to color each bar value.' } },
      withColorByField(value): { options+: { colorByField: value } },
      '#withFullHighlight': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Enables mode which highlights the entire bar area and shows tooltip when cursor\nhovers over highlighted area' } },
      withFullHighlight(value=true): { options+: { fullHighlight: value } },
      '#withGroupWidth': { 'function': { args: [{ default: 0.69999999999999996, enums: null, name: 'value', type: 'number' }], help: 'Controls the width of groups. 1 = max with, 0 = min width.' } },
      withGroupWidth(value=0.69999999999999996): { options+: { groupWidth: value } },
      '#withOrientation': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Controls the orientation of the bar chart, either vertical or horizontal.' } },
      withOrientation(value): { options+: { orientation: value } },
      '#withShowValue': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'This controls whether values are shown on top or to the left of bars.' } },
      withShowValue(value): { options+: { showValue: value } },
      '#withStacking': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Controls whether bars are stacked or not, either normally or in percent mode.' } },
      withStacking(value): { options+: { stacking: value } },
      '#withXField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Manually select which field from the dataset to represent the x field.' } },
      withXField(value): { options+: { xField: value } },
      '#withXTickLabelMaxLength': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'Sets the max length that a label can have before it is truncated.' } },
      withXTickLabelMaxLength(value): { options+: { xTickLabelMaxLength: value } },
      '#withXTickLabelRotation': { 'function': { args: [{ default: 0, enums: null, name: 'value', type: 'integer' }], help: 'Controls the rotation of the x axis labels.' } },
      withXTickLabelRotation(value=0): { options+: { xTickLabelRotation: value } },
      '#withXTickLabelSpacing': { 'function': { args: [{ default: 0, enums: null, name: 'value', type: 'integer' }], help: 'Controls the spacing between x axis labels.\nnegative values indicate backwards skipping behavior' } },
      withXTickLabelSpacing(value=0): { options+: { xTickLabelSpacing: value } },
    },
}
+ { panelOptions+: { '#withType':: {} } }
