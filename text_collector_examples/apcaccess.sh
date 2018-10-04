#!/usr/bin/env bash
#
# Description: Expose metrics from apcaccess
#
# Author: Erwin Oegema <erwin.oegema@gmail.com>
# Required tools:
# - apcaccess
# - bc
# - awk
# - cut

# Make sure we also terminate on pipe operations
set -o pipefail

# It is recommended to check https://linux.die.net/man/8/apcaccess for
# more properties that could be exposed
# Example response:
#
# APC      : 001,000,0000
# DATE     : 2018-10-03 15:36:29 +0200
# HOSTNAME : hostname
# VERSION  : 3.14.14 (31 May 2016) debian
# UPSNAME  : backups700
# CABLE    : USB Cable
# DRIVER   : USB UPS Driver
# UPSMODE  : Stand Alone
# STARTTIME: 2018-09-27 18:28:59 +0200
# MODEL    : Back-UPS ES 700G
# STATUS   : ONLINE
# LINEV    : 234.0 Volts
# LOADPCT  : 10.0 Percent
# BCHARGE  : 100.0 Percent
# TIMELEFT : 27.6 Minutes
# MBATTCHG : 5 Percent
# MINTIMEL : 3 Minutes
# MAXTIME  : 0 Seconds
# SENSE    : Medium
# LOTRANS  : 180.0 Volts
# HITRANS  : 266.0 Volts
# ALARMDEL : 30 Seconds
# BATTV    : 13.6 Volts
# LASTXFER : Unacceptable line voltage changes
# NUMXFERS : 1
# XONBATT  : 2018-10-01 14:20:44 +0200
# TONBATT  : 0 Seconds
# CUMONBATT: 10 Seconds
# XOFFBATT : 2018-10-01 14:20:54 +0200
# STATFLAG : 0x05000008
# SERIALNO : xxx
# BATTDATE : 2018-04-21
# NOMINV   : 230 Volts
# NOMBATTV : 12.0 Volts
# FIRMWARE : 871.O4 .I USB FW:O4
# END APC  : 2018-10-03 15:37:04 +0200

function write_line {
	local key="$1"
	local val="$2"
	local type="$3"
	local help="$4"

	echo '# HELP' "$key" "$help"
	echo '# TYPE' "$key" "$type"
	echo "$key" "$val" "$data_date"
}

function try_convert_number {
	local val="$1"
	val_nmbr=$(echo $val | cut -d ' ' -f 1 | awk '{$1=$1};1')
	val_type=$(echo $val | cut -d ' ' -f 2 | awk '{$1=$1};1')

	case "$val_type" in
	# Convert to seconds
	Minutes) val_nmbr=$(bc <<< "$val_nmbr * 60") ;;
	Hours)   val_nmbr=$(bc <<< "$val_nmbr * 3600") ;;
	esac

	echo "$val_nmbr"
}

function convert_date {
	local val="$1"
	# Format date as time since epoch, in seconds
	date -d "$val" +%s
}

# Helper props
ups_model="x"
ups_name=""
ups_status=""
ups_shutdown_charge_min=""
ups_shutdown_time_max=""
ups_shutdown_time_min=""
ups_last_transfer_reason=""
ups_battery_date=""
data_date=""

result=$(/sbin/apcaccess)
error_code=$?
if [[ $? -ne 0 ]]; then
	(>&2 echo "apcaccess exited with error code $error_code!")
	exit $error_code
fi

while read line; do
	var_key=$(echo $line | cut -d : -f 1  | awk '{$1=$1};1')
	var_val=$(echo $line | cut -d : -f 2- | awk '{$1=$1};1')

	case "$var_key" in
	STARTTIME)	write_line 'node_apc_start_time' $(convert_date "$var_val") 'gauge' 'The startup time of the APC' ;;
	XONBATT)	write_line 'node_apc_battery_backup_start' $(convert_date "$var_val") 'gauge' 'Last time the battery backup started being used (power lost)' ;;
	XOFFBATT)	write_line 'node_apc_battery_backup_end' $(convert_date "$var_val") 'gauge' 'Last time the battery backup stopped being used (power recovered)' ;;
	BATTV)		write_line 'node_apc_battery_voltage' $(try_convert_number "$var_val") 'gauge' 'Current voltage of the battery inside the APC' ;;
	BCHARGE)	write_line 'node_apc_battery_charge' $(try_convert_number "$var_val") 'gauge' 'Current charge percentage of the battery inside the APC' ;;
	LINEV)		write_line 'node_apc_line_voltage' $(try_convert_number "$var_val") 'gauge' 'Current voltage of the power connected to the APC' ;;
	TIMELEFT)	write_line 'node_apc_time_left_seconds' $(try_convert_number "$var_val") 'gauge' 'Time left until battery is empty, in seconds' ;;
	LOADPCT)	write_line 'node_apc_load_capacity_percentage' $(try_convert_number "$var_val") 'gauge' 'Current load capacity on the APC' ;;
	NUMXFERS)	write_line 'node_apc_transfers' $(try_convert_number "$var_val") 'counter' 'Amount of times the APC had to switch to battery backup' ;;

	# Read generic info
	UPSNAME)	ups_name="$var_val" ;;
	MODEL)		ups_model="$var_val" ;;
	STATUS)		ups_status="$var_val" ;;
	LASTXFER)	ups_last_transfer_reason="$var_val" ;;
	BATTDATE)	ups_battery_date="$(convert_date $var_val)" ;;
	# This row contains the date the data was received from the APC.
	# We'll be using it to tag our records with the current date
	DATE)
		data_date="$(convert_date "$var_val")"
		# Lets add some milliseconds
		data_date="${data_date}000"
	;;
	# *) echo '# Unused prop ' "$var_key" "$var_val"
	esac
done <<< "$result"

first_kv=1
function write_kv {
	local key="$1"
	local val="$2"

	# Skip empty values
	if [ "$val" == "" ]; then
		return
	fi


	if [ "$first_kv" == "1" ]; then
		first_kv=0
	else
		echo -n ','
	fi

	echo -n "$key=\"$val\""
}

echo '# HELP node_apc_info APC Info key'
echo -n 'node_apc_ups_info{'

write_kv 'model' "$ups_model"
write_kv 'name' "$ups_name"
write_kv 'status' "$ups_status"
write_kv 'battery_date' "$ups_battery_date"
write_kv 'last_transfer_reason' "$ups_last_transfer_reason"
write_kv 'shutdown_time_min' "$ups_shutdown_time_min"
write_kv 'shutdown_time_max' "$ups_shutdown_time_max"
write_kv 'shutdown_charge_min' "$ups_shutdown_charge_min"

echo "} 1 $data_date"
