#!/bin/sh

ABS_PATH=$(cd "$(dirname "$1")"; pwd -P)$(basename "$1") || exit

echo 'Running migrations...'

for i in $(ls -v ${ABS_PATH}/migrations/*.cql);
    do cat $i | cqlsh;
done;

echo 'Migrations done.'