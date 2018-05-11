local grafana = import 'grafonnet/grafana.libsonnet';
local singlestat = grafana.singlestat;
local prometheus = grafana.prometheus;

{
  new(title, query)::
    singlestat.new(
      title,
      datasource='$datasource',
      span=3,
      format='percentunit',
      valueName='current',
      colors=[
        'rgba(245, 54, 54, 0.9)',
        'rgba(237, 129, 40, 0.89)',
        'rgba(50, 172, 45, 0.97)',
      ],
      thresholds='50, 80',
      valueMaps=[
        {
          op: '=',
          text: 'N/A',
          value: 'null',
        },
      ],
    )
    .addTarget(
      prometheus.target(
        query
      )
    ) + {
      gauge: {
        maxValue: 100,
        minValue: 0,
        show: true,
        thresholdLabels: false,
        thresholdMarkers: true,
      },
      withTextNullValue(text):: self {
        valueMaps: [
          {
            op: '=',
            text: text,
            value: 'null',
          },
        ],
      },
      withSpanSize(size):: self {
        span: size,
      },
      withLowerBeingBetter():: self {
        colors: [
          'rgba(50, 172, 45, 0.97)',
          'rgba(237, 129, 40, 0.89)',
          'rgba(245, 54, 54, 0.9)',
        ],
        thresholds: '80, 90',
      },
    },
}
