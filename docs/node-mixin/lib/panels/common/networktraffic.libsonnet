// Panels to graph network traffic in and out
local timeseries = import '../timeseries.libsonnet';
timeseries {
  new(
    title=null,
    description=null,
    datasource=null,
  )::
    super.new(
      title,
      description,
      datasource,
    )
    + self.withDecimals(1)
    + self.withUnits('bps')
    + self.withNegativeYByRegex('transmit|tx|out')
    + self.withAxisLabel('out(-) | in(+)'),
}
