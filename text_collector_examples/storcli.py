#!/usr/bin/env python3
"""
Script to parse StorCLI's JSON output and expose
MegaRAID health as Prometheus metrics.

Tested against StorCLI 'Ver 1.14.12 Nov 25, 2014'.

StorCLI reference manual:
http://docs.avagotech.com/docs/12352476

Advanced Software Options (ASO) not exposed as metrics currently.

JSON key abbreviations used by StorCLI are documented in the standard command
output, i.e.  when you omit the trailing 'J' from the command.

Formatting done with YAPF:
$ yapf -i --style '{COLUMN_LIMIT: 99}' storcli.py
"""

from __future__ import print_function
from datetime import datetime
import argparse
import collections
import json
import os
import shlex
import subprocess

DESCRIPTION = """Parses StorCLI's JSON output and exposes MegaRAID health as
    Prometheus metrics."""
VERSION = '0.0.3'

storcli_path = ''
metric_prefix = 'megaraid_'
metric_list = {}
metric_list = collections.defaultdict(list)


def main(args):
    """ main """
    global storcli_path
    storcli_path = args.storcli_path
    data = get_storcli_json('/cALL show all J')

    try:
        # All the information is collected underneath the Controllers key
        data = data['Controllers']

        for controller in data:
            response = controller['Response Data']
            if response['Version']['Driver Name'] == 'megaraid_sas':
                handle_megaraid_controller(response)
            elif response['Version']['Driver Name'] == 'mpt3sas':
                handle_sas_controller(response)
    except KeyError:
        pass

    print_all_metrics(metric_list)


def handle_sas_controller(response):
    (controller_index, baselabel) = get_basic_controller_info(response)
    add_metric('healthy', baselabel, int(response['Status']['Controller Status'] == 'OK'))
    add_metric('ports', baselabel, response['HwCfg']['Backend Port Count'])
    try:
        # The number of physical disks is half of the number of items in this dict
        # Every disk is listed twice - once for basic info, again for detailed info
        add_metric('physical_drives', baselabel,
                   len(response['Physical Device Information'].keys()) / 2)
    except AttributeError:
        pass
    # Split up string to not trigger CodeSpell issues
    add_metric('temperature', baselabel,
               int(response['HwCfg']['ROC temperature(Degree Celc' + 'ius)']))
    for key, basic_disk_info in response['Physical Device Information'].items():
        if 'Detailed Information' in key:
            continue
        create_metrcis_of_physical_drive(basic_disk_info[0],
                                         response['Physical Device Information'], controller_index)


def handle_megaraid_controller(response):
    (controller_index, baselabel) = get_basic_controller_info(response)

    # BBU Status Optimal value is 0 for cachevault and 32 for BBU
    add_metric('battery_backup_healthy', baselabel,
               int(response['Status']['BBU Status'] in [0, 32]))
    add_metric('degraded', baselabel, int(response['Status']['Controller Status'] == 'Degraded'))
    add_metric('failed', baselabel, int(response['Status']['Controller Status'] == 'Failed'))
    add_metric('healthy', baselabel, int(response['Status']['Controller Status'] == 'Optimal'))
    add_metric('drive_groups', baselabel, response['Drive Groups'])
    add_metric('virtual_drives', baselabel, response['Virtual Drives'])
    add_metric('physical_drives', baselabel, response['Physical Drives'])
    add_metric('ports', baselabel, response['HwCfg']['Backend Port Count'])
    add_metric('scheduled_patrol_read', baselabel,
               int('hrs' in response['Scheduled Tasks']['Patrol Read Reoccurrence']))
    add_metric('temperature', baselabel, int(response['HwCfg']['ROC temperature(Degree Celsius)']))
    for cvidx, cvinfo in enumerate(response['Cachevault_Info']):
        add_metric('cv_temperature', baselabel + ',cvidx="' + str(cvidx) + '"', int(cvinfo['Temp'].replace('C','')))

    time_difference_seconds = -1
    system_time = datetime.strptime(response['Basics'].get('Current System Date/time'),
                                    "%m/%d/%Y, %H:%M:%S")
    controller_time = datetime.strptime(response['Basics'].get('Current Controller Date/Time'),
                                        "%m/%d/%Y, %H:%M:%S")
    if system_time and controller_time:
        time_difference_seconds = abs(system_time - controller_time).seconds
        add_metric('time_difference', baselabel, time_difference_seconds)

    for virtual_drive in response['VD LIST']:
        vd_position = virtual_drive.get('DG/VD')
        drive_group, volume_group = -1, -1
        if vd_position:
            drive_group = vd_position.split('/')[0]
            volume_group = vd_position.split('/')[1]
        vd_baselabel = 'controller="{0}",DG="{1}",VG="{2}"'.format(controller_index, drive_group,
                                                                volume_group)
        vd_info_label = vd_baselabel + ',name="{0}",cache="{1}",type="{2}",state="{3}"'.format(
            str(virtual_drive.get('Name')).strip(),
            str(virtual_drive.get('Cache')).strip(),
            str(virtual_drive.get('TYPE')).strip(),
            str(virtual_drive.get('State')).strip())
        add_metric('vd_info', vd_info_label, 1)

    if response['Physical Drives'] > 0:
        data = get_storcli_json('/cALL/eALL/sALL show all J')
        drive_info = data['Controllers'][controller_index]['Response Data']
    for physical_drive in response['PD LIST']:
        create_metrcis_of_physical_drive(physical_drive, drive_info, controller_index)


