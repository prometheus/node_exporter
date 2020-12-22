#!/bin/bash
#
#
# Description: Expose metrics for number of boots of the system
# Dependencies: none
#
# The script creates one metric series which is the number of recorded
# boots in the /var/log/wtmp file. If wtmp is rotated, the next time the
# node exporter is restarted it will recalculate the value, and may over
# count restarts by one.
#
# Author: Clayton Coleman <smarterclayton@gmail.com>

set -euo pipefail

v=$(mktemp --suffix=.boots.working)
chmod uga+r "${v}"

reboots=$( last reboot --time-format=iso -R | grep -cE '^reboot' )
if [[ "${reboots}" -gt 0 ]]; then
  echo '# HELP node_boots_total reports a single series which is the number of times this system has been booted excluding the current boot. If the value is zero, this is the first time the system has booted. The value is always non-negative.' >>"${v}"
  echo '# TYPE node_boots_total counter' >>"${v}"
  echo "node_boots_total{} $(( reboots - 1 ))" >>"${v}"
fi

mv "${v}" boots.prom
