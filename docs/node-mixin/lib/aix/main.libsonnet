local g = import '../g.libsonnet';
local nodelib = import '../linux/main.libsonnet';
local alerts = import './alerts.libsonnet';
local config = import './config.libsonnet';
local panels = import './panels.libsonnet';
local targets = import './targets.libsonnet';


// inherit nodelib
nodelib
{

  new():
    super.new()
    + nodelib.withConfigMixin(config)
    +
    {
      local this = self,
      local parentGrafana = super.grafana,
      local parentPrometheus = super.prometheus,

      grafana+: {
        // drop backToFleet link
        links+: {
          local link = g.dashboard.link,
          backToFleet:: {},
          backToOverview:
            link.link.new('Back to ' + this.config.dashboardNamePrefix + 'overview', '/d/' + this.grafana.dashboards['nodes-darwin.json'].uid)
            + link.link.options.withKeepTime(true),
        },
        annotations: {
          // keep only reboot annotation
          reboot: parentGrafana.annotations.reboot,
        },
        // override targets (memory)
        targets+: targets.new(this),
        // override panels (update description and targets in panels)
        panels+: panels.new(this),

        // keep only overview and logs(optionally) dashes
        dashboards:
          {
            'nodes-aix.json':
              parentGrafana.dashboards['nodes.json']
              + g.dashboard.withUid(
                (if this.config.uid == 'aix' then std.md5('nodes-aix.json') else this.config.uid + '-overview')
              ),
          }
          +
          (
            if this.config.enableLokiLogs
            then
              {
                'logs-aix.json': parentGrafana.dashboards['logs.json'],
              }
            else {}
          ),
      },
      prometheus+: {
        recordingRules: {},
        alerts: alerts.new(this, parentPrometheus),
      },
    },

}