def get_basic_controller_info(response):
    controller_index = response['Basics']['Controller']
    baselabel = 'controller="{0}"'.format(controller_index)

    controller_info_label = baselabel + ',model="{0}",serial="{1}",fwversion="{2}"'.format(
        str(response['Basics']['Model']).strip(),
        str(response['Basics']['Serial Number']).strip(),
        str(response['Version']['Firmware Version']).strip(),
    )
    add_metric('controller_info', controller_info_label, 1)

    return (controller_index, baselabel)


def create_metrcis_of_physical_drive(physical_drive, detailed_info_array, controller_index):
    enclosure = physical_drive.get('EID:Slt').split(':')[0]
    slot = physical_drive.get('EID:Slt').split(':')[1]

    pd_baselabel = 'controller="{0}",enclosure="{1}",slot="{2}"'.format(controller_index, enclosure,
                                                                     slot)
    pd_info_label = pd_baselabel + \
        ',disk_id="{0}",interface="{1}",media="{2}",model="{3}",DG="{4}",state="{5}"'.format(
            str(physical_drive.get('DID')).strip(),
            str(physical_drive.get('Intf')).strip(),
            str(physical_drive.get('Med')).strip(),
            str(physical_drive.get('Model')).strip(),
            str(physical_drive.get('DG')).strip(),
            str(physical_drive.get('State')).strip())

    drive_identifier = 'Drive /c' + str(controller_index) + '/e' + str(enclosure) + '/s' + str(
        slot)
    if enclosure == ' ':
        drive_identifier = 'Drive /c' + str(controller_index) + '/s' + str(slot)
    try:
        info = detailed_info_array[drive_identifier + ' - Detailed Information']
        state = info[drive_identifier + ' State']
        attributes = info[drive_identifier + ' Device attributes']
        settings = info[drive_identifier + ' Policies/Settings']

        add_metric('pd_shield_counter', pd_baselabel, state['Shield Counter'])
        add_metric('pd_media_errors', pd_baselabel, state['Media Error Count'])
        add_metric('pd_other_errors', pd_baselabel, state['Other Error Count'])
        add_metric('pd_predictive_errors', pd_baselabel, state['Predictive Failure Count'])
        add_metric('pd_smart_alerted', pd_baselabel,
                   int(state['S.M.A.R.T alert flagged by drive'] == 'Yes'))
        add_metric('pd_link_speed_gbps', pd_baselabel, attributes['Link Speed'].split('.')[0])
        add_metric('pd_device_speed_gbps', pd_baselabel, attributes['Device Speed'].split('.')[0])
        add_metric('pd_commissioned_spare', pd_baselabel,
                   int(settings['Commissioned Spare'] == 'Yes'))
        add_metric('pd_emergency_spare', pd_baselabel, int(settings['Emergency Spare'] == 'Yes'))
        pd_info_label += ',firmware="{0}"'.format(attributes['Firmware Revision'].strip())
    except KeyError:
        pass
    add_metric('pd_info', pd_info_label, 1)


def add_metric(name, labels, value):
    global metric_list
    try:
        metric_list[name].append({
            'labels': labels,
            'value': float(value),
        })
    except ValueError:
        pass


def print_all_metrics(metrics):
    for metric, measurements in metrics.items():
        print('# HELP {0}{1} MegaRAID {2}'.format(metric_prefix, metric, metric.replace('_', ' ')))
        print('# TYPE {0}{1} gauge'.format(metric_prefix, metric))
        for measurement in measurements:
            if measurement['value'] != 'Unknown':
                print('{0}{1}{2} {3}'.format(metric_prefix, metric, '{' + measurement['labels'] + '}',
                                         measurement['value']))


def get_storcli_json(storcli_args):
    """Get storcli output in JSON format."""
    # Check if storcli is installed and executable
    if not (os.path.isfile(storcli_path) and os.access(storcli_path, os.X_OK)):
        SystemExit(1)
    storcli_cmd = shlex.split(storcli_path + ' ' + storcli_args)
    proc = subprocess.Popen(
        storcli_cmd, shell=False, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    output_json = proc.communicate()[0]
    data = json.loads(output_json.decode("utf-8"))

    if data["Controllers"][0]["Command Status"]["Status"] != "Success":
        SystemExit(1)
    return data


if __name__ == "__main__":
    PARSER = argparse.ArgumentParser(
        description=DESCRIPTION, formatter_class=argparse.ArgumentDefaultsHelpFormatter)
    PARSER.add_argument(
        '--storcli_path', default='/opt/MegaRAID/storcli/storcli64', help='path to StorCLi binary')
    PARSER.add_argument('--version', action='version', version='%(prog)s {0}'.format(VERSION))
    ARGS = PARSER.parse_args()

    main(ARGS)
