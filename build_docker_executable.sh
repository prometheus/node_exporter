#!/bin/sh
#
# Attempts to generate static binary executable suitable for use in docker images.
#
# Needed for generating docker images in non-linux environments (osx, windows, etc.)
#
# Usage: ./build_docker_executable.sh [options] ...
#
# Options:
#
#   ldflags=<value>
#     values to pass with -ldflags for "go build"
#
#   tags=<value>
#     values to pass with -tags for "go build" (-s for stripped binary executable)
#
# Example:
#
#   ./build_docker_executable.sh ldflags="-s" tags="noloadavg notime"
#

GITHUB_USER="prometheus"
GITHUB_REPO="node_exporter"
LOCAL_TARGET="./node_exporter"

DOCKER_TAG="${GITHUB_REPO}_build"
GITHUB_PATH="${GITHUB_USER}/${GITHUB_REPO}"
EXE_PATH="/bin/${GITHUB_REPO}"

LDFLAGS=""
TAGS=""
for arg in "${@}"
do
	echo "$arg" | grep = > /dev/null
	if [ $? != 0 ]
	then
		echo "Invalid argument passed: $arg"
		exit 1
	fi
	key=`echo "$arg" | awk -F= '{ print $1 }'`
	value=`echo "$arg" | awk -F= '{ $1=""; print substr($0,1) }'`
	if [ "$key" = "" ]
	then
		echo "Empty key passed.  Expected key=value."
		exit 1
	fi
	if [ "$value" = "" ]
	then
		echo "Empty value passed for $key.  Expected key=value."
		exit 1
	fi
	case $key in
		ldflags)
			LDFLAGS="$value"
			;;
		tags)
			TAGS="$value"
			;;
		*)
			echo "Invalid option: $key"
			exit 1
			;;
	esac
	
done

echo "Building static binary executable as docker image build..."
docker build -t ${DOCKER_TAG} --build-arg ldflags="${LDFLAGS}" --build-arg tags="${TAGS}" -f Dockerfile.build .
if [ $? != 0 ]
then
	echo "Failed to build docker image with binary executable" >&2
	exit 1
fi

echo "Running idle build image container in order to access its contents..."
CONTAINER_ID=`docker run -d -t ${DOCKER_TAG}`
if [ -z "${CONTAINER_ID}" ]
then
	echo "Failed to get static build container id" >&2
	exit 1
fi

echo "Copying static binary executable out of build image container..."
docker cp ${CONTAINER_ID}:${EXE_PATH} ${LOCAL_TARGET}
if [ $? != 0 ]
then
	echo "Failed to cp ${EXE_PATH} out of static build container" >&2
	exit 1
fi

echo "Stopping build image container..."
docker stop ${CONTAINER_ID}
if [ $? != 0 ]
then
	echo "Failed to stop static build container" >&2
	exit 1
fi

echo "Destroying build image container..."
docker rm ${CONTAINER_ID}
if [ $? != 0 ]
then
	echo "Failed to remove static build container" >&2
	exit 1
fi

echo "Getting build image id..."
IMAGE_ID=`docker images -q ${DOCKER_TAG}`
if [ -z "${IMAGE_ID}" ]
then
	echo "Failed to determine static build image id" >&2
	exit 1
fi

echo "Destroying build image..."
docker rmi ${IMAGE_ID}
if [ $? != 0 ]
then
	echo "Failed to remove static build image" >&2
	exit 1
fi

echo "Finished build of ${LOCAL_TARGET}"
