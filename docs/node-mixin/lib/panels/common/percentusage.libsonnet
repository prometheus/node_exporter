// Panels to display metrics that can go from 0 to 100%. (cpu utilization, memory utilization etc). Full utilization is considered an issue.
local statPanel = import '../stat.libsonnet';
statPanel {
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
    + self.withUnits("percent")
    + self.withMax(100)
    + self.withMin(0)
    + self.withColor(mode="continuous-BlYlRd")
    {
      options+: {
        "reduceOptions": {
        "values": false,
        "calcs": [
            "lastNotNull"
        ],
        "fields": ""
        },
      }
    }
}
