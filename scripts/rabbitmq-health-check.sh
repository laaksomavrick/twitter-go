#!/bin/sh

# Scratch file for figuring out rabbit init container.
# Not in use, see ready service. Keeping around for posterity and as a testament to my efforts.

i=0;
while [ $i -le 10 ]; do
    RABBIT_STATE=$(curl --max-time 8 --silent --show-error --fail ${AMQP_MGMT_URL});
    if [[ $RABBIT_STATE == *no* ]]; then
        echo ok;
        exit 0;
    fi
    echo waiting...;
    ((i++));
    sleep $i;
done;
exit 1;