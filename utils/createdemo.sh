#!/bin/bash

echo
echo "Creating namespaces..."
for i in resources/namespace*; do kubectl create -f resources/${i##*/}; done

echo
echo "Creating pods..."
for i in resources/pod*; do kubectl create -f  resources/${i##*/}; done

echo
echo "Labeling pods..."
kubectl label pods redis destiny=doomed
kubectl label pods mysql destiny=doomed --namespace development
kubectl label pods caddy destiny=doomed --namespace staging
