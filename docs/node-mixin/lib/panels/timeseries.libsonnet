local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local genericPanel = import 'panel.libsonnet';
genericPanel
{
  new(
    title=null,
    description=null,
    datasource=null,
  ):: self +
    grafana.graphPanel.new(
      title=title,
      description=description,
      datasource=datasource,
    )
    +
    {
      type: 'timeseries',
    }
    + self.withFillOpacity(10)
    + self.withGradientMode('opacity')
    + self.withLineInterpolation('smooth')
    + self.withShowPoints('never')
    + self.withTooltip(mode='multi', sort='none')
    + self.withLegend(mode='list', calcs=[]),

  withTooltip(mode=null, sort='none'):: self {
    options+: {
      tooltip: {
        mode: 'multi',
        sort: 'none',
      },
    },
  },
  withLineInterpolation(value):: self {
    fieldConfig+: {
      defaults+: {
        custom+: {
          lineInterpolation: value,
        },
      },
    },
  },
  withShowPoints(value):: self {
    fieldConfig+: {
      defaults+: {
        custom+: {
          showPoints: value,
        },
      },
    },
  },
  withStacking(stack):: self {
    fieldConfig+: {
      defaults+: {
        custom+: {
          stacking: {
            mode: stack,
            group: 'A',
          },
        },
      },
    },
  },
  withGradientMode(mode):: self {
    fieldConfig+: {
      defaults+: {
        custom+: {
          gradientMode: mode,
        },
      },
    },
  },
  addDataLink(title, url):: self {

    fieldConfig+: {
      defaults+: {
        links: [
          {
            title: title,
            url: url,
          },
        ],
      },
    },
  },

  withFillOpacity(opacity):: self {
    fieldConfig+: {
      defaults+: {
        custom+: {
          fillOpacity: opacity,
        },
      },
    },

  },

  withAxisLabel(label):: self {
    fieldConfig+: {
      defaults+: {
        custom+: {
          axisLabel: label,
        },
      },
    },
  },

  withNegativeYByRegex(regex):: self {
      fieldConfig+: {
        overrides+: [
{          matcher: {
            id: 'byRegexp',
            options: '/'+regex+'/',
          },
          properties: [
            {
              "id": "custom.transform",
              "value": "negative-Y"
            },
          ]}

        ]
      }
      

  }
}
