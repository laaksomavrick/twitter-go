#!/bin/sh

ABS_PATH=$(cd "$(dirname "$1")"; pwd -P)/$(basename "$1") || exit

echo 'Formatting with gofmt...'

# find ~ \*.sql
for i in $(find ${ABS_PATH} -name \*.go);
    # do cat $i | cqlsh;
    # do echo $i;
    do gofmt -w $i;
done;

echo 'Formatting done.'