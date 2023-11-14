# playlist

grafonnet.playlist

## Subpackages

* [items](items.md)

## Index

* [`fn withInterval(value="5m")`](#fn-withinterval)
* [`fn withItems(value)`](#fn-withitems)
* [`fn withItemsMixin(value)`](#fn-withitemsmixin)
* [`fn withName(value)`](#fn-withname)
* [`fn withUid(value)`](#fn-withuid)

## Fields

### fn withInterval

```jsonnet
withInterval(value="5m")
```

PARAMETERS:

* **value** (`string`)
   - default value: `"5m"`

Interval sets the time between switching views in a playlist.
FIXME: Is this based on a standardized format or what options are available? Can datemath be used?
### fn withItems

```jsonnet
withItems(value)
```

PARAMETERS:

* **value** (`array`)

The ordered list of items that the playlist will iterate over.
FIXME! This should not be optional, but changing it makes the godegen awkward
### fn withItemsMixin

```jsonnet
withItemsMixin(value)
```

PARAMETERS:

* **value** (`array`)

The ordered list of items that the playlist will iterate over.
FIXME! This should not be optional, but changing it makes the godegen awkward
### fn withName

```jsonnet
withName(value)
```

PARAMETERS:

* **value** (`string`)

Name of the playlist.
### fn withUid

```jsonnet
withUid(value)
```

PARAMETERS:

* **value** (`string`)

Unique playlist identifier. Generated on creation, either by the
creator of the playlist of by the application.