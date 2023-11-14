# contactPoint

grafonnet.alerting.contactPoint

**NOTE**: The schemas for all different contact points is under development, this means we can't properly express them in Grafonnet yet. The way this works now may change heavily.


## Index

* [`fn withDisableResolveMessage(value=true)`](#fn-withdisableresolvemessage)
* [`fn withName(value)`](#fn-withname)
* [`fn withProvenance(value)`](#fn-withprovenance)
* [`fn withSettings(value)`](#fn-withsettings)
* [`fn withSettingsMixin(value)`](#fn-withsettingsmixin)
* [`fn withType(value)`](#fn-withtype)
* [`fn withUid(value)`](#fn-withuid)

## Fields

### fn withDisableResolveMessage

```jsonnet
withDisableResolveMessage(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


### fn withName

```jsonnet
withName(value)
```

PARAMETERS:

* **value** (`string`)

Name is used as grouping key in the UI. Contact points with the
same name will be grouped in the UI.
### fn withProvenance

```jsonnet
withProvenance(value)
```

PARAMETERS:

* **value** (`string`)


### fn withSettings

```jsonnet
withSettings(value)
```

PARAMETERS:

* **value** (`object`)


### fn withSettingsMixin

```jsonnet
withSettingsMixin(value)
```

PARAMETERS:

* **value** (`object`)


### fn withType

```jsonnet
withType(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"alertmanager"`, `" dingding"`, `" discord"`, `" email"`, `" googlechat"`, `" kafka"`, `" line"`, `" opsgenie"`, `" pagerduty"`, `" pushover"`, `" sensugo"`, `" slack"`, `" teams"`, `" telegram"`, `" threema"`, `" victorops"`, `" webhook"`, `" wecom"`


### fn withUid

```jsonnet
withUid(value)
```

PARAMETERS:

* **value** (`string`)

UID is the unique identifier of the contact point. The UID can be
set by the user.