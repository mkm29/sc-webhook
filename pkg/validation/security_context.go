package validation

import (
	"encoding/json"

	"github.com/mkm29/sc-webhook/pkg/utils"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// array of errors, JSON serializable
type JSONErrs []error

// securityContextValidator is a container for validating the name of pods
type SecurityContextValidator struct {
	Logger logrus.FieldLogger
}

// securityContextValidator implements the podValidator interface
var _ PodValidator = (*SecurityContextValidator)(nil)

// return a JSON representation of an errors array
func (je JSONErrs) MarshalJSON() ([]byte, error) {
	res := make([]interface{}, len(je))
	for i, e := range je {
		if _, ok := e.(json.Marshaler); ok {
			res[i] = e // e knows how to marshal itself
		} else {
			res[i] = e.Error() // Fallback to the error string
		}
	}
	return json.Marshal(res)
}

// Name returns the name of securityContextValidator
func (n SecurityContextValidator) Name() string {
	return "security_context_validator"
}

// Validate inspects the security context of a given pod and returns validation.
// The returned validation is only valid if the pod has a valid security context
// that is configured to not run as root
func (n SecurityContextValidator) Validate(pod *corev1.Pod) (validation, error) {
	hasSC := utils.HasValidSecurityContext(pod)
	if !hasSC {
		return validation{Valid: false, Reason: "pod does not have a valid security context"}, nil
	}
	return validation{Valid: true, Reason: "pod has a valid security context"}, nil
}

func HasValidSecurityContext(pod *corev1.Pod) {
	panic("unimplemented")
}
