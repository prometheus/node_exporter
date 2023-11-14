// This file is generated, do not manually edit.
(import '../../clean/panel.libsonnet')
+ {
  '#': { help: 'grafonnet.panel.heatmap', name: 'heatmap' },
  panelOptions+:
    {
      '#withType': { 'function': { args: [], help: '' } },
      withType(): { type: 'heatmap' },
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
            },
        },
    },
  options+:
    {
      '#withCalculate': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Controls if the heatmap should be calculated from data' } },
      withCalculate(value=true): { options+: { calculate: value } },
      '#withCalculation': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withCalculation(value): { options+: { calculation: value } },
      '#withCalculationMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
      withCalculationMixin(value): { options+: { calculation+: value } },
      calculation+:
        {
          '#withXBuckets': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withXBuckets(value): { options+: { calculation+: { xBuckets: value } } },
          '#withXBucketsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withXBucketsMixin(value): { options+: { calculation+: { xBuckets+: value } } },
          xBuckets+:
            {
              '#withMode': { 'function': { args: [{ default: null, enums: ['size', 'count'], name: 'value', type: 'string' }], help: '' } },
              withMode(value): { options+: { calculation+: { xBuckets+: { mode: value } } } },
              '#withScale': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
              withScale(value): { options+: { calculation+: { xBuckets+: { scale: value } } } },
              '#withScaleMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
              withScaleMixin(value): { options+: { calculation+: { xBuckets+: { scale+: value } } } },
              scale+:
                {
                  '#withLinearThreshold': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
                  withLinearThreshold(value): { options+: { calculation+: { xBuckets+: { scale+: { linearThreshold: value } } } } },
                  '#withLog': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
                  withLog(value): { options+: { calculation+: { xBuckets+: { scale+: { log: value } } } } },
                  '#withType': { 'function': { args: [{ default: null, enums: ['linear', 'log', 'ordinal', 'symlog'], name: 'value', type: 'string' }], help: 'TODO docs' } },
                  withType(value): { options+: { calculation+: { xBuckets+: { scale+: { type: value } } } } },
                },
              '#withValue': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'The number of buckets to use for the axis in the heatmap' } },
              withValue(value): { options+: { calculation+: { xBuckets+: { value: value } } } },
            },
          '#withYBuckets': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withYBuckets(value): { options+: { calculation+: { yBuckets: value } } },
          '#withYBucketsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: '' } },
          withYBucketsMixin(value): { options+: { calculation+: { yBuckets+: value } } },
          yBuckets+:
            {
              '#withMode': { 'function': { args: [{ default: null, enums: ['size', 'count'], name: 'value', type: 'string' }], help: '' } },
              withMode(value): { options+: { calculation+: { yBuckets+: { mode: value } } } },
              '#withScale': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
              withScale(value): { options+: { calculation+: { yBuckets+: { scale: value } } } },
              '#withScaleMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
              withScaleMixin(value): { options+: { calculation+: { yBuckets+: { scale+: value } } } },
              scale+:
                {
                  '#withLinearThreshold': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
                  withLinearThreshold(value): { options+: { calculation+: { yBuckets+: { scale+: { linearThreshold: value } } } } },
                  '#withLog': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
                  withLog(value): { options+: { calculation+: { yBuckets+: { scale+: { log: value } } } } },
                  '#withType': { 'function': { args: [{ default: null, enums: ['linear', 'log', 'ordinal', 'symlog'], name: 'value', type: 'string' }], help: 'TODO docs' } },
                  withType(value): { options+: { calculation+: { yBuckets+: { scale+: { type: value } } } } },
                },
              '#withValue': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'The number of buckets to use for the axis in the heatmap' } },
              withValue(value): { options+: { calculation+: { yBuckets+: { value: value } } } },
            },
        },
      '#withCellGap': { 'function': { args: [{ default: 1, enums: null, name: 'value', type: 'integer' }], help: 'Controls gap between cells' } },
      withCellGap(value=1): { options+: { cellGap: value } },
      '#withCellRadius': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'Controls cell radius' } },
      withCellRadius(value): { options+: { cellRadius: value } },
      '#withCellValues': { 'function': { args: [{ default: {}, enums: null, name: 'value', type: 'object' }], help: 'Controls cell value unit' } },
      withCellValues(value={}): { options+: { cellValues: value } },
      '#withCellValuesMixin': { 'function': { args: [{ default: {}, enums: null, name: 'value', type: 'object' }], help: 'Controls cell value unit' } },
      withCellValuesMixin(value): { options+: { cellValues+: value } },
      cellValues+:
        {
          '#withCellValues': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Controls cell value options' } },
          withCellValues(value): { options+: { cellValues+: { CellValues: value } } },
          '#withCellValuesMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Controls cell value options' } },
          withCellValuesMixin(value): { options+: { cellValues+: { CellValues+: value } } },
          CellValues+:
            {
              '#withDecimals': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'Controls the number of decimals for cell values' } },
              withDecimals(value): { options+: { cellValues+: { decimals: value } } },
              '#withUnit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Controls the cell value unit' } },
              withUnit(value): { options+: { cellValues+: { unit: value } } },
            },
        },
      '#withColor': { 'function': { args: [{ default: { exponent: 0.5, fill: 'dark-orange', reverse: false, scheme: 'Oranges', steps: 64 }, enums: null, name: 'value', type: 'object' }], help: 'Controls the color options' } },
      withColor(value={ exponent: 0.5, fill: 'dark-orange', reverse: false, scheme: 'Oranges', steps: 64 }): { options+: { color: value } },
      '#withColorMixin': { 'function': { args: [{ default: { exponent: 0.5, fill: 'dark-orange', reverse: false, scheme: 'Oranges', steps: 64 }, enums: null, name: 'value', type: 'object' }], help: 'Controls the color options' } },
      withColorMixin(value): { options+: { color+: value } },
      color+:
        {
          '#withHeatmapColorOptions': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Controls various color options' } },
          withHeatmapColorOptions(value): { options+: { color+: { HeatmapColorOptions: value } } },
          '#withHeatmapColorOptionsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Controls various color options' } },
          withHeatmapColorOptionsMixin(value): { options+: { color+: { HeatmapColorOptions+: value } } },
          HeatmapColorOptions+:
            {
              '#withExponent': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'Controls the exponent when scale is set to exponential' } },
              withExponent(value): { options+: { color+: { exponent: value } } },
              '#withFill': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Controls the color fill when in opacity mode' } },
              withFill(value): { options+: { color+: { fill: value } } },
              '#withMax': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'Sets the maximum value for the color scale' } },
              withMax(value): { options+: { color+: { max: value } } },
              '#withMin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'Sets the minimum value for the color scale' } },
              withMin(value): { options+: { color+: { min: value } } },
              '#withMode': { 'function': { args: [{ default: null, enums: ['opacity', 'scheme'], name: 'value', type: 'string' }], help: 'Controls the color mode of the heatmap' } },
              withMode(value): { options+: { color+: { mode: value } } },
              '#withReverse': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Reverses the color scheme' } },
              withReverse(value=true): { options+: { color+: { reverse: value } } },
              '#withScale': { 'function': { args: [{ default: null, enums: ['linear', 'exponential'], name: 'value', type: 'string' }], help: 'Controls the color scale of the heatmap' } },
              withScale(value): { options+: { color+: { scale: value } } },
              '#withScheme': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Controls the color scheme used' } },
              withScheme(value): { options+: { color+: { scheme: value } } },
              '#withSteps': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'Controls the number of color steps' } },
              withSteps(value): { options+: { color+: { steps: value } } },
            },
        },
      '#withExemplars': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Controls exemplar options' } },
      withExemplars(value): { options+: { exemplars: value } },
      '#withExemplarsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Controls exemplar options' } },
      withExemplarsMixin(value): { options+: { exemplars+: value } },
      exemplars+:
        {
          '#withColor': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Sets the color of the exemplar markers' } },
          withColor(value): { options+: { exemplars+: { color: value } } },
        },
      '#withFilterValues': { 'function': { args: [{ default: { le: 1.0000000000000001e-09 }, enums: null, name: 'value', type: 'object' }], help: 'Filters values between a given range' } },
      withFilterValues(value={ le: 1.0000000000000001e-09 }): { options+: { filterValues: value } },
      '#withFilterValuesMixin': { 'function': { args: [{ default: { le: 1.0000000000000001e-09 }, enums: null, name: 'value', type: 'object' }], help: 'Filters values between a given range' } },
      withFilterValuesMixin(value): { options+: { filterValues+: value } },
      filterValues+:
        {
          '#withFilterValueRange': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Controls the value filter range' } },
          withFilterValueRange(value): { options+: { filterValues+: { FilterValueRange: value } } },
          '#withFilterValueRangeMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Controls the value filter range' } },
          withFilterValueRangeMixin(value): { options+: { filterValues+: { FilterValueRange+: value } } },
          FilterValueRange+:
            {
              '#withGe': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'Sets the filter range to values greater than or equal to the given value' } },
              withGe(value): { options+: { filterValues+: { ge: value } } },
              '#withLe': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'Sets the filter range to values less than or equal to the given value' } },
              withLe(value): { options+: { filterValues+: { le: value } } },
            },
        },
      '#withLegend': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Controls legend options' } },
      withLegend(value): { options+: { legend: value } },
      '#withLegendMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Controls legend options' } },
      withLegendMixin(value): { options+: { legend+: value } },
      legend+:
        {
          '#withShow': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Controls if the legend is shown' } },
          withShow(value=true): { options+: { legend+: { show: value } } },
        },
      '#withRowsFrame': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Controls frame rows options' } },
      withRowsFrame(value): { options+: { rowsFrame: value } },
      '#withRowsFrameMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Controls frame rows options' } },
      withRowsFrameMixin(value): { options+: { rowsFrame+: value } },
      rowsFrame+:
        {
          '#withLayout': { 'function': { args: [{ default: null, enums: ['le', 'ge', 'unknown', 'auto'], name: 'value', type: 'string' }], help: '' } },
          withLayout(value): { options+: { rowsFrame+: { layout: value } } },
          '#withValue': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Sets the name of the cell when not calculating from data' } },
          withValue(value): { options+: { rowsFrame+: { value: value } } },
        },
      '#withShowValue': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '| *{\n\tlayout: ui.HeatmapCellLayout & "auto" // TODO: fix after remove when https://github.com/grafana/cuetsy/issues/74 is fixed\n}\nControls the display of the value in the cell' } },
      withShowValue(value): { options+: { showValue: value } },
      '#withTooltip': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Controls tooltip options' } },
      withTooltip(value): { options+: { tooltip: value } },
      '#withTooltipMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Controls tooltip options' } },
      withTooltipMixin(value): { options+: { tooltip+: value } },
      tooltip+:
        {
          '#withShow': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Controls if the tooltip is shown' } },
          withShow(value=true): { options+: { tooltip+: { show: value } } },
          '#withYHistogram': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Controls if the tooltip shows a histogram of the y-axis values' } },
          withYHistogram(value=true): { options+: { tooltip+: { yHistogram: value } } },
        },
      '#withYAxis': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Configuration options for the yAxis' } },
      withYAxis(value): { options+: { yAxis: value } },
      '#withYAxisMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'Configuration options for the yAxis' } },
      withYAxisMixin(value): { options+: { yAxis+: value } },
      yAxis+:
        {
          '#withAxisCenteredZero': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withAxisCenteredZero(value=true): { options+: { yAxis+: { axisCenteredZero: value } } },
          '#withAxisColorMode': { 'function': { args: [{ default: null, enums: ['text', 'series'], name: 'value', type: 'string' }], help: 'TODO docs' } },
          withAxisColorMode(value): { options+: { yAxis+: { axisColorMode: value } } },
          '#withAxisGridShow': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: '' } },
          withAxisGridShow(value=true): { options+: { yAxis+: { axisGridShow: value } } },
          '#withAxisLabel': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: '' } },
          withAxisLabel(value): { options+: { yAxis+: { axisLabel: value } } },
          '#withAxisPlacement': { 'function': { args: [{ default: null, enums: ['auto', 'top', 'right', 'bottom', 'left', 'hidden'], name: 'value', type: 'string' }], help: 'TODO docs' } },
          withAxisPlacement(value): { options+: { yAxis+: { axisPlacement: value } } },
          '#withAxisSoftMax': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
          withAxisSoftMax(value): { options+: { yAxis+: { axisSoftMax: value } } },
          '#withAxisSoftMin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
          withAxisSoftMin(value): { options+: { yAxis+: { axisSoftMin: value } } },
          '#withAxisWidth': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
          withAxisWidth(value): { options+: { yAxis+: { axisWidth: value } } },
          '#withScaleDistribution': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
          withScaleDistribution(value): { options+: { yAxis+: { scaleDistribution: value } } },
          '#withScaleDistributionMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'TODO docs' } },
          withScaleDistributionMixin(value): { options+: { yAxis+: { scaleDistribution+: value } } },
          scaleDistribution+:
            {
              '#withLinearThreshold': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
              withLinearThreshold(value): { options+: { yAxis+: { scaleDistribution+: { linearThreshold: value } } } },
              '#withLog': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: '' } },
              withLog(value): { options+: { yAxis+: { scaleDistribution+: { log: value } } } },
              '#withType': { 'function': { args: [{ default: null, enums: ['linear', 'log', 'ordinal', 'symlog'], name: 'value', type: 'string' }], help: 'TODO docs' } },
              withType(value): { options+: { yAxis+: { scaleDistribution+: { type: value } } } },
            },
          '#withDecimals': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'Controls the number of decimals for yAxis values' } },
          withDecimals(value): { options+: { yAxis+: { decimals: value } } },
          '#withMax': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'Sets the maximum value for the yAxis' } },
          withMax(value): { options+: { yAxis+: { max: value } } },
          '#withMin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'number' }], help: 'Sets the minimum value for the yAxis' } },
          withMin(value): { options+: { yAxis+: { min: value } } },
          '#withReverse': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'Reverses the yAxis' } },
          withReverse(value=true): { options+: { yAxis+: { reverse: value } } },
          '#withUnit': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Sets the yAxis unit' } },
          withUnit(value): { options+: { yAxis+: { unit: value } } },
        },
    },
}
+ { panelOptions+: { '#withType':: {} } }
