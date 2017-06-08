#!/bin/sh
/usr/bin/apt-get update
/usr/bin/apt-get --just-print upgrade | /usr/bin/awk '/^Inst/ {print $5}' | /usr/bin/sort  | /usr/bin/uniq -c | /usr/bin/awk '{ gsub("\\\\", "\\\\", $2); gsub("\"", "\\\"", $2); print "apt_upgrades_pending{origin=\"" $2 "\"} " $1 }' > /var/lib/node_exporter/textfile_collector/debian_updates.prom.$$
/bin/mv /var/lib/node_exporter/textfile_collector/debian_updates.prom.$$ /var/lib/node_exporter/textfile_collector/debian_updates.prom
