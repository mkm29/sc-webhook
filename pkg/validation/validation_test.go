package validation

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidatePodImage(t *testing.T) {
	// set REGISTRY env var
	os.Setenv("REGISTRY", "docker.io")
	ev, ok := GetRegistry()
	if !ok {
		t.Errorf("%s environment variable is not set", REGISTRY)
	}
	v := NewValidator(logger())

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "secure",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "secure-container",
				Image: fmt.Sprintf("%s/busybox", ev),
			}},
		},
	}

	validations := []PodValidator{
		ImageValidator{v.Logger},
	}
	val, err := v.ValidatePod(pod, validations)
	assert.Nil(t, err)
	assert.True(t, val.Valid)

	// change the image to a bad image
	pod.Spec.Containers[0].Image = "nginx"

	val, err = v.ValidatePod(pod, validations)
	assert.Nil(t, err)
	assert.False(t, val.Valid)
}

func TestValidatePodSecurityContext(t *testing.T) {
	v := NewValidator(logger())
	trueVal := true
	falseVal := false

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "secure",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "insecure-container",
				Image: "busybox",
			}},
		},
	}

	validations := []PodValidator{
		SecurityContextValidator{v.Logger},
	}
	val, err := v.ValidatePod(pod, validations)
	assert.Nil(t, err)
	assert.False(t, val.Valid)

	// add security context to container
	pod.Spec.Containers[0].SecurityContext = &corev1.SecurityContext{
		RunAsNonRoot:             &trueVal,
		Privileged:               &falseVal,
		AllowPrivilegeEscalation: &falseVal,
	}
	val, err = v.ValidatePod(pod, validations)
	assert.Nil(t, err)
	assert.True(t, val.Valid)
}

func logger() *logrus.Entry {
	mute := logrus.StandardLogger()
	mute.Out = ioutil.Discard
	return mute.WithField("logger", "test")
}
