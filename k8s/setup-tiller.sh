#!/bin/sh
kubectl create namespace twtr-dev
kubectl create serviceaccount tiller --namespace twtr-dev
kubectl create -f role-tiller.yaml
kubectl create -f rolebinding-tiller.yaml
helm init --service-account tiller --tiller-namespace twtr-dev