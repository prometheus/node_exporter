#!/usr/bin/env bash
#
# Description: Expose metrics from apcaccess
#
# Author: Erwin Oegema <erwin.oegema@gmail.com>


set -ue

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
	local KEY="$1"
	local VAL="$2"
	local TYPE="$3"
	local HELP="$4"

	echo '# HELP' "$KEY" "$HELP"
	echo '# TYPE' "$KEY" "$TYPE"
	echo "$KEY" "$VAL"
}

function try_convert_number {
	local VAL="$1"
	VAL_NMBR=$(echo $VAL | cut -d ' ' -f 1 | awk '{$1=$1};1')
	VAL_TYPE=$(echo $VAL | cut -d ' ' -f 2 | awk '{$1=$1};1')

	case "$VAL_TYPE" in
	# Convert to seconds
	Minutes) VAL_NMBR=$(bc <<< "$VAL_NMBR * 60") ;;
	Hours)   VAL_NMBR=$(bc <<< "$VAL_NMBR * 3600") ;;
	esac

	echo "$VAL_NMBR"
}

function convert_date {
	local VAL="$1"
	# Format date as time since epoch, in seconds
	date -d "$VAL" +%s
}

# Helper props
UPS_MODEL="x"
UPS_NAME=""
UPS_STATUS=""
UPS_SHUTDOWN_CHARGE_MIN=""
UPS_SHUTDOWN_TIME_MAX=""
UPS_SHUTDOWN_TIME_MIN=""
UPS_LAST_TRANSFER_REASON=""
UPS_BATTERY_DATE=""

RESULT=$(/sbin/apcaccess)

while read LINE; do
	VAR_KEY=$(echo $LINE | cut -d : -f 1 | awk '{$1=$1};1')
	VAR_VAL=$(echo $LINE | cut -d : -f 2- | awk '{$1=$1};1')

	case "$VAR_KEY" in
	STARTTIME)  write_line 'node_apc_start_time' $(convert_date "$VAR_VAL") 'untyped' 'The startup time of the APC' ;;
	XONBATT) 	write_line 'node_apc_battery_backup_start' $(convert_date "$VAR_VAL") 'untyped' 'Last time the battery backup started being used (power lost)' ;;
	XOFFBATT) 	write_line 'node_apc_battery_backup_end' $(convert_date "$VAR_VAL") 'untyped' 'Last time the battery backup stopped being used (power recovered)' ;;
	BATTV) 		write_line 'node_apc_battery_voltage' $(try_convert_number "$VAR_VAL") 'gauge' 'Current voltage of the battery inside the APC' ;;
	BCHARGE)	write_line 'node_apc_battery_charge' $(try_convert_number "$VAR_VAL") 'gauge' 'Current charge percentage of the battery inside the APC' ;;
	LINEV) 		write_line 'node_apc_line_voltage' $(try_convert_number "$VAR_VAL") 'gauge' 'Current voltage of the power connected to the APC' ;;
	TIMELEFT)	write_line 'node_apc_time_left_seconds' $(try_convert_number "$VAR_VAL") 'gauge' 'Time left until battery is empty, in seconds' ;;
	LOADPCT)	write_line 'node_apc_load_capacity_percentage' $(try_convert_number "$VAR_VAL") 'gauge' 'Current load capacity on the APC' ;;
	NUMXFERS)	write_line 'node_apc_transfers' $(try_convert_number "$VAR_VAL") 'counter' 'Amount of times the APC had to switch to battery backup' ;;
	
	# Read generic info
	UPSNAME) 	UPS_NAME="$VAR_VAL"; ;;
	MODEL)   	UPS_MODEL="$VAR_VAL" ;;
	STATUS)  	UPS_STATUS="$VAR_VAL" ;;
	LASTXFER)	UPS_LAST_TRANSFER_REASON="$VAR_VAL"; ;;
	BATTDATE)	UPS_BATTERY_DATE="$(convert_date $VAR_VAL)"; ;;
	# *) echo '# Unused prop ' "$VAR_KEY" "$VAR_VAL"
	esac
done <<< "$(echo -e "$RESULT")"

FIRST_KV=1
function write_kv {
	local KEY="$1"
	local VAL="$2"

	# Skip empty values
	if [ "$VAL" == "" ]; then
		return
	fi


	if [ "$FIRST_KV" == "1" ]; then
		FIRST_KV=0
	else
		echo -n ','
	fi

	echo -n "$KEY=\"$VAL\""
}

echo '# HELP node_apc_info APC Info key'
echo -n 'node_apc_ups_info{'

write_kv 'model' "$UPS_MODEL"
write_kv 'name' "$UPS_NAME"
write_kv 'status' "$UPS_STATUS"
write_kv 'battery_date' "$UPS_BATTERY_DATE"
write_kv 'last_transfer_reason' "$UPS_LAST_TRANSFER_REASON"
write_kv 'shutdown_time_min' "$UPS_SHUTDOWN_TIME_MIN"
write_kv 'shutdown_time_max' "$UPS_SHUTDOWN_TIME_MAX"
write_kv 'shutdown_charge_min' "$UPS_SHUTDOWN_CHARGE_MIN"

echo '} 1'
