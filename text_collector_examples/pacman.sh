#!/bin/bash
#
#
# Description: Expose metrics from pacman updates
# If installed The bash script *checkupdates*, included with the
# *pacman-contrib* package, is used to calculate the number of pending updates.
# Otherwise *pacman* is used for calculation.
#
# Author: Sven Haardiek <sven@haardiek.de>

set -o errexit
set -o nounset
set -o pipefail

if [ -x /usr/bin/checkupdates ]
then
    updates=$(/usr/bin/checkupdates | wc -l)
    cache=0
else
    if ! updates=$(/usr/bin/pacman -Qu | wc -l)
    then
        updates=0
    fi
    cache=1
fi

echo "# HELP updates_pending number of pending updates from pacman"
echo "# TYPE updates_pending gauge"
echo "pacman_updates_pending $updates"

echo "# HELP pacman_updates_pending_from_cache pending updates information are from cache"
echo "# TYPE pacman_updates_pending_from_cache gauge"
echo "pacman_updates_pending_from_cache $cache"
