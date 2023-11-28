local alerts = import './alerts.libsonnet';
local annotations = import './annotations.libsonnet';
local config = import './config.libsonnet';
local dashboards = import './dashboards.libsonnet';
local datasources = import './datasources.libsonnet';
local g = import './g.libsonnet';
local links = import './links.libsonnet';
local panels = import './panels.libsonnet';
local rules = import './rules.libsonnet';
local targets = import './targets.libsonnet';
local variables = import './variables.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';

{
  withConfigMixin(config): {
    config+: config,
  },

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
      dashboards: dashboards.new(this),
    },

    prometheus: {
      alerts: alerts.new(this),
      recordingRules: rules.new(this),
    },

  },
}
