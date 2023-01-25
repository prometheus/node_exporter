local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local grafana70 = import 'github.com/grafana/grafonnet-lib/grafonnet-7.0/grafana.libsonnet';
local gaugePanel = grafana70.panel.gauge;
local table = grafana70.panel.table;
local nodePanels = import 'panels.libsonnet';
local nodeTimeseries = nodePanels.nodeTimeseries;
{

  new(config=null, platform=null):: {

    local prometheusDatasourceTemplate = {
      current: {
        text: 'default',
        value: 'default',
      },
      hide: 0,
      label: 'Data Source',
      name: 'datasource',
      options: [],
      query: 'prometheus',
      refresh: 1,
      regex: '',
      type: 'datasource',
    },

    local instanceTemplatePrototype =
      template.new(
        'instance',
        '$datasource',
        '',
        refresh='time',
        label='Instance',
      ),
    local instanceTemplate =
      if platform == 'Darwin' then
        instanceTemplatePrototype
        { query: 'label_values(node_uname_info{%(nodeExporterSelector)s, sysname="Darwin"}, instance)' % config }
      else
        instanceTemplatePrototype
        { query: 'label_values(node_uname_info{%(nodeExporterSelector)s, sysname!="Darwin"}, instance)' % config },

    // return common templates
    templates: [
      prometheusDatasourceTemplate,
      instanceTemplate,
    ],
  },


}
