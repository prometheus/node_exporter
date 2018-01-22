#!/bin/bash
# Script informed by the collectd monitoring script for smartmontools (using smartctl)
# by Samuel B. <samuel_._behan_(at)_dob_._sk> (c) 2012
# source at: http://devel.dob.sk/collectd-scripts/

# TODO: This probably needs to be a little more complex.  The raw numbers can have more
#       data in them than you'd think.
#       http://arstechnica.com/civis/viewtopic.php?p=22062211

#
# 2017-11-07 extended by Jan Walzer
#  - fixed some minor bugs
#  - extended the list of smart-fields with the ones I found in my DC
#  - implemented 'exact' timestamps for delivery of the metric
#

START_TS="$(TZ=UTC date +%s.%3N)"

parse_smartctl_attributes_awk="$(cat << 'SMARTCTLAWK'
$1 ~ /^ *[0-9]+$/ && $2 ~ /^[a-zA-Z0-9_-]+$/ {
  gsub(/-/, "_");
  printf "%s_value{%s,smart_id=\"%s\"} %d %s\n", $2, labels, $1, $4, ts
  printf "%s_worst{%s,smart_id=\"%s\"} %d %s\n", $2, labels, $1, $5, ts
  printf "%s_threshold{%s,smart_id=\"%s\"} %d %s\n", $2, labels, $1, $6, ts
  printf "%s_raw_value{%s,smart_id=\"%s\"} %e %s\n", $2, labels, $1, $10, ts
}
SMARTCTLAWK
)"

smartmon_attrs="$(cat << 'SMARTMONATTRS'
airflow_temperature_cel
available_reservd_space
ave_block-erase_count
bckgnd_program_page_cnt
command_timeout
crc_error_count
cumulativ_corrected_ecc
current_pending_sector
end-to-end_error
end_to_end_error
erase_fail_count
error_correction_count
g_sense_error_rate
hardware_ecc_recovered
host_program_page_count
host_reads_32mib
host_reads_mib
host_writes_32mib
host_writes_mib
load_cycle_count
media_wearout_indicator
nand_writes_1gib
nand_writes_32mib
offline_uncorrectable
percent_lifetime_used
power-off_retract_count
power_cycle_count
power_loss_cap_test
power_on_hours
program_fail_cnt_total
program_fail_count
program_fail_count_chip
raw_read_error_rate
reallocate_nand_blk_cnt
reallocated_event_count
reallocated_sector_ct
reported_uncorrect
reserved_block_count
runtime_bad_block
sata_downshift_count
sata_interfac_downshift
seek_error_rate
seek_time_performance
spin_retry_count
spin_up_time
start_stop_count
success_rain_recov_cnt
temperature_case
temperature_celsius
temperature_internal
thermal_throttle
throughput_performance
total_host_sector_write
total_lbas_read
total_lbas_written
udma_crc_error_count
unexpect_power_loss_ct
unknown_attribute
unknown_ssd_attribute
unsafe_shutdown_count
unused_reserve_nand_blk
unused_rsvd_blk_cnt_tot
wear_leveling_count
workld_host_reads_perc
workld_media_wear_indic
workload_minutes
write_error_rate
SMARTMONATTRS
)"
# shellcheck disable=SC2086
smartmon_attrs="$(echo ${smartmon_attrs} | xargs | tr ' ' '|')"

parse_smartctl_attributes() {
  local disk="$1"
  local disk_type="$2"
  local labels="disk=\"${disk}\",type=\"${disk_type}\""
# shellcheck disable=SC2018
# shellcheck disable=SC2019
  sed 's/^ \+//g' \
    | awk -v labels="${labels}" -v ts="$TS" "${parse_smartctl_attributes_awk}" 2>/dev/null \
    | tr A-Z a-z \
    | grep -E "(${smartmon_attrs})"
}

