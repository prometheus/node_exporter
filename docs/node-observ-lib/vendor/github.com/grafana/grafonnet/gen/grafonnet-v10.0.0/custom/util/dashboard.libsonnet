local d = import 'github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet';
local xtd = import 'github.com/jsonnet-libs/xtd/main.libsonnet';

{
  local root = self,

  '#getOptionsForCustomQuery':: d.func.new(
    |||
      `getOptionsForCustomQuery` provides values for the `options` and `current` fields.
      These are required for template variables of type 'custom'but do not automatically
      get populated by Grafana when importing a dashboard from JSON.

      This is a bit of a hack and should always be called on functions that set `type` on
      a template variable. Ideally Grafana populates these fields from the `query` value
      but this provides a backwards compatible solution.
    |||,
    args=[d.arg('query', d.T.string)],
  ),
  getOptionsForCustomQuery(query, multi): {
    local values = root.parseCustomQuery(query),
    current: root.getCurrentFromValues(values, multi),
    options: root.getOptionsFromValues(values),
  },

  getCurrentFromValues(values, multi): {
    selected: false,
    text: if multi then [values[0].key] else values[0].key,
    value: if multi then [values[0].value] else values[0].value,
  },

  getOptionsFromValues(values):
    std.mapWithIndex(
      function(i, item) {
        selected: i == 0,
        text: item.key,
        value: item.value,
      },
      values
    ),

  parseCustomQuery(query):
    std.map(
      function(v)
        // Split items into key:value pairs
        local split = std.splitLimit(v, ' : ', 1);
        {
          key: std.stripChars(split[0], ' '),
          value:
            if std.length(split) == 2
            then std.stripChars(split[1], ' ')
            else self.key,
        },
      xtd.string.splitEscape(query, ',')  // Split query by comma, unless the comma is escaped
    ),
}
