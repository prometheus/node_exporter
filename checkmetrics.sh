#!/usr/bin/env bash

if [[ ( -z "$1" ) || ( -z "$2" ) ]]; then
    echo "usage: ./checkmetrics.sh /usr/bin/promtool e2e-output.txt"
    exit 1
fi

# Ignore known issues in auto-generated and network specific collectors.
lint=$($1 check metrics < "$2" 2>&1 | grep -v -E "^node_(entropy|memory|netstat|wifi_station)_")

if [[ -n $lint ]]; then
    echo -e "Some Prometheus metrics do not follow best practices:\n"
    echo "$lint"

    exit 1
fi
