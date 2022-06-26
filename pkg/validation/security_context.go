package validation

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// array of errors, JSON serializable
type JSONErrs []error

// securityContextValidator is a container for validating the name of pods
type securityContextValidator struct {
	Logger logrus.FieldLogger
}

// securityContextValidator implements the podValidator interface
var _ podValidator = (*securityContextValidator)(nil)

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
func (n securityContextValidator) Name() string {
	return "security_context_validator"
}

// Validate inspects the security context of a given pod and returns validation.
// The returned validation is only valid if the pod has a valid security context
// that is configured to not run as root
func (n securityContextValidator) Validate(pod *corev1.Pod) (validation, error) {
	v := validation{}
	if pod.Spec.SecurityContext == nil {
		v.Valid = false
		v.Reason = fmt.Sprintf("pod %s has no security context", pod.Name)
		return v, nil
	}
	if pod.Spec.SecurityContext.RunAsNonRoot == nil {
		v.Valid = false
		v.Reason = fmt.Sprintf("pod %s has no security context.RunAsNonRoot", pod.Name)
		return v, nil
	} else if *pod.Spec.SecurityContext.RunAsNonRoot == false {
		v.Valid = false
		v.Reason = fmt.Sprintf("pod %s has security context.RunAsNonRoot set to true", pod.Name)
		return v, nil
	}
	return validation{Valid: true, Reason: "pod has a security context"}, nil
}
