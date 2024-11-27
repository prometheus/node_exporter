local g = import '../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
{
  new(this):
    {
      local link = g.dashboard.link,
      backToFleet:
        link.link.new('Back to ' + this.config.dashboardNamePrefix + 'fleet', '/d/' + this.grafana.dashboards['fleet.json'].uid)
        + link.link.options.withKeepTime(true),
      backToOverview:
        link.link.new('Back to ' + this.config.dashboardNamePrefix + 'overview', '/d/' + this.grafana.dashboards['nodes.json'].uid)
        + link.link.options.withKeepTime(true),
      otherDashboards:
        link.dashboards.new('All ' + this.config.dashboardNamePrefix + ' dashboards', this.config.dashboardTags)
        + link.dashboards.options.withIncludeVars(true)
        + link.dashboards.options.withKeepTime(true)
        + link.dashboards.options.withAsDropdown(true),
    },
}
