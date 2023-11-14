# preferences

grafonnet.preferences

## Index

* [`fn withHomeDashboardUID(value)`](#fn-withhomedashboarduid)
* [`fn withLanguage(value)`](#fn-withlanguage)
* [`fn withQueryHistory(value)`](#fn-withqueryhistory)
* [`fn withQueryHistoryMixin(value)`](#fn-withqueryhistorymixin)
* [`fn withTheme(value)`](#fn-withtheme)
* [`fn withTimezone(value)`](#fn-withtimezone)
* [`fn withWeekStart(value)`](#fn-withweekstart)
* [`obj queryHistory`](#obj-queryhistory)
  * [`fn withHomeTab(value)`](#fn-queryhistorywithhometab)

## Fields

### fn withHomeDashboardUID

```jsonnet
withHomeDashboardUID(value)
```

PARAMETERS:

* **value** (`string`)

UID for the home dashboard
### fn withLanguage

```jsonnet
withLanguage(value)
```

PARAMETERS:

* **value** (`string`)

Selected language (beta)
### fn withQueryHistory

```jsonnet
withQueryHistory(value)
```

PARAMETERS:

* **value** (`object`)


### fn withQueryHistoryMixin

```jsonnet
withQueryHistoryMixin(value)
```

PARAMETERS:

* **value** (`object`)


### fn withTheme

```jsonnet
withTheme(value)
```

PARAMETERS:

* **value** (`string`)

light, dark, empty is default
### fn withTimezone

```jsonnet
withTimezone(value)
```

PARAMETERS:

* **value** (`string`)

The timezone selection
TODO: this should use the timezone defined in common
### fn withWeekStart

```jsonnet
withWeekStart(value)
```

PARAMETERS:

* **value** (`string`)

day of the week (sunday, monday, etc)
### obj queryHistory


#### fn queryHistory.withHomeTab

```jsonnet
queryHistory.withHomeTab(value)
```

PARAMETERS:

* **value** (`string`)

one of: '' | 'query' | 'starred';