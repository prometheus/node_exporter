local config = import './config.libsonnet';
local g = import './g.libsonnet';
local panels = import './panels.libsonnet';
local targets = import './targets.libsonnet';
local nodelib = import 'node-observ-lib/main.libsonnet';


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
          backToFleet:: {},
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
            overview: parentGrafana.dashboards.overview,
          }
          +
          (
            if this.config.enableLokiLogs
            then
              {
                logs: parentGrafana.dashboards.logs,
              }
          ),
      },
      prometheus+: {
        recordingRules: {},
        alerts:
          {
            groups:
              [
                {
                  name: group.name,
                  rules: [
                    rule
                    for rule in group.rules
                    if std.length(std.find(rule.alert, this.config.alertsMacKeep)) > 0
                  ],
                }
                for group in parentPrometheus.alerts.groups
              ],

          },

      },
    },

}
