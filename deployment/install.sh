#!/bin/bash

./deployment/webhook-create-signed-cert.sh
kubectl get secret lxcfs-admission-webhook-certs

kubectl create -f deployment/deployment.yaml
kubectl create -f deployment/service.yaml
cat ./deployment/mutatingwebhook.yaml | ./deployment/webhook-patch-ca-bundle.sh > ./deployment/mutatingwebhook-ca-bundle.yaml
kubectl create -f deployment/mutatingwebhook-ca-bundle.yaml

