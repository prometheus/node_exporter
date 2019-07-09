#!/bin/sh
#
# Custom logic to capture any specific process and monitor the status of same
# It takes input the way you would like to capture the process and what name you would like the process to have 
# Usage: add this to crontab:
# argument 1 : the way you want to capture the process , it could be a custom one lines bash script as example below in between () these brackets
# argument 2 : the process name you would like to add to that captured process
# example I am pullig splunk main process info by using a custom logic in ()
# */5 * * * * sh processmon.sh '(grep splunkd | grep process-runner)' splunk | sponge /var/lib/node_exporter/yourprocess_info.prom
#
#
# Author: Yadu Mathur <mathur.yadu@gmail.com>


#!/bin/bash
IFS=$'\n'
#capture your one liner logic for specific process
processname=$2
cli=`echo $1 | tr -d "(" | tr -d ")"`
echo "ps -ef | $cli | grep -v grep | grep -v processmon | cut -d' ' -f 2" > .mycommand.sh
pid=`sh .mycommand.sh`

if [ "$pid" != "" ]
then
 status=1
  echo "node_"$processname"_status{process_name=\"${processname}\"} 1"
else
 status=0
  echo "node_"$processname"_status{process_name=\"${processname}\"} 0"
fi

## You can add multiple piece of this stub of code to monitor multiple process
