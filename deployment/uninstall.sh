#!/bin/bash

kubectl delete -f deployment/mutatingwebhook-ca-bundle.yaml
kubectl delete -f deployment/service.yaml
kubectl delete -f deployment/deployment.yaml
kubectl delete secret lxcfs-admission-webhook-certs
kubectl delete CertificateSigningRequest lxcfs-admission-webhook-svc.default
kubectl delete MutatingWebhookConfiguration mutating-lxcfs-admission-webhook-cfg -n default
kubectl delete ValidatingWebhookConfiguration validation-lxcfs-admission-webhook-cfg -n default
