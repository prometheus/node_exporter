local macoslib = import './lib/macos/main.libsonnet';
{
  _config:: {},
  _macosLib::
    macoslib.new()
    + macoslib.withConfigMixin(self._config),
  grafanaDashboards+:: self._macosLib.grafana.dashboards,
  prometheusAlerts+:: self._macosLib.prometheus.alerts,
  prometheusRules+:: self._macosLib.prometheus.recordingRules,
}
