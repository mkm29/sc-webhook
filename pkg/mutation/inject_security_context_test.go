package mutation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInjectSecurityContextMutate(t *testing.T) {
	trueVal := true
	falseVal := false
	want := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{
			// SecurityContext: &corev1.PodSecurityContext{
			// 	RunAsNonRoot: &trueVal,
			// },
			Containers: []corev1.Container{{
				Name: "test",
				SecurityContext: &corev1.SecurityContext{
					AllowPrivilegeEscalation: &falseVal,
					Privileged:               &falseVal,
					RunAsNonRoot:             &trueVal,
				},
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
	falseVal := false
	c := corev1.Container{
		Name:  "test",
		Image: "busybox",
		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: &falseVal,
			RunAsNonRoot:             &trueVal,
			Privileged:               &falseVal,
		},
	}

	// create Pod spec
	pn := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				corev1.Container{
					Name:  "test",
					Image: "busybox",
				},
			},
		},
	}
	// create new Pod spec with security context
	py := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{c},
		},
	}

	assert.True(t, HasValidSecurityContext(py))
	assert.False(t, HasValidSecurityContext(pn))
}
