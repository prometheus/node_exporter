local macoslib = import './macos/main.libsonnet';
local macos = macoslib.new();

{
  grafanaDashboards+:: macos.grafana.dashboards,
  prometheusAlerts+:: macos.prometheus.alerts,
  prometheusRules+:: macos.prometheus.recordingRules,
}
