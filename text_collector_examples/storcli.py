#!/usr/bin/env python
"""
Script to parse StorCLI's JSON output and expose
MegaRAID health as Prometheus metrics.

Tested against StorCLI 'Ver 1.14.12 Nov 25, 2014'.

StorCLI reference manual:
http://docs.avagotech.com/docs/12352476

Advanced Software Options (ASO) not exposed as metrics currently.

JSON key abbreviations used by StorCLI are documented in the standard command
output, i.e.  when you omit the trailing 'J' from the command.
"""

from __future__ import print_function
import argparse
import json
import os
import subprocess

DESCRIPTION = """Parses StorCLI's JSON output and exposes MegaRAID health as
    Prometheus metrics."""
VERSION = '0.0.1'


def main(args):
    """ main """

    # exporter variables
    metric_prefix = 'megaraid_'
    metric_controller_labels = '{{controller="{}", model="{}"}}'

    data = json.loads(get_storcli_json(args.storcli_path))

    # It appears that the data we need will always be present in the first
    # item in the Controllers array
    status = data['Controllers'][0]

    metrics = {
        'status_code': status['Command Status']['Status Code'],
        'controllers': status['Response Data']['Number of Controllers'],
    }

    for name, value in metrics.iteritems():
        print('# HELP {}{} MegaRAID {}'.format(metric_prefix, name, name.replace('_', ' ')))
        print('# TYPE {}{} gauge'.format(metric_prefix, name))
        print("{}{} {}".format(metric_prefix, name, value))

    controller_info = []
    controller_metrics = {}
    overview = []

    try:
        overview = status['Response Data']['System Overview']
    except KeyError:
        pass

    for controller in overview:
        controller_index = controller['Ctl']
        model = controller['Model']
        controller_info.append(metric_controller_labels.format(controller_index, model))

        controller_metrics = {
            # FIXME: Parse dimmer switch options
            # 'dimmer_switch':          controller['DS'],

            'battery_backup_healthy':   int(controller['BBU'] == 'Opt'),
            'degraded':                 int(controller['Hlth'] == 'Dgd'),
            'drive_groups':             controller['DGs'],
            'emergency_hot_spare':      int(controller['EHS'] == 'Y'),
            'failed':                   int(controller['Hlth'] == 'Fld'),
            'healthy':                  int(controller['Hlth'] == 'Opt'),
            'physical_drives':          controller['PDs'],
            'ports':                    controller['Ports'],
            'scheduled_patrol_read':    int(controller['sPR'] == 'On'),
            'virtual_drives':           controller['VDs'],

            # Reverse StorCLI's logic to make metrics consistent
            'drive_groups_optimal':     int(controller['DNOpt'] == 0),
            'virtual_drives_optimal':   int(controller['VNOpt'] == 0),
            }

    for name, value in controller_metrics.iteritems():
        print('# HELP {}{} MegaRAID {}'.format(metric_prefix, name, name.replace('_', ' ')))
        print('# TYPE {}{} gauge'.format(metric_prefix, name))
        print('{}{}{{controller="{}"}} {}'.format(metric_prefix, name,
                                                  controller_index, value))

    if controller_info:
        print('# HELP {}{} MegaRAID controller info'.format(metric_prefix, 'controller_info'))
        print('# TYPE {}{} gauge'.format(metric_prefix, 'controller_info'))
    for labels in controller_info:
        print('{}{}{} {}'.format(metric_prefix, 'controller_info', labels, 1))


def get_storcli_json(storcli_path):
    """Get storcli output in JSON format."""

    # Check if storcli is installed
    if os.path.isfile(storcli_path) and os.access(storcli_path, os.X_OK):
        storcli_cmd = [storcli_path, 'show', 'all', 'J']
        proc = subprocess.Popen(storcli_cmd, shell=False,
                                stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        output_json = proc.communicate()[0]
    else:
        # Create an empty dummy-JSON where storcli not installed.
        dummy_json = {"Controllers":[{
            "Command Status": {"Status Code": 0, "Status": "Success",
                               "Description": "None"},
            "Response Data": {"Number of Controllers": 0}}]}
        output_json = json.dumps(dummy_json)

    return output_json

if __name__ == "__main__":
    PARSER = argparse.ArgumentParser(description=DESCRIPTION,
                                     formatter_class=argparse.ArgumentDefaultsHelpFormatter)
    PARSER.add_argument('--storcli_path',
                        default='/opt/MegaRAID/storcli/storcli64',
                        help='path to StorCLi binary')
    PARSER.add_argument('--version',
                        action='version',
                        version='%(prog)s {}'.format(VERSION))
    ARGS = PARSER.parse_args()

    main(ARGS)
