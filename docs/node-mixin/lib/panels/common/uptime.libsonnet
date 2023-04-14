local statPanel = import '../stat.libsonnet';
statPanel {
  new(
    title='Uptime',
    description=null,
    datasource=null,
  )::
    super.new(
      title,
      description,
      datasource,
    )
    + self.withDecimals(1)
    + self.withGraphMode('none')
    + self.withTextSize(value=20)
    + self.withUnits('dtdurations')
    + self.withThresholds(
      mode='absolute',
      steps=[
        {
          color: 'orange',
          value: null,
        },
        {
          color: 'text',
          value: 300,
        },
      ]
    )
    + self.withColor(mode='thresholds')
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
      },
    },
}
