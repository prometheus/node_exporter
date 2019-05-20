#!/usr/bin/env bash
# Note: This script uses "mdadm --detail" to get some of the metrics, so it must be run as root.
#       It is designed to be run periodically in a cronjob, and output to /var/lib/node_exporter/textfile_collector/md_info_detail.prom
#       $ cat /etc/cron.d/prometheus_md_info_detail
#       * * * * * bash /var/lib/node_exporter/md_info_detail.sh > /var/lib/node_exporter/md_info_detail.prom.$$ && mv /var/lib/node_exporter/md_info_detail.prom.$$ /var/lib/node_exporter/md_info_detail.prom

set -eu

for MD_DEVICE in /dev/md/*; do
  # Subshell to avoid eval'd variables from leaking between iterations
  (
    # Resolve symlink to discover device, e.g. /dev/md127
    MD_DEVICE_NUM=$(readlink -f "${MD_DEVICE}")

    # Remove /dev/ prefix
    MD_DEVICE_NUM=${MD_DEVICE_NUM#/dev/}
    MD_DEVICE=${MD_DEVICE#/dev/md/}

    # Query sysfs for info about md device
    SYSFS_BASE="/sys/devices/virtual/block/${MD_DEVICE_NUM}/md"
    MD_LAYOUT=$(cat "${SYSFS_BASE}/layout")
    MD_LEVEL=$(cat "${SYSFS_BASE}/level")
    MD_METADATA_VERSION=$(cat "${SYSFS_BASE}/metadata_version")
    MD_NUM_RAID_DISKS=$(cat "${SYSFS_BASE}/raid_disks")

    # Remove 'raid' prefix from RAID level
    MD_LEVEL=${MD_LEVEL#raid}

    # Output disk metrics
    for RAID_DISK in ${SYSFS_BASE}/rd[0-9]*; do
      DISK=$(readlink -f "${RAID_DISK}/block")
      DISK_DEVICE=$(basename "${DISK}")
      RAID_DISK_DEVICE=$(basename "${RAID_DISK}")
      RAID_DISK_INDEX=${RAID_DISK_DEVICE#rd}
      RAID_DISK_STATE=$(cat "${RAID_DISK}/state")

      DISK_SET=""
      # Determine disk set using logic from mdadm: https://github.com/neilbrown/mdadm/commit/2c096ebe4b
      if [[ ${RAID_DISK_STATE} == "in_sync" && ${MD_LEVEL} == 10 && $((MD_LAYOUT & ~0x1ffff)) ]]; then
        NEAR_COPIES=$((MD_LAYOUT & 0xff))
        FAR_COPIES=$(((MD_LAYOUT >> 8) & 0xff))
        COPIES=$((NEAR_COPIES * FAR_COPIES))

        if [[ $((MD_NUM_RAID_DISKS % COPIES == 0)) && $((COPIES <= 26)) ]]; then
          DISK_SET=$((RAID_DISK_INDEX % COPIES))
        fi
      fi

      echo -n "node_md_disk_info{disk_device=\"${DISK_DEVICE}\", md_device=\"${MD_DEVICE_NUM}\""
      if [[ -n ${DISK_SET} ]]; then
        SET_LETTERS=({A..Z})
        echo -n ", md_set=\"${SET_LETTERS[${DISK_SET}]}\""
      fi
      echo "} 1"
    done

    # Get output from mdadm --detail (Note: root/sudo required)
    MDADM_DETAIL_OUTPUT=$(mdadm --detail /dev/"${MD_DEVICE_NUM}")

    # Output RAID "Devices", "Size" and "Event" metrics, from the output of "mdadm --detail"
    while IFS= read -r line ; do
      # Filter out these keys that have numeric values that increment up
      if echo "$line" | grep -E -q "Devices :|Array Size :| Used Dev Size :|Events :"; then
        MDADM_DETAIL_KEY=$(echo "$line" | cut -d ":" -f 1 | tr -cd '[a-zA-Z0-9]._-')
        MDADM_DETAIL_VALUE=$(echo "$line" | cut -d ":" -f 2 | cut -d " " -f 2 | sed 's:^ ::')
        echo "node_md_info_${MDADM_DETAIL_KEY}{md_device=\"${MD_DEVICE_NUM}\", md_name=\"${MD_DEVICE}\", raid_level=\"${MD_LEVEL}\", md_num_raid_disks=\"${MD_NUM_RAID_DISKS}\", md_metadata_version=\"${MD_METADATA_VERSION}\"} ${MDADM_DETAIL_VALUE}"
      fi
    done  <<< "$MDADM_DETAIL_OUTPUT"

    # Output RAID detail metrics info from the output of "mdadm --detail"
    # NOTE: Sending this info as labels rather than separate metrics, because some of them can be strings.
    echo -n "node_md_info{md_device=\"${MD_DEVICE_NUM}\", md_name=\"${MD_DEVICE}\", raid_level=\"${MD_LEVEL}\", md_num_raid_disks=\"${MD_NUM_RAID_DISKS}\", md_metadata_version=\"${MD_METADATA_VERSION}\""
    while IFS= read -r line ; do
      # Filter for lines with a ":", to use for Key/Value pairs in labels
      if echo "$line" | grep -E -q ":" ; then
        # Exclude lines with these keys, as they're values are numbers that increment up and captured in individual metrics above
        if echo "$line" | grep -E -qv "Array Size|Used Dev Size|Events|Update Time" ; then
          echo -n ", "
          MDADM_DETAIL_KEY=$(echo "$line" | cut -d ":" -f 1 | tr -cd '[a-zA-Z0-9]._-')
          MDADM_DETAIL_VALUE=$(echo "$line" | cut -d ":" -f 2- | sed 's:^ ::')
          echo -n "${MDADM_DETAIL_KEY}=\"${MDADM_DETAIL_VALUE}\""
        fi
      fi
    done  <<< "$MDADM_DETAIL_OUTPUT"
    echo "} 1"
  )
done
