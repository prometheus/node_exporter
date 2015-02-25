#!/bin/bash

# We strip the tag part of git describe and append it to the currently defined
# version in the makefile
DESCRIBE=`git describe`
FROM_TAG_INDEX=`expr index "$DESCRIBE" "-"`
echo ${DESCRIBE:$FROM_TAG_INDEX}
