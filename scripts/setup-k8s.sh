#!/bin/sh

# This script assumes that the user has kubectl set up in their path, and that kubectl is configured properly.

# Set up k8s namespace
kubectl create namespace twtr-dev

# Set up tiller with RBAC to this namespace
kubectl create serviceaccount tiller --namespace twtr-dev
kubectl create -f role-tiller.yaml
kubectl create -f rolebinding-tiller.yaml
helm init --service-account tiller --tiller-namespace twtr-dev

# Create a storage class for persistent volumes
kubectl create -f create-storage-gce.yaml