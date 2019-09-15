package main

import (
	"encoding/json"
	jsonpatch "github.com/evanphx/json-patch"

	corev1 "k8s.io/api/core/v1"

	"testing"
)

func TestCreatePodPatch(t *testing.T) {

	var pod corev1.Pod

	pod.Name = "test"
	pod.Namespace = "testNS"
	pod.Annotations = map[string]string{}
	pod.Annotations[admissionWebhookAnnotationStatusKey] = "testAnnotation"
	pod.Spec.Containers = []corev1.Container{
		{
			Image: "test_image",
		},
	}
	testCreatePodPatch(t, &pod)
}

func TestCreatePodPatch2(t *testing.T) {

	var pod corev1.Pod

	pod.Name = "test"
	pod.Namespace = "testNS"
	pod.Annotations = map[string]string{}
	pod.Annotations[admissionWebhookAnnotationStatusKey] = "testAnnotation"
	pod.Spec.Containers = []corev1.Container{
		{
			Image: "test_image",
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "test",
					MountPath: "/etc/test",
				},
			},
		},
	}
	pod.Spec.Volumes = []corev1.Volume{
		{
			Name: "test",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/test",
				},
			},
		},
	}
	testCreatePodPatch(t, &pod)
}

func testCreatePodPatch(t *testing.T, pod *corev1.Pod) {
	patchJSON, err := createPodPatch(pod)

	if err != nil {
		t.Errorf("error in createPodPatch: %v", err)
	} else {
		t.Logf("patch :\n%s", string(patchJSON))
	}

	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		t.Errorf("error in createPodPatch: %v", err)
	}

	podJSON, err := json.Marshal(pod)
	if err != nil {
		t.Errorf("error in json.Marshal: %v", err)
	}

	modified, err := patch.Apply(podJSON)
	if err != nil {
		t.Errorf("error in createPodPatch: %v", err)
	}

	t.Logf("modified Pod:\n%s", string(modified))

	var modifiedPod corev1.Pod
	err = json.Unmarshal(modified, &modifiedPod)
	if err != nil {
		t.Errorf("error in createPodPatch: %v", err)
	}

	modifiedJSON, err := json.Marshal(modifiedPod)
	if err != nil {
		t.Errorf("error in json.Marshal: %v", err)
	}
	t.Logf("modified Pod:\n%s", string(modifiedJSON))

}
