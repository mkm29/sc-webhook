package mutation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInjectSecurityContextMutate(t *testing.T) {
	trueVal := true
	want := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{
			SecurityContext: &corev1.PodSecurityContext{
				RunAsNonRoot: &trueVal,
			},
			Containers: []corev1.Container{{
				Name: "test",
			}},
		},
	}

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name: "test",
			}},
		},
	}

	got, err := injectSecurityContext{Logger: logger()}.Mutate(pod)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}

func TestHasSecurityContext(t *testing.T) {
	trueVal := true
	c := corev1.Container{
		Name:  "test",
		Image: "busybox",
	}

	// create Pod spec
	pn := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{c},
		},
	}
	// create new Pod spec with security context
	py := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{c},
			SecurityContext: &corev1.PodSecurityContext{
				RunAsNonRoot: &trueVal,
			},
		},
	}

	assert.True(t, HasValidSecurityContext(py))
	assert.False(t, HasValidSecurityContext(pn))
}
