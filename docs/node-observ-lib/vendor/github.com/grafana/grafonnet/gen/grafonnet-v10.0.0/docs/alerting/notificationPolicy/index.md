# notificationPolicy

grafonnet.alerting.notificationPolicy

## Subpackages

* [matcher](matcher.md)

## Index

* [`fn withContactPoint(value)`](#fn-withcontactpoint)
* [`fn withContinue(value=true)`](#fn-withcontinue)
* [`fn withGroupBy(value)`](#fn-withgroupby)
* [`fn withGroupByMixin(value)`](#fn-withgroupbymixin)
* [`fn withGroupInterval(value)`](#fn-withgroupinterval)
* [`fn withGroupWait(value)`](#fn-withgroupwait)
* [`fn withMatchers(value)`](#fn-withmatchers)
* [`fn withMatchersMixin(value)`](#fn-withmatchersmixin)
* [`fn withMuteTimeIntervals(value)`](#fn-withmutetimeintervals)
* [`fn withMuteTimeIntervalsMixin(value)`](#fn-withmutetimeintervalsmixin)
* [`fn withPolicy(value)`](#fn-withpolicy)
* [`fn withPolicyMixin(value)`](#fn-withpolicymixin)
* [`fn withRepeatInterval(value)`](#fn-withrepeatinterval)

## Fields

### fn withContactPoint

```jsonnet
withContactPoint(value)
```

PARAMETERS:

* **value** (`string`)


### fn withContinue

```jsonnet
withContinue(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


### fn withGroupBy

```jsonnet
withGroupBy(value)
```

PARAMETERS:

* **value** (`array`)


### fn withGroupByMixin

```jsonnet
withGroupByMixin(value)
```

PARAMETERS:

* **value** (`array`)


### fn withGroupInterval

```jsonnet
withGroupInterval(value)
```

PARAMETERS:

* **value** (`string`)


### fn withGroupWait

```jsonnet
withGroupWait(value)
```

PARAMETERS:

* **value** (`string`)


### fn withMatchers

```jsonnet
withMatchers(value)
```

PARAMETERS:

* **value** (`array`)

Matchers is a slice of Matchers that is sortable, implements Stringer, and
provides a Matches method to match a LabelSet against all Matchers in the
slice. Note that some users of Matchers might require it to be sorted.
### fn withMatchersMixin

```jsonnet
withMatchersMixin(value)
```

PARAMETERS:

* **value** (`array`)

Matchers is a slice of Matchers that is sortable, implements Stringer, and
provides a Matches method to match a LabelSet against all Matchers in the
slice. Note that some users of Matchers might require it to be sorted.
### fn withMuteTimeIntervals

```jsonnet
withMuteTimeIntervals(value)
```

PARAMETERS:

* **value** (`array`)


### fn withMuteTimeIntervalsMixin

```jsonnet
withMuteTimeIntervalsMixin(value)
```

PARAMETERS:

* **value** (`array`)


### fn withPolicy

```jsonnet
withPolicy(value)
```

PARAMETERS:

* **value** (`array`)


### fn withPolicyMixin

```jsonnet
withPolicyMixin(value)
```

PARAMETERS:

* **value** (`array`)


### fn withRepeatInterval

```jsonnet
withRepeatInterval(value)
```

PARAMETERS:

* **value** (`string`)

