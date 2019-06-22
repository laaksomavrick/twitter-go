#!/bin/sh

kubectl get pods | grep Evicted | awk '{print $1}' | xargs kubectl delete pod