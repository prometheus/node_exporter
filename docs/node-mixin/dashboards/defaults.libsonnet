{
  local grafanaDashboards = super.grafanaDashboards,
  grafanaDashboards::
    {
      [fname]: grafanaDashboards[fname] { uid: std.md5(fname) }
      for fname in std.objectFields(grafanaDashboards)
    },
}
