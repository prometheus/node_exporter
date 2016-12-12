#!/bin/sh
#
# Attempts to generate static binary executable suitable for use in docker images.
#
# Needed for generating docker images in non-linux environments (osx, windows, etc.)
#
# Usage: ./build_docker_executable.sh [options]
#
# Options:
#
#   stripped
#     produces a smaller binary, stripped of debug info
#


GITHUB_USER="prometheus"
GITHUB_REPO="node_exporter"
LOCAL_TARGET="./node_exporter"

DOCKER_TAG="${GITHUB_REPO}_build"
GITHUB_PATH="${GITHUB_USER}/${GITHUB_REPO}"
EXE_PATH="/bin/${GITHUB_REPO}"

LDFLAGS=""
for arg in "${@}"
do
	case $arg in
		stripped)
			LDFLAGS="${LDFLAGS} -s"
			;;
	esac
	
done

export GO_LD_FLAGS

echo "Building static binary executable as docker image build..."
docker build -t ${DOCKER_TAG} --build-arg ldflags="${LDFLAGS}" -f Dockerfile.build .
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
