#!/bin/bash
_tag=`git describe`
another_var=`expr index "$_tag" "-"`
echo ${_tag:$another_var}
