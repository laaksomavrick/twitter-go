#!/bin/sh

ABS_PATH=$(cd "$(dirname "$1")"; pwd -P)$(basename "$1") || exit

echo 'Running integration tests...'

for i in $(ls -v ${ABS_PATH}/tests/integration/*_test.go);
    do go test -v $i;
done;

echo 'Integration tests done.'