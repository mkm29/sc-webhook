package mutation

import (
	"github.com/mkm29/sc-webhook/pkg/utils"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

func NewTrue() *bool {
	b := true
	return &b
}
func NewFalse() *bool {
	b := false
	return &b
}

// injectSecurityContext is a container for the mutation injecting environment vars
type injectSecurityContext struct {
	Logger logrus.FieldLogger
}

// injectSecurityContext implements the podMutator interface
var _ podMutator = (*injectSecurityContext)(nil)

// Name returns the struct name
func (sc injectSecurityContext) Name() string {
	return "inject_security_context"
}

// Mutate returns a new mutated pod according to set env rules
func (sc injectSecurityContext) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {
	sc.Logger = sc.Logger.WithField("mutation", sc.Name())
	mpod := pod.DeepCopy()

	// this needs to be done on the Pod level, not the container level
	securityContext := corev1.SecurityContext{
		Privileged:               NewFalse(),
		RunAsNonRoot:             NewTrue(),
		AllowPrivilegeEscalation: NewFalse(),
	}

	// inject env vars into pod
	sc.Logger.Debugf("pod security context injected %s", securityContext)
	injectValidSecurityContext(mpod, securityContext)

	return mpod, nil
}

// injectSecurityContextVar injects a var in both containers and init containers of a pod
func injectValidSecurityContext(pod *corev1.Pod, sc corev1.SecurityContext) {
	if !utils.HasValidSecurityContext(pod) {
		// inject the security context into each container
		for i := range pod.Spec.Containers {
			pod.Spec.Containers[i].SecurityContext = &sc
		}
	}
}
