#!/usr/bin/env bash
set -eu

# Dependencies: nvme-cli (package)
# Based on code from
# - https://github.com/prometheus/node_exporter/blob/master/text_collector_examples/smartmon.sh
# - https://github.com/vorlon/check_nvme/blob/master/check_nvme.sh
# Author: Henk <henk@wearespindle.com>

output_format_awk="$(
  cat <<'OUTPUTAWK'
BEGIN { v = "" }
v != $1 {
  print "# HELP nvme_" $1 " SMART metric " $1;
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

# Cut device name to /dev/nvmeX
device_list="$(nvme list | awk '/^\/dev/{print $1}' | cut -c1-10)"

for device in ${device_list}; do
  CHECK=$(nvme smart-log "${device}")
  DISK="$(echo "${device}" | cut -c6-10)"

  value_critical_warning=$(echo "$CHECK" | awk '$1 == "critical_warning" {print $3}')
  echo "critical_warning_total{disk=\"${DISK}\"} ${value_critical_warning}"

  value_media_errors=$(echo "$CHECK" | awk '$1 == "media_errors" {print $3}')
  echo "media_errors_total{disk=\"${DISK}\"} ${value_media_errors}"

  value_num_err_log_entries=$(echo "$CHECK" | awk '$1 == "num_err_log_entries" {print $3}')
  echo "num_err_log_entries_total{disk=\"${DISK}\"} ${value_num_err_log_entries}"

  value_power_cycles=$(echo "$CHECK" | awk '$1 == "power_cycles" {print $3}' | sed 's/,\+//g')
  echo "power_cycles_total{disk=\"${DISK}\"} ${value_power_cycles}"

  value_power_on_hours=$(echo "$CHECK" | awk '$1 == "power_on_hours" {print $3}' | sed 's/,\+//g')
  echo "power_on_hours_total{disk=\"${DISK}\"} ${value_power_on_hours}"

  value_temperature=$(echo "$CHECK" | awk '$1 == "temperature" {print $3}')
  echo "temperature_celcius{disk=\"${DISK}\"} ${value_temperature}"

  value_controller_busy_time=$(echo "$CHECK" | awk '$1 == "controller_busy_time" {print $3}' | sed 's/,\+//g')
  echo "controller_busy_time_seconds{disk=\"${DISK}\"} ${value_controller_busy_time}"

  value_available_spare=$(echo "$CHECK" | sed 's/%$//' | awk '$1 == "available_spare" {print $3 / 100}')
  echo "available_spare_ratio{disk=\"${DISK}\"} ${value_available_spare}"

  value_available_spare_threshold=$(echo "$CHECK" | sed 's/%$//' | awk '$1 == "available_spare_threshold" {print $3 / 100}')
  echo "available_spare_threshold_ratio{disk=\"${DISK}\"} ${value_available_spare_threshold}"

  value_percentage_used=$(echo "$CHECK" | sed 's/%$//' | awk '$1 == "percentage_used" {print $3 / 100}')
  echo "percentage_used_ratio{disk=\"${DISK}\"} ${value_percentage_used}"

  value_data_units_written=$(echo "$CHECK" | awk '$1 == "data_units_written" {print $3}' | sed 's/,\+//g')
  echo "data_units_written_total{disk=\"${DISK}\"} ${value_data_units_written}"

  value_data_units_read=$(echo "$CHECK" | awk '$1 == "data_units_read" {print $3}' | sed 's/,\+//g')
  echo "data_units_read_total{disk=\"${DISK}\"} ${value_data_units_read}"

  value_host_read_commands=$(echo "$CHECK" | awk '$1 == "host_read_commands" {print $3}' | sed 's/,\+//g')
  echo "host_read_commands_total{disk=\"${DISK}\"} ${value_host_read_commands}"

  value_host_write_commands=$(echo "$CHECK" | awk '$1 == "host_write_commands" {print $3}' | sed 's/,\+//g')
  echo "host_write_commands_total{disk=\"${DISK}\"} ${value_host_write_commands}"
done | format_output
