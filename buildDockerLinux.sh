###
 # @Author: David
 # @version: 
 # @Description: 
 # @Date: 2020-01-23 16:19:26
 # @LastEditors: David
 # @LastEditTime: 2020-03-12 16:32:28
 ###
#!/bin/bash
VERSION=`date "+%Y%m%d%H"`
APP_NAME="node_exporter"
GO_ARCH=`go env| grep GOARCH|awk -F\" '{print $2}'`
cd `dirname $0`
WORKDIR=`pwd| sed 's#.*/src#/go/src#g'`

docker run --rm -it -v ${GOPATH}:/go -w ${WORKDIR} golang:latest go build  -o ${APP_NAME}.linux.${GO_ARCH} -v


echo "Linux Packaging Binaries..."
mkdir -p tmp/${APP_NAME}
mv ${APP_NAME}.linux.${GO_ARCH} tmp/${APP_NAME}/
cp README.md tmp/${APP_NAME}/ 
mkdir -p ./dist/
#tar -czf $@ -C tmp $(APP_NAME);
tar -cvzf ${APP_NAME}-${VERSION}.linux.${GO_ARCH}.tar.gz -C tmp  .
mv ${APP_NAME}-${VERSION}.linux.${GO_ARCH}.tar.gz ./dist/
rm -rf tmp
echo
echo "Package ${APP_NAME}-${VERSION}.linux.${GO_ARCH}.tar.gz saved in dist directory"