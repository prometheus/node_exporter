local genericPanel = import 'panel.libsonnet';
local grafana70 = import 'github.com/grafana/grafonnet-lib/grafonnet-7.0/grafana.libsonnet';
local table = grafana70.panel.table;
genericPanel
{
  new(
    title=null,
    description=null,
    datasource=null,
  ):: self +
    table.new(
      title=title,
      description=description,
      datasource=datasource,
    )
    
}
