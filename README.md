# Kubernetes Admission Webhook for LXCFS

This project shows how to build and deploy an [AdmissionWebhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#admission-webhooks) for [LXCFS](https://github.com/lxc/lxcfs).

## Prerequisites

Kubernetes 1.9.0 or above with the `admissionregistration.k8s.io/v1beta1` API enabled. Verify that by the following command:
```
kubectl api-versions | grep admissionregistration.k8s.io/v1beta1
```
The result should be:
```
admissionregistration.k8s.io/v1beta1
```

In addition, the `MutatingAdmissionWebhook` and `ValidatingAdmissionWebhook` admission controllers should be added and listed in the correct order in the admission-control flag of kube-apiserver.

## Build

1. Setup dep

   The repo uses [dep](https://github.com/golang/dep) as the dependency management tool for its Go codebase. Install `dep` by the following command:

```
go get -u github.com/golang/dep/cmd/dep
```

2. Build and push docker image
   
```
./build
```

## Deploy 
 
1. Deploy lxcfs to worker nodes

```
kubectl apply -f deployment/lxcfs-daemonset.yaml
```

2. Install injector with lxcfs-admission-webhook

```
deployment/install.sh
```

## Test

1. Enable the namespace for injection

```
kubectl label namespace default lxcfs-admission-webhook=enabled
```

Note: All the new created pod under the namespace will be injected with LXCFS


2. Deploy the test deployment
 
```
kubectl apply -f deployment/web.yaml
```

3. Inspect the resource inside container


```
$ kubectl get pod

NAME                                                 READY   STATUS    RESTARTS   AGE
lxcfs-admission-webhook-deployment-f4bdd6f66-5wrlg   1/1     Running   0          8m29s
lxcfs-pqs2d                                          1/1     Running   0          55m
lxcfs-zfh99                                          1/1     Running   0          55m
web-7c5464f6b9-6zxdf                                 1/1     Running   0          8m10s
web-7c5464f6b9-nktff                                 1/1     Running   0          8m10s

$ kubectl exec -ti web-7c5464f6b9-6zxdf sh
# free
             total       used       free     shared    buffers     cached
Mem:        262144       2744     259400          0          0        312
-/+ buffers/cache:       2432     259712
Swap:            0          0          0
#
```

## Cleanup

1. Uninstall injector with lxcfs-admission-webhook

```
deployment/uninstall.sh
```

2. Uninstall lxcfs from cluster nodes

```
kubectl delete -f deployment/lxcfs-daemonset.yaml
```

## How does it work?

If you want to know webhooks in depth, please check [it](https://aliyun.com/blog/k8s-admission-webhooks/) out!


