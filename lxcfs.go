package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultAnnotation = "initializer.kubernetes.io/lxcfs"
	defaultNamespace  = "default"
)

var (
	annotation        string
	configmap         string
	initializerName   string
	namespace         string
	requireAnnotation bool
)

// -v /var/lib/lxcfs/proc/cpuinfo:/proc/cpuinfo:ro
// -v /var/lib/lxcfs/proc/diskstats:/proc/diskstats:ro
// -v /var/lib/lxcfs/proc/meminfo:/proc/meminfo:ro
// -v /var/lib/lxcfs/proc/stat:/proc/stat:ro
// -v /var/lib/lxcfs/proc/swaps:/proc/swaps:ro
// -v /var/lib/lxcfs/proc/uptime:/proc/uptime:ro
// -v /var/lib/lxcfs/proc/loadavg:/proc/loadavg:ro
var volumeMountsTemplate = []corev1.VolumeMount{

	{
		Name:      "lxcfs-proc-cpuinfo",
		MountPath: "/proc/cpuinfo",
		ReadOnly:  true,
	},
	{
		Name:      "lxcfs-proc-meminfo",
		MountPath: "/proc/meminfo",
		ReadOnly:  true,
	},
	{
		Name:      "lxcfs-proc-diskstats",
		MountPath: "/proc/diskstats",
		ReadOnly:  true,
	},
	{
		Name:      "lxcfs-proc-stat",
		MountPath: "/proc/stat",
		ReadOnly:  true,
	},
    // docker-ce do not allow to mount 'swap' and 'loadavg'
    // reference: https://github.com/docker/docker-ce/blob/master/components/engine/contrib/apparmor/template.go line 113
	//{
	//	Name:      "lxcfs-proc-swaps",
	//	MountPath: "/proc/swaps",
	//	ReadOnly:  true,
	//},
	{
		Name:      "lxcfs-proc-uptime",
		MountPath: "/proc/uptime",
		ReadOnly:  true,
	},
	//{
	//	Name:      "lxcfs-proc-loadavg",
	//	MountPath: "/proc/loadavg",
	//	ReadOnly:  true,
	//},
	{
		Name:      "lxcfs-sys-devices-system-cpu-online",
		MountPath: "/sys/devices/system/cpu/online",
		ReadOnly:  true,
	},
}
var volumesTemplate = []corev1.Volume{
	{
		Name: "lxcfs-proc-cpuinfo",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/lib/lxcfs/proc/cpuinfo",
			},
		},
	},
	{
		Name: "lxcfs-proc-diskstats",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/lib/lxcfs/proc/diskstats",
			},
		},
	},
	{
		Name: "lxcfs-proc-meminfo",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/lib/lxcfs/proc/meminfo",
			},
		},
	},
	{
		Name: "lxcfs-proc-stat",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/lib/lxcfs/proc/stat",
			},
		},
	},
	//{
	//	Name: "lxcfs-proc-swaps",
	//	VolumeSource: corev1.VolumeSource{
	//		HostPath: &corev1.HostPathVolumeSource{
	//			Path: "/var/lib/lxcfs/proc/swaps",
	//		},
	//	},
	//},
	{
		Name: "lxcfs-proc-uptime",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/lib/lxcfs/proc/uptime",
			},
		},
	},
	//{
	//	Name: "lxcfs-proc-loadavg",
	//	VolumeSource: corev1.VolumeSource{
	//		HostPath: &corev1.HostPathVolumeSource{
	//			Path: "/var/lib/lxcfs/proc/loadavg",
	//		},
	//	},
	//},
	{
		Name: "lxcfs-sys-devices-system-cpu-online",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/lib/lxcfs/sys/devices/system/cpu/online",
			},
		},
	},
}

// main mutation process
func (whsvr *WebhookServer) mutatePod(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	req := ar.Request
	var (
		objectMeta                      *metav1.ObjectMeta
		resourceNamespace, resourceName string
	)

	glog.Infof("AdmissionReview for Kind=%v, Namespace=%v Name=%v (%v) UID=%v patchOperation=%v UserInfo=%v",
		req.Kind, req.Namespace, req.Name, resourceName, req.UID, req.Operation, req.UserInfo)

	var pod corev1.Pod

	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		glog.Errorf("Could not unmarshal raw object to pod: %v", err)
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}
	resourceName, resourceNamespace, objectMeta = pod.Name, pod.Namespace, &pod.ObjectMeta

	if !mutationRequired(ignoredNamespaces, objectMeta) {
		glog.Infof("Skipping validation for %s/%s due to policy check", resourceNamespace, resourceName)
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}

	patchBytes, err := createPodPatch(&pod)
	if err != nil {
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	patchType := v1beta1.PatchTypeJSONPatch

	glog.Infof("AdmissionResponse: patch=%v\n", string(patchBytes))
	return &v1beta1.AdmissionResponse{
		UID:       req.UID,
		Allowed:   true,
		Patch:     patchBytes,
		PatchType: &patchType,
	}
}

func createPodPatch(pod *corev1.Pod) ([]byte, error) {

	var patches []patchOperation

	var op = patchOperation{
		Path: "/metadata/annotations",
		Value: map[string]string{
			admissionWebhookAnnotationStatusKey: "mutated",
		},
	}

	if pod.Annotations == nil {
		op.Op = "add"
	} else {
		op.Op = "add"
		if pod.Annotations[admissionWebhookAnnotationStatusKey] != "" {
			op.Op = " replace"
		}
		op.Path = "/metadata/annotations/" + escapeJSONPointerValue(admissionWebhookAnnotationStatusKey)
		op.Value = "mutated"
	}

	patches = append(patches, op)

	containers := pod.Spec.Containers

	// Modify the Pod spec to include the LXCFS volumes, then op the original pod.
	for i := range containers {
		if containers[i].VolumeMounts == nil {
			path := fmt.Sprintf("/spec/containers/%d/volumeMounts", i)
			op = patchOperation{
				Op:    "add",
				Path:  path,
				Value: volumeMountsTemplate,
			}
			patches = append(patches, op)
		} else {
			path := fmt.Sprintf("/spec/containers/%d/volumeMounts/-", i)
			for _, volumeMount := range volumeMountsTemplate {
				op = patchOperation{
					Op:    "add",
					Path:  path,
					Value: volumeMount,
				}
				patches = append(patches, op)
			}
		}
	}

	if pod.Spec.Volumes == nil {
		op = patchOperation{
			Op:    "add",
			Path:  "/spec/volumes",
			Value: volumesTemplate,
		}
		patches = append(patches, op)
	} else {
		for _, volume := range volumesTemplate {
			op = patchOperation{
				Op:    "add",
				Path:  "/spec/volumes/-",
				Value: volume,
			}
			patches = append(patches, op)
		}
	}

	patchBytes, err := json.Marshal(patches)
	if err != nil {
		glog.Warningf("error in json.Marshal %s: %v", pod.Name, err)
		return nil, err
	}
	return patchBytes, nil
}

// validate deployments and services
func (whsvr *WebhookServer) validatePod(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Allowed: true,
	}
}

func escapeJSONPointerValue(in string) string {
	step := strings.Replace(in, "~", "~0", -1)
	return strings.Replace(step, "/", "~1", -1)
}
