#!/usr/bin/env bash
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

    # Output RAID array metrics
    # NOTE: Metadata version is a label rather than a separate metric because the version can be a string
    echo "node_md_info{md_device=\"${MD_DEVICE_NUM}\", md_name=\"${MD_DEVICE}\", raid_level=\"${MD_LEVEL}\", md_metadata_version=\"${MD_METADATA_VERSION}\"} 1"
  )
done
