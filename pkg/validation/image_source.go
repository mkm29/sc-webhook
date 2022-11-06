package validation

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

const REGISTRY = "REGISTRY"

// imageValidator is a container for validating the name of pods
type ImageValidator struct {
	Logger logrus.FieldLogger
}

// imageValidator implements the podValidator interface
var _ PodValidator = (*ImageValidator)(nil)

// Name returns the name of imageValidator
func (n ImageValidator) Name() string {
	return "image_source_validator"
}

// get REGISTRY from environment variable
func GetRegistry() (string, bool) {
	// first see if the key is present
	val, ok := os.LookupEnv(REGISTRY)
	if !ok {
		return "", false
	}
	return val, true
}

// Validate inspects the security context of a given pod and returns validation.
// The returned validation is only valid if the pod has a valid security context
// that is configured to not run as root
func (n ImageValidator) Validate(pod *corev1.Pod) (validation, error) {
	v := validation{}
	// get approved registry from environment variable
	registry, ok := GetRegistry()
	if !ok {
		return v, fmt.Errorf(fmt.Sprintf("%s environment variable is not set", REGISTRY))
	}
	// check if the image comes from a certain registry
	for _, container := range pod.Spec.Containers {
		if !strings.Contains(container.Image, registry) {
			v.Valid = false
			v.Reason = fmt.Sprintf("Image is not from approved registry: %s", container.Image)
			return v, nil
		}
	}
	return validation{Valid: true, Reason: "pod is from approved registry"}, nil
}
