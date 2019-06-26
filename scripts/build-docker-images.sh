#!/bin/sh

ABS_PATH=$(cd "$(dirname "$1")"; pwd -P)$(basename "$1") || exit

echo 'Building docker images...'

for i in $(ls -v ${ABS_PATH}/services/*/Dockerfile);
    do docker build -f $i -t twtr-$(basename $(dirname $i)) .
done;

echo 'Done building docker images.'