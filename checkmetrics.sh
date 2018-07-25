#!/usr/bin/env bash

if [[ ( -z "$1" ) || ( -z "$2" ) ]]; then
    echo "usage: ./checkmetrics.sh /usr/bin/promtool e2e-output.txt"
    exit 1
fi

# Only check node_exporter's metrics, as the Prometheus Go client currently
# exposes a metric with a unit of microseconds.  Once that is fixed, remove
# this filter.
lint=$($1 check metrics < $2 2>&1 | grep "node_")

if [[ ! -z $lint ]]; then
    echo -e "Some Prometheus metrics do not follow best practices:\n"
    echo "$lint"

    exit 1
fi