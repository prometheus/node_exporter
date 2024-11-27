local alerts = import './alerts/alerts.libsonnet';
local annotations = import './annotations.libsonnet';
local config = import './config.libsonnet';
local dashboards = import './dashboards.libsonnet';
local datasources = import './datasources.libsonnet';
local g = import './g.libsonnet';
local links = import './links.libsonnet';
local panels = import './panels/main.libsonnet';
local rows = import './rows/main.libsonnet';
local rules = import './rules/rules.libsonnet';
local targets = import './targets/main.libsonnet';
local variables = import './variables.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';

{
  new(): {

    local this = self,
    config: config,
    grafana: {
      variables: variables.new(this),
      targets: targets.new(this),
      panels: panels.new(this),
      annotations: annotations.new(this),
      // common links here used across all dashboards
      links: links.new(this),
      rows: rows.new(this),
      dashboards: dashboards.new(this),
    },

    prometheus: {
      alerts: alerts.new(this),
      recordingRules: rules.new(this),
    },


  },
  withConfigMixin(config): {
    //backward compatible:  handle both formats string and array for instanceLabels, groupLabels
    local _patch =
      (
        if std.objectHasAll(config, 'instanceLabels')
        then
          { instanceLabels: if std.isString(config.instanceLabels) then std.split(',', config.instanceLabels) else config.instanceLabels }
        else {}
      ) +
      (
        if std.objectHasAll(config, 'groupLabels')
        then
          {
            groupLabels: if std.isString(config.groupLabels) then std.split(',', config.groupLabels) else config.groupLabels,
          }
        else {}
      ),
    local groupLabels = if std.isString(config.groupLabels) then std.split(',', config.groupLabels) else config.groupLabels,
    config+: config + _patch,
  },
}
