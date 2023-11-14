---
permalink: /url/
---

# url

```jsonnet
local url = import "github.com/jsonnet-libs/xtd/url.libsonnet"
```

`url` provides functions to deal with URLs

## Index

* [`fn encodeQuery(params)`](#fn-encodequery)
* [`fn escapeString(str, excludedChars=[])`](#fn-escapestring)
* [`fn join(splitObj)`](#fn-join)
* [`fn parse(url)`](#fn-parse)

## Fields

### fn encodeQuery

```ts
encodeQuery(params)
```

`encodeQuery` takes an object of query parameters and returns them as an escaped `key=value` string

### fn escapeString

```ts
escapeString(str, excludedChars=[])
```

`escapeString` escapes the given string so it can be safely placed inside an URL, replacing special characters with `%XX` sequences

### fn join

```ts
join(splitObj)
```

`join` joins URLs from the object generated from `parse`

### fn parse

```ts
parse(url)
```

`parse` parses absolute and relative URLs.

<scheme>://<netloc>/<path>;parameters?<query>#<fragment>

Inspired by Python's urllib.urlparse, following several RFC specifications.
