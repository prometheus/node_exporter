local aixlib = import './lib/aix/main.libsonnet';
{
  _config:: {},
  _aixLib::
    aixlib.new()
    + aixlib.withConfigMixin(self._config),
  grafanaDashboards+:: self._aixLib.grafana.dashboards,
  prometheusAlerts+:: self._aixLib.prometheus.alerts,
  prometheusRules+:: self._aixLib.prometheus.recordingRules,
}
