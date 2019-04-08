#!/usr/bin/env bash
set -eu

# Dependencies: nvme-cli, jq (packages)
# Based on code from
# - https://github.com/prometheus/node_exporter/blob/master/text_collector_examples/smartmon.sh
# - https://github.com/prometheus/node_exporter/blob/master/text_collector_examples/mellanox_hca_temp
# - https://github.com/vorlon/check_nvme/blob/master/check_nvme.sh
#
# Author: Henk <henk@wearespindle.com>

# Check if we are root
if [ "$EUID" -ne 0 ]; then
  echo "${0##*/}: Please run as root!" >&2
  exit 1
fi

# Check if programs are installed
if ! command -v nvme >/dev/null 2>&1; then
  echo "${0##*/}: nvme is not installed. Aborting." >&2
  exit 1
fi

output_format_awk="$(
  cat <<'OUTPUTAWK'
BEGIN { v = "" }
v != $1 {
  print "# HELP nvme_" $1 " SMART metric " $1;
  if ($1 ~ /_total$/)
    print "# TYPE nvme_" $1 " counter";
  else
    print "# TYPE nvme_" $1 " gauge";
  v = $1
}
{print "nvme_" $0}
OUTPUTAWK
)"

format_output() {
  sort | awk -F'{' "${output_format_awk}"
}

# Get the nvme-cli version
nvme_version="$(nvme version | awk '$1 == "nvme" {print $3}')"
echo "nvmecli{version=\"${nvme_version}\"} 1" | format_output

# Get devices
device_list="$(nvme list | awk '/^\/dev/{print $1}')"

# Loop through the NVMe devices
for device in ${device_list}; do
  json_check="$(nvme smart-log -o json "${device}")"
  disk="$(echo "${device}" | cut -c6-10)"

  # The temperature value in JSON is in Kelvin, we want Celsius
  value_temperature="$(echo "$json_check" | jq '.temperature - 273')"
  echo "temperature_celcius{device=\"${disk}\"} ${value_temperature}"

  value_available_spare="$(echo "$json_check" | jq '.avail_spare / 100')"
  echo "available_spare_ratio{device=\"${disk}\"} ${value_available_spare}"

  value_available_spare_threshold="$(echo "$json_check" | jq '.spare_thresh / 100')"
  echo "available_spare_threshold_ratio{device=\"${disk}\"} ${value_available_spare_threshold}"

  value_percentage_used="$(echo "$json_check" | jq '.percent_used / 100')"
  echo "percentage_used_ratio{device=\"${disk}\"} ${value_percentage_used}"

  value_critical_warning="$(echo "$json_check" | jq '.critical_warning')"
  echo "critical_warning_total{device=\"${disk}\"} ${value_critical_warning}"

  value_media_errors="$(echo "$json_check" | jq '.media_errors')"
  echo "media_errors_total{device=\"${disk}\"} ${value_media_errors}"

  value_num_err_log_entries="$(echo "$json_check" | jq '.num_err_log_entries')"
  echo "num_err_log_entries_total{device=\"${disk}\"} ${value_num_err_log_entries}"

  value_power_cycles="$(echo "$json_check" | jq '.power_cycles')"
  echo "power_cycles_total{device=\"${disk}\"} ${value_power_cycles}"

  value_power_on_hours="$(echo "$json_check" | jq '.power_on_hours')"
  echo "power_on_hours_total{device=\"${disk}\"} ${value_power_on_hours}"

  value_controller_busy_time="$(echo "$json_check" | jq '.controller_busy_time')"
  echo "controller_busy_time_seconds{device=\"${disk}\"} ${value_controller_busy_time}"

  value_data_units_written="$(echo "$json_check" | jq '.data_units_written')"
  echo "data_units_written_total{device=\"${disk}\"} ${value_data_units_written}"

  value_data_units_read="$(echo "$json_check" | jq '.data_units_read')"
  echo "data_units_read_total{device=\"${disk}\"} ${value_data_units_read}"

  value_host_read_commands="$(echo "$json_check" | jq '.host_read_commands')"
  echo "host_read_commands_total{device=\"${disk}\"} ${value_host_read_commands}"

  value_host_write_commands="$(echo "$json_check" | jq '.host_write_commands')"
  echo "host_write_commands_total{device=\"${disk}\"} ${value_host_write_commands}"
done | format_output
