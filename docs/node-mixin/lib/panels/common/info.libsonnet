// Info panel text (number or text)
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
    + self.withColor(color='text')
    + self.withTextSize(value=30)
    + self.withGraphMode('none')
    +
    {
      options+: {
        reduceOptions: {
          values: false,
          calcs: [
            'lastNotNull',
          ],
          fields: '',
        },
        graphMode: 'none',
      },
    },
}
