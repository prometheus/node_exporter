// This file is generated, do not manually edit.
(import '../../clean/panel.libsonnet')
+ {
  '#': { help: 'grafonnet.panel.xyChart', name: 'xyChart' },
  panelOptions+:
    {
      '#withType': { 'function': { args: [], help: '' } },
      withType(): { type: 'xychart' },
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
      '#withDims': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withDims(value): { options+: { dims: value } },
      '#withDimsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withDimsMixin(value): { options+: { dims+: value } },
      dims+:
        {
          '#withExclude': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
          withExclude(value): { options+: { dims+: { exclude: (if std.isArray(value)
                                                               then value
                                                               else [value]) } } },
          '#withExcludeMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
          withExcludeMixin(value): { options+: { dims+: { exclude+: (if std.isArray(value)
                                                                     then value
                                                                     else [value]) } } },
          '#withFrame': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
          withFrame(value): { options+: { dims+: { frame: value } } },
          '#withX': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withX(value): { options+: { dims+: { x: value } } },
        },
      '#withSeries': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withSeries(value): { options+: { series: (if std.isArray(value)
                                                then value
                                                else [value]) } },
      '#withSeriesMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
      withSeriesMixin(value): { options+: { series+: (if std.isArray(value)
                                                      then value
                                                      else [value]) } },
      series+:
        {
          '#': { help: '', name: 'series' },
          '#withHideFrom': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
          withHideFrom(value): { hideFrom: value },
          '#withHideFromMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
          withHideFromMixin(value): { hideFrom+: value },
          hideFrom+:
            {
              '#withLegend': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withLegend(value=true): { hideFrom+: { legend: value } },
              '#withTooltip': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withTooltip(value=true): { hideFrom+: { tooltip: value } },
              '#withViz': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
              withViz(value=true): { hideFrom+: { viz: value } },
            },
          '#withAxisCenteredZero': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withAxisCenteredZero(value=true): { axisCenteredZero: value },
          '#withAxisColorMode': { 'function': { args: [{ default: null, enums: ['text', 'series'], name: 'value', type: 'string' }], help: 'TODO docs' } },
          withAxisColorMode(value): { axisColorMode: value },
          '#withAxisGridShow': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withAxisGridShow(value=true): { axisGridShow: value },
          '#withAxisLabel': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withAxisLabel(value): { axisLabel: value },
          '#withAxisPlacement': { 'function': { args: [{ default: null, enums: ['auto', 'top', 'right', 'bottom', 'left', 'hidden'], name: 'value', type: 'string' }], help: 'TODO docs' } },
          withAxisPlacement(value): { axisPlacement: value },
          '#withAxisSoftMax': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
          withAxisSoftMax(value): { axisSoftMax: value },
          '#withAxisSoftMin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
          withAxisSoftMin(value): { axisSoftMin: value },
          '#withAxisWidth': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
          withAxisWidth(value): { axisWidth: value },
          '#withScaleDistribution': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
          withScaleDistribution(value): { scaleDistribution: value },
          '#withScaleDistributionMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
          withScaleDistributionMixin(value): { scaleDistribution+: value },
          scaleDistribution+:
            {
              '#withLinearThreshold': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
              withLinearThreshold(value): { scaleDistribution+: { linearThreshold: value } },
              '#withLog': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
              withLog(value): { scaleDistribution+: { log: value } },
              '#withType': { 'function': { args: [{ default: null, enums: ['linear', 'log', 'ordinal', 'symlog'], name: 'value', type: 'string' }], help: 'TODO docs' } },
              withType(value): { scaleDistribution+: { type: value } },
            },
          '#withLabel': { 'function': { args: [{ default: null, enums: ['auto', 'never', 'always'], name: 'value', type: 'string' }], help: 'TODO docs' } },
          withLabel(value): { label: value },
          '#withLabelValue': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withLabelValue(value): { labelValue: value },
          '#withLabelValueMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withLabelValueMixin(value): { labelValue+: value },
          labelValue+:
            {
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'fixed: T -- will be added by each element' } },
              withField(value): { labelValue+: { field: value } },
              '#withFixed': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withFixed(value): { labelValue+: { fixed: value } },
              '#withMode': { 'function': { args: [{ default: null, enums: ['fixed', 'field', 'template'], name: 'value', type: 'string' }], help: '' } },
              withMode(value): { labelValue+: { mode: value } },
            },
          '#withLineColor': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withLineColor(value): { lineColor: value },
          '#withLineColorMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withLineColorMixin(value): { lineColor+: value },
          lineColor+:
            {
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'fixed: T -- will be added by each element' } },
              withField(value): { lineColor+: { field: value } },
              '#withFixed': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withFixed(value): { lineColor+: { fixed: value } },
            },
          '#withLineStyle': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
          withLineStyle(value): { lineStyle: value },
          '#withLineStyleMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
          withLineStyleMixin(value): { lineStyle+: value },
          lineStyle+:
            {
              '#withDash': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
              withDash(value): { lineStyle+: { dash: (if std.isArray(value)
                                                      then value
                                                      else [value]) } },
              '#withDashMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: '' } },
              withDashMixin(value): { lineStyle+: { dash+: (if std.isArray(value)
                                                            then value
                                                            else [value]) } },
              '#withFill': { 'function': { args: [{ default: null, enums: ['solid', 'dash', 'dot', 'square'], name: 'value', type: 'string' }], help: '' } },
              withFill(value): { lineStyle+: { fill: value } },
            },
          '#withLineWidth': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: '' } },
          withLineWidth(value): { lineWidth: value },
          '#withPointColor': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withPointColor(value): { pointColor: value },
          '#withPointColorMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withPointColorMixin(value): { pointColor+: value },
          pointColor+:
            {
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'fixed: T -- will be added by each element' } },
              withField(value): { pointColor+: { field: value } },
              '#withFixed': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
              withFixed(value): { pointColor+: { fixed: value } },
            },
          '#withPointSize': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withPointSize(value): { pointSize: value },
          '#withPointSizeMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withPointSizeMixin(value): { pointSize+: value },
          pointSize+:
            {
              '#withField': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'fixed: T -- will be added by each element' } },
              withField(value): { pointSize+: { field: value } },
              '#withFixed': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
              withFixed(value): { pointSize+: { fixed: value } },
              '#withMax': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
              withMax(value): { pointSize+: { max: value } },
              '#withMin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
              withMin(value): { pointSize+: { min: value } },
              '#withMode': { 'function': { args: [{ default: null, enums: ['linear', 'quad'], name: 'value', type: 'string' }], help: '' } },
              withMode(value): { pointSize+: { mode: value } },
            },
          '#withShow': { 'function': { args: [{ default: null, enums: ['points', 'lines', 'points+lines'], name: 'value', type: 'string' }], help: '' } },
          withShow(value): { show: value },
          '#withName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withName(value): { name: value },
          '#withX': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withX(value): { x: value },
          '#withY': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withY(value): { y: value },
        },
      '#withSeriesMapping': { 'function': { args: [{ default: null, enums: ['auto', 'manual'], name: 'value', type: 'string' }], help: '' } },
      withSeriesMapping(value): { options+: { seriesMapping: value } },
    },
}
+ { panelOptions+: { '#withType':: {} } }
