package utils

import (
	corev1 "k8s.io/api/core/v1"
)

// HasEnvVar returns true if environment variable exists false otherwise
func HasValidSecurityContext(pod *corev1.Pod) bool {
	// security context can exist on 2 levels: pod and containers
	hasSc := false
	// check if Pod has a security context
	// check that the security context has a run as non root value
	// if !*pod.Spec.SecurityContext.RunAsNonRoot {
	// 	hasSc = true
	// }
	// check if containers have a security context
	for i := range pod.Spec.Containers {
		// *pod.Spec.Containers[i].SecurityContext.RunAsNonRoot && *pod.Spec.Containers[i].SecurityContext.Privileged && *pod.Spec.Containers[i].SecurityContext.AllowPrivilegeEscalation
		if pod.Spec.Containers[i].SecurityContext != nil {
			if pod.Spec.Containers[i].SecurityContext.RunAsNonRoot != nil && pod.Spec.Containers[i].SecurityContext.Privileged != nil && pod.Spec.Containers[i].SecurityContext.AllowPrivilegeEscalation != nil {
				hasSc = true
				continue
			}
		}
	}

	return hasSc
}
