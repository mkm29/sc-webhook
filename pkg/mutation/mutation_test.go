package mutation

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMutatePodPatch(t *testing.T) {
	t.Run("mutate security context", func(t *testing.T) {
		m := NewMutator(logger())
		got, err := m.MutatePodPatch(pod())
		if err != nil {
			t.Fatal(err)
		}

		p := patch()
		g := string(got)
		assert.Equal(t, p, g)
	})

	t.Run("protected namespace do not mutate", func(t *testing.T) {
		m := NewMutator(logger())
		p := pod()
		p.ObjectMeta.Namespace = "kube-system"
		got, err := m.MutatePodPatch(p)
		if err != nil {
			t.Fatal(err)
		}
		// should be nil
		assert.Nil(t, got)
	})
}

func BenchmarkMutatePodPatch(b *testing.B) {
	m := NewMutator(logger())
	pod := pod()

	for i := 0; i < b.N; i++ {
		_, err := m.MutatePodPatch(pod)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func pod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "insecurePod",
			Labels: map[string]string{
				"dx.rtx.com/security-requested": "true",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "insecureContainer",
				Image: "busybox",
			}},
		},
	}
}

func patch() string {
	patch := `[
			{"op":"add","path":"/spec/containers/0/securityContext","value":
				{"allowPrivilegeEscalation":false,"privileged":false,"runAsNonRoot":true}
			}
]`

	patch = strings.ReplaceAll(patch, "\n", "")
	patch = strings.ReplaceAll(patch, "\t", "")
	patch = strings.ReplaceAll(patch, " ", "")

	return patch
}

func logger() *logrus.Entry {
	mute := logrus.StandardLogger()
	mute.Out = ioutil.Discard
	return mute.WithField("logger", "test")
}
