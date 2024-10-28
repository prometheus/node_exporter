local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local genericPanel = import 'panel.libsonnet';
genericPanel {
  new(
    title=null,
    description=null,
    datasource=null,
  ):: self +
      grafana.statPanel.new(
        title=title,
        description=description,
        datasource=datasource,
      ),
  withGraphMode(mode='none'):: self {
    options+:
      {
        graphMode: mode,
      },
  },
  withTextSize(value='auto', title='auto'):: self {
    options+:
      { text: {
        valueSize: value,
        titleSize: title,
      } },
  },

}
