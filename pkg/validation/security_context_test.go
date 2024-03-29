package validation

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNSecurityContextValidatorValidate(t *testing.T) {
	t.Run("valid security context", func(t *testing.T) {
		trueVal := true
		falseVal := false
		pod := &corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name: "securePod",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  "secureContainer",
					Image: "busybox",
					SecurityContext: &corev1.SecurityContext{
						RunAsNonRoot:             &trueVal,
						AllowPrivilegeEscalation: &falseVal,
						Privileged:               &falseVal,
					},
				}},
			},
		}

		v, err := SecurityContextValidator{logger()}.Validate(pod)
		assert.Nil(t, err)
		assert.True(t, v.Valid)
	})

	t.Run("no security context", func(t *testing.T) {
		pod := &corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name: "no-security-context",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  "insecureContainer",
					Image: "busybox",
				}},
			},
		}

		v, err := SecurityContextValidator{logger()}.Validate(pod)
		assert.Nil(t, err)
		assert.False(t, v.Valid)
	})

	t.Run("Pod in kube-system namespace", func(t *testing.T) {
		os.Setenv("EXCLUDE_NAMESPACES", "kube-system")
		trueVal := true
		falseVal := false
		pod := &corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "securePod",
				Namespace: "kube-system",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  "secureContainer",
					Image: "busybox",
					SecurityContext: &corev1.SecurityContext{
						RunAsNonRoot:             &trueVal,
						AllowPrivilegeEscalation: &falseVal,
						Privileged:               &falseVal,
					},
				}},
			},
		}

		v, err := SecurityContextValidator{logger()}.Validate(pod)
		assert.Nil(t, err)
		assert.True(t, v.Valid)
	})
}
