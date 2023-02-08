// generic grafana dashboard
{
  //feed grafonnet panel
  new()::{},

  withUnits(unit):: self {

    fieldConfig+: {
      defaults+: {
        unit: unit
      },
    },
  },

  withLegend(show=true, mode="table", placement="bottom", calcs=["min","mean","max","lastNotNull"]):: self {
    options+: {
      "legend": {
        "showLegend": show,
        "displayMode": mode,
        "placement": placement,
        "calcs": calcs,
      }
    }
  },
  withDecimals(decimals):: self {

    fieldConfig+: {
      defaults+: {
        decimals: decimals,
      },
    },
  },

  withThresholds(mode="absolute", steps=null):: self {

    fieldConfig+: {
      defaults+: {
        "thresholds": {
            "mode": mode,
            "steps": steps
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
  withColor(color=null, mode="fixed"):: self {
        fieldConfig+: {
        defaults+: {
            "color": {
                mode: mode,
                fixedColor: if mode == 'fixed' then color else null
            },
        },
    },
  },

  withTransform():: self {
    
    merge():: self
    {
      transformations+: [
        {
          "id": "merge",
          "options": {}
        },
      ]
    },
    filterFieldsByName(pattern=null):: self
    {
      transformations+: [
          {
            "id": "filterFieldsByName",
            "options": {
              "include": {
                "pattern": pattern
              }
            }
          }
      ]
    },
    joinByField(
      mode="outer",
      field=null
    ):: self {
      transformations+: [
        {
          "id": "joinByField",
          "options": {
            "byField": field,
            "mode": mode
          }
        },
      ]
    },
    organize(
      excludeByName={},
      indexByName={},
      renameByName={},

      ):: self
    {
      transformations+: [
        {
          "id": "organize",
          "options": {
            "excludeByName": excludeByName,
            "indexByName": indexByName,
            "renameByName": renameByName,
          }
        },
      ]
    }
  },
}
