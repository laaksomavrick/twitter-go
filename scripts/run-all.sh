#!/bin/sh

trap 'killall' INT

killall() {
    trap '' INT TERM     # ignore INT and TERM while shutting down
    echo "**** Shutting down... ****"     # added double quotes
    kill -TERM 0         # fixed order, send TERM not INT
    wait
    echo DONE
}

ABS_PATH=$(cd "$(dirname "$1")"; pwd -P)$(basename "$1") || exit

echo 'Booting go services...'

for i in $(find ${ABS_PATH}/services/**/cmd/*.go); do 
    echo "Booting $i";
    go run $i &
done;

cat