#!/bin/sh

ABS_PATH=$(cd "$(dirname "$1")"; pwd -P)$(basename "$1") || exit
HOST_NAME=gcr.io
GOOGLE_PROJECT_ID=precise-clock-244301
TAG=latest

echo 'Building docker images...'

for i in $(ls -v ${ABS_PATH}/services/*/Dockerfile); do
    IMAGE_NAME=twtr-$(basename $(dirname $i))
    docker build -f $i -t ${IMAGE_NAME} .
    docker tag ${IMAGE_NAME} ${HOST_NAME}/${GOOGLE_PROJECT_ID}/${IMAGE_NAME}:${TAG}
    docker push ${HOST_NAME}/${GOOGLE_PROJECT_ID}/${IMAGE_NAME}
done;

echo 'Done building docker images.'