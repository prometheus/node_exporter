#!/bin/sh
#
# Return curl duration and HTTP return code. HTTP code is stored in "code" label.
#
# +Usage+
# Input
# - accepts space separated list of targets where a target can be valid URL or rely on curl defaults
#
# Cron Example
# */5 * * * * check_http.sh localhost:8080 google.com | sponge /var/lib/node_exporter/check_http.prom
#
#
# Label code="000" indicates an error running curl and should not be confused as an HTTP return code.
#
# Author: Patrick Freeman <will.pat.free@gmail.com>
echo "# HELP node_curl_time_milliseconds The total time, in seconds, that the full operation lasted. The time will be displayed with millisecond resolution."
echo "# TYPE node_curl_time_milliseconds gauge"
/usr/bin/curl -k -sSf -so /dev/null -w "node_curl_time_milliseconds{code=\"%{http_code}\",url=\"%{url_effective}\"} %{time_total}\n" "$@" 2>/dev/null | grep -e '^#' -e 'node_curl_time_milliseconds'
