# items



## Index

* [`fn withTitle(value)`](#fn-withtitle)
* [`fn withType(value)`](#fn-withtype)
* [`fn withValue(value)`](#fn-withvalue)

## Fields

### fn withTitle

```jsonnet
withTitle(value)
```

PARAMETERS:

* **value** (`string`)

Title is an unused property -- it will be removed in the future
### fn withType

```jsonnet
withType(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"dashboard_by_uid"`, `"dashboard_by_id"`, `"dashboard_by_tag"`

Type of the item.
### fn withValue

```jsonnet
withValue(value)
```

PARAMETERS:

* **value** (`string`)

Value depends on type and describes the playlist item.

 - dashboard_by_id: The value is an internal numerical identifier set by Grafana. This
 is not portable as the numerical identifier is non-deterministic between different instances.
 Will be replaced by dashboard_by_uid in the future. (deprecated)
 - dashboard_by_tag: The value is a tag which is set on any number of dashboards. All
 dashboards behind the tag will be added to the playlist.
 - dashboard_by_uid: The value is the dashboard UID