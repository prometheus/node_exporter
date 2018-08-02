#!/bin/bash
######################################################################

user=`whoami`
if [ "${user}" != "root" ]; then
   echo "This script must run as root"
   exit -1
fi

VERSION="0.16.0"
if [ "$#" -gt 0 ]; then
    VERSION=$1
fi

echo "*** Build Snap of version ${VERSION} of prometheus-node-exporter ***"

# S/W version of prometheus-node-exporter used to find source tarball
export EXPORTER_VERSION=${VERSION}
rm -f snapcraft.yaml
sed -e "s/EXPORTER_VERSION/${EXPORTER_VERSION}/g" prometheus-node-exporter.yaml > snapcraft.yaml

snapcraft clean prometheus-node-exporter
rm -rf parts prime snap stage

# include manifest file in SNAP file
export SNAPCRAFT_BUILD_INFO=1
snapcraft snap

# remove temp file, does not remove build directories
snapcraft clean prometheus-node-exporter
rm -f snapcraft.yaml
