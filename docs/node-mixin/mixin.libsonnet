local nodelib = import './lib/linux/main.libsonnet';

{
  _config:: {},
  _linuxLib::
    nodelib.new()
    + nodelib.withConfigMixin(self._config),
  grafanaDashboards+:: self._linuxLib.grafana.dashboards,
  prometheusAlerts+:: self._linuxLib.prometheus.alerts,
  prometheusRules+:: self._linuxLib.prometheus.recordingRules,
}
