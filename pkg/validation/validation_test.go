package validation

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidatePod(t *testing.T) {
	ev, ok := GetRegistry()
	if !ok {
		t.Errorf("%s environment variable is not set", REGISTRY)
	}
	v := NewValidator(logger())
	trueVal := true

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "secure",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "secure-container",
				Image: fmt.Sprintf("%s/busybox", ev),
			}},
			SecurityContext: &corev1.PodSecurityContext{
				RunAsNonRoot: &trueVal,
			},
		},
	}

	val, err := v.ValidatePod(pod)
	assert.Nil(t, err)
	assert.True(t, val.Valid)

	// change the image to a bad image
	pod.Spec.Containers[0].Image = "nginx"

	val, err = v.ValidatePod(pod)
	assert.Nil(t, err)
	assert.False(t, val.Valid)
}

func logger() *logrus.Entry {
	mute := logrus.StandardLogger()
	mute.Out = ioutil.Discard
	return mute.WithField("logger", "test")
}
