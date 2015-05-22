## 0.9.0 / unreleased

* [BUGFIX] Fix `/proc/net/dev` parsing.
* [CLEANUP] Remove the `attributes` collector, use `textfile` instead.
* [CLEANUP] Replace last uses of the configuration file with flags.
* [IMPROVEMENT] Remove cgo dependency.
* [IMPROVEMENT] Sort collector names when printing.
* [FEATURE] IPVS stats collector.

## 0.8.1 / 2015-05-17

* [MAINTENANCE] Use the common Prometheus build infrastructure.
* [MAINTENANCE] Update former Google Code imports.
* [IMPROVEMENT] Log the version at startup.
* [FEATURE] TCP stats collector

## 0.8.0 / 2015-03-09
* [CLEANUP] Introduced semantic versioning and changelog. From now on,
  changes will be reported in this file.
