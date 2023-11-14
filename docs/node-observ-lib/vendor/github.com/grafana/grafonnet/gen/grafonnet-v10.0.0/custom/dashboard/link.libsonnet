local d = import 'github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet';

{
  '#withLinks':: d.func.new(
    |||
      Dashboard links are displayed at the top of the dashboard, these can either link to other dashboards or to external URLs.

      `withLinks` takes an array of [link objects](./link.md).

      The [docs](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/manage-dashboard-links/#dashboard-links) give a more comprehensive description.

      Example:

      ```jsonnet
      local g = import 'g.libsonnet';
      local link = g.dashboard.link;

      g.dashboard.new('Title dashboard')
      + g.dashboard.withLinks([
        link.link.new('My title', 'https://wikipedia.org/'),
      ])
      ```
    |||,
    [d.arg('value', d.T.array)],
  ),
  '#withLinksMixin':: self['#withLinks'],

  link+: {
    '#':: d.package.newSub(
      'link',
      |||
        Dashboard links are displayed at the top of the dashboard, these can either link to other dashboards or to external URLs.

        The [docs](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/manage-dashboard-links/#dashboard-links) give a more comprehensive description.

        Example:

        ```jsonnet
        local g = import 'g.libsonnet';
        local link = g.dashboard.link;

        g.dashboard.new('Title dashboard')
        + g.dashboard.withLinks([
          link.link.new('My title', 'https://wikipedia.org/'),
        ])
        ```
      |||,
    ),

    dashboards+: {
      '#new':: d.func.new(
        |||
          Create links to dashboards based on `tags`.
        |||,
        args=[
          d.arg('title', d.T.string),
          d.arg('tags', d.T.array),
        ]
      ),
      new(title, tags):
        self.withTitle(title)
        + self.withType('dashboards')
        + self.withTags(tags),

      '#withTitle':: {},
      '#withType':: {},
      '#withTags':: {},
    },

    link+: {
      '#new':: d.func.new(
        |||
          Create link to an arbitrary URL.
        |||,
        args=[
          d.arg('title', d.T.string),
          d.arg('url', d.T.string),
        ]
      ),
      new(title, url):
        self.withTitle(title)
        + self.withType('link')
        + self.withUrl(url),

      '#withTitle':: {},
      '#withType':: {},
      '#withUrl':: {},
    },
  },
}
