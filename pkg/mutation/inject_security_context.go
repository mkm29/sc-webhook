package mutation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

func NewTrue() *bool {
	b := true
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
	// create a pod security context
	securityContext := corev1.PodSecurityContext{
		RunAsNonRoot: NewTrue(),
	}

	// inject env vars into pod
	sc.Logger.Debugf("pod security context injected %s", securityContext)
	injectValidSecurityContext(mpod, securityContext)

	return mpod, nil
}

// injectSecurityContextVar injects a var in both containers and init containers of a pod
func injectValidSecurityContext(pod *corev1.Pod, sc corev1.PodSecurityContext) {
	if !HasValidSecurityContext(pod) {
		pod.Spec.SecurityContext = &sc
	}
}

// HasEnvVar returns true if environment variable exists false otherwise
func HasValidSecurityContext(pod *corev1.Pod) bool {
	// check if Pod has a security context
	if pod.Spec.SecurityContext == nil {
		return false
	}
	// check that the security context has a run as non root value
	if pod.Spec.SecurityContext.RunAsNonRoot == nil || *pod.Spec.SecurityContext.RunAsNonRoot == false {
		return false
	}
	return true
}
