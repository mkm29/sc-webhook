package validation

import (
	"fmt"
	"strings"

	"github.com/mkm29/sc-webhook/pkg/utils"
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

func getRegistry() string {
	registry, ok := utils.GetEnvironmentVariable(REGISTRY)
	if !ok {
		return ""
	}
	return registry
}

// Validate inspects the security context of a given pod and returns validation.
// The returned validation is only valid if the pod has a valid security context
// that is configured to not run as root
func (n ImageValidator) Validate(pod *corev1.Pod) (validation, error) {
	v := validation{}
	// get list of namespaces to ignore
	xns := utils.GetExcludedNamespaces()
	// check if the pod is in the excluded namespaces
	for _, ns := range xns {
		if pod.ObjectMeta.Namespace == ns {
			return validation{Valid: true, Reason: fmt.Sprintf("pod is in protected namespace: %s", ns)}, nil
		}
	}
	// get approved registry from environment variable
	registry := getRegistry()
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