parse_smartctl_info() {
  local -i smart_available=0 smart_enabled=0 smart_healthy=0
  local disk="$1" disk_type="$2"
# shellcheck disable=SC2162
  while read line ; do
    info_type="$(echo "${line}" | cut -f1 -d: | tr ' ' '_')"
    info_value="$(echo "${line}" | cut -f2- -d: | sed 's/^ \+//g')"
    case "${info_type}" in
      Model_Family) model_family="${info_value}" ;;
      Device_Model) device_model="${info_value}" ;;
      Serial_Number) serial_number="${info_value}" ;;
      Firmware_Version) fw_version="${info_value}" ;;
      Vendor) vendor="${info_value}" ;;
      Product) product="${info_value}" ;;
      Revision) revision="${info_value}" ;;
      Logical_Unit_id) lun_id="${info_value}" ;;
    esac
    if [[ "${info_type}" == 'SMART_support_is' ]] ; then
      case "${info_value:0:7}" in
        Enabled) smart_enabled=1 ;;
        Availab) smart_available=1 ;;
        Unavail) smart_available=0 ;;
      esac
    fi
    if [[ "${info_type}" == 'SMART_overall-health_self-assessment_test_result' ]] ; then
      case "${info_value:0:6}" in
        PASSED) smart_healthy=1 ;;
      esac
    elif [[ "${info_type}" == 'SMART_Health_Status' ]] ; then
      case "${info_value:0:2}" in
        OK) smart_healthy=1 ;;
      esac
    fi
  done
  if [[ -n "${vendor}" ]] ; then
    echo "device_info{disk=\"${disk}\",type=\"${disk_type}\",vendor=\"${vendor}\",product=\"${product}\",revision=\"${revision}\",lun_id=\"${lun_id}\"} 1 $TS"
  else
    echo "device_info{disk=\"${disk}\",type=\"${disk_type}\",model_family=\"${model_family}\",device_model=\"${device_model}\",serial_number=\"${serial_number}\",firmware_version=\"${fw_version}\"} 1 $TS"
  fi
  echo "device_smart_available{disk=\"${disk}\",type=\"${disk_type}\"} ${smart_available} $TS"
  echo "device_smart_enabled{disk=\"${disk}\",type=\"${disk_type}\"} ${smart_enabled} $TS"
  echo "device_smart_healthy{disk=\"${disk}\",type=\"${disk_type}\"} ${smart_healthy} $TS"
}

output_format_awk="$(cat << 'OUTPUTAWK'
BEGIN { v = "" }
v != $1 {
  print "# HELP smartmon_" $1 " SMART metric " $1;
  print "# TYPE smartmon_" $1 " gauge";
  v = $1
}
{print "smartmon_" $0}
OUTPUTAWK
)"

format_output() {
  sort \
  | awk -F'{' "${output_format_awk}"
}

TS="$(TZ=UTC date +%s%3N)"
smartctl_version="$(/usr/sbin/smartctl -V | head -n1  | awk '$1 == "smartctl" {print $2}')"

echo "smartctl_version{version=\"${smartctl_version}\"} 1 $TS" | format_output

if [[ "$(expr "${smartctl_version}" : '\([0-9]*\)\..*')" -lt 6 ]] ; then
  exit
fi

device_list="$(/usr/sbin/smartctl --scan-open | awk '/^\/dev/{print $1 "|" $3}')"

for device in ${device_list}; do
# shellcheck disable=SC2086
  disk="$(echo ${device} | cut -f1 -d'|')"
# shellcheck disable=SC2086
  type="$(echo ${device} | cut -f2 -d'|')"
# shellcheck disable=SC2030
  TS="$(TZ=UTC date +%s.%N)"
  echo "smartctl_run{disk=\"${disk}\",type=\"${type}\"} $TS"
  # Get the SMART information and health
  TS="$(TZ=UTC date +%s%3N)"
  /usr/sbin/smartctl -i -H -d "${type}" "${disk}" | parse_smartctl_info "${disk}" "${type}"
  # Get the SMART attributes
  TS="$(TZ=UTC date +%s%3N)"
  /usr/sbin/smartctl -A -d "${type}" "${disk}" | parse_smartctl_attributes "${disk}" "${type}"
done | format_output

DURATION="$(TZ=UTC date -u -d "-${START_TS}seconds" +%s.%3N)"

# shellcheck disable=SC2031
echo "smartctl_run_duration{} $DURATION $TS" | format_output
