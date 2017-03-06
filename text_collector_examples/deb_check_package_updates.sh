#!/bin/sh
/usr/bin/apt-get update
/bin/echo "" > /var/lib/node_exporter/textfile_collector/debian_updates.prom
/usr/bin/apt-get --just-print upgrade | /bin/grep -- Debian-Security | /bin/grep ^Inst | /usr/bin/wc -l | /bin/sed 's/^/debian_updates_pending{type="security"} /' > /var/lib/node_exporter/textfile_collector/debian_updates.prom
/usr/bin/apt-get --just-print upgrade | /bin/grep -v -- Debian-Security | /bin/grep ^Inst | /usr/bin/wc -l | /bin/sed 's/^/debian_updates_pending{type="others"} /' >> /var/lib/node_exporter/textfile_collector/debian_updates.prom
