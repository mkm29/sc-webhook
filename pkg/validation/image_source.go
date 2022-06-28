package validation

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

const REGISTRY_BASE_URL = "REGISTRY_BASE_URL"

// imageValidator is a container for validating the name of pods
type imageValidator struct {
	Logger logrus.FieldLogger
}

// imageValidator implements the podValidator interface
var _ podValidator = (*imageValidator)(nil)

// Name returns the name of imageValidator
func (n imageValidator) Name() string {
	return "image_source_validator"
}

// get REGISTRY_BASE_URL from environment variable
func GetRegistry() (string, bool) {
	// first see if the key is present
	val, ok := os.LookupEnv(REGISTRY_BASE_URL)
	if !ok {
		return "", false
	}
	return val, true
}

// Validate inspects the security context of a given pod and returns validation.
// The returned validation is only valid if the pod has a valid security context
// that is configured to not run as root
func (n imageValidator) Validate(pod *corev1.Pod) (validation, error) {
	v := validation{}
	// get approved registry from environment variable
	registry, ok := GetRegistry()
	if !ok {
		return v, fmt.Errorf(fmt.Sprintf("%s environment variable is not set", REGISTRY_BASE_URL))
	}
	// check if the image comes from a certain registry
	for _, container := range pod.Spec.Containers {
		if !strings.Contains(container.Image, registry) {
			v.Valid = false
			v.Reason = "Image source is not from approved registry"
			return v, nil
		}
	}
	return validation{Valid: true, Reason: "pod is from approved registry"}, nil
}
