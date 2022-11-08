package validation

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNImageValidatorValidate(t *testing.T) {
	t.Run("valid image sources", func(t *testing.T) {
		// set the REGISTRY environment variable
		os.Setenv(REGISTRY, "private.docker.io")
		ev := getRegistry()
		containers := []corev1.Container{
			{
				Name:  "good-container-1",
				Image: fmt.Sprintf("%s/busybox", ev),
			},
			{
				Name:  "bad-container",
				Image: fmt.Sprintf("%s/nginx", ev),
			},
		}
		pod := &corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name: "approved-pod",
			},
			Spec: corev1.PodSpec{
				Containers: containers,
			},
		}

		v, err := ImageValidator{logger()}.Validate(pod)
		assert.Nil(t, err)
		assert.True(t, v.Valid)
	})

	t.Run("image not from an approved registry", func(t *testing.T) {
		os.Setenv(REGISTRY, "private.docker.io")
		ev := getRegistry()
		os.Setenv(REGISTRY, ev)
		// define 2 containers
		containers := []corev1.Container{
			{
				Name:  "good-container",
				Image: fmt.Sprintf("%s/busybox", ev),
			},
			{
				Name:  "bad-container",
				Image: "busybox",
			},
		}
		pod := &corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name: "image-not-approved-pod",
			},
			Spec: corev1.PodSpec{
				Containers: containers,
			},
		}

		v, err := SecurityContextValidator{logger()}.Validate(pod)
		assert.Nil(t, err)
		assert.False(t, v.Valid)
	})
}

func TestExcludeNamespacePod(t *testing.T) {
	t.Run("Pod in kube-system namespace", func(t *testing.T) {
		os.Setenv("EXCLUDE_NAMESPACES", "kube-system")
		containers := []corev1.Container{
			{
				Name:  "good-container-1",
				Image: "busybox",
			},
		}
		pod := &corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "approved-pod",
				Namespace: "kube-system",
			},
			Spec: corev1.PodSpec{
				Containers: containers,
			},
		}

		v, err := ImageValidator{logger()}.Validate(pod)
		assert.Nil(t, err)
		assert.True(t, v.Valid)
	})
}
