{
  new(panel)::
    panel
    {
      type: 'timeseries',
      options+: {
        tooltip: {
          mode: 'multi',
        },
      },
    }
    + self.withFillOpacity(10)
    + self.withGradientMode('opacity')
    + self.withLineInterpolation('smooth')
    + self.withShowPoints('never')
    + self.withTooltip(mode='multi', sort='none'),
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
  withMin(value):: self {
    fieldConfig+: {
      defaults+: {
        min: value,
      },
    },
  },
  withMax(value):: self {
    fieldConfig+: {
      defaults+: {
        max: value,
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

  withUnit(unit):: self {

    fieldConfig+: {
      defaults+: {
        unit: unit,
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
}
