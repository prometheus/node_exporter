# Version 0.16.0 Upgrade Guide

The `node_exporter` 0.16.0 and newer renamed many metrics in order to conform with Prometheus [naming best practices].

In order to allow easy upgrades, there are several options.

## Update dashboards

Grafana users can add multiple queries in order to display both the old and new data simultaneously.

## Use recording rules

We have provided a [sample recording rule set that translates old metrics to new ones] and the [one that translates new metrics format to old one] to create duplicate metrics (it translates "old" metrics format to new one).  This has a minor disadvantage that it creates a lot of extra data, and re-aligns the timestamps of the data.

## Run both old and new versions simultaneously.

It's possible to run both the old and new exporter on different ports, and include an additional scrape job in Prometheus.  It's recommended to enable only the collectors that have name changes that you care about.

[naming best practices]: https://prometheus.io/docs/practices/naming/
[sample recording rule set that translates old metrics to new ones]: example-16-compatibility-rules.yml
[one that translates new metrics format to old one]: example-16-compatibility-rules-new-to-old.yml
