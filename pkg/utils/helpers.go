package utils

import (
	"os"
	"strings"

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

// get REGISTRY from environment variable
func GetEnvironmentVariable(ev string) (string, bool) {
	// first see if the key is present
	val, ok := os.LookupEnv(ev)
	if !ok {
		return "", false
	}
	return val, true
}

// Get list of excluded namespaces from environment variable
func GetExcludedNamespaces() []string {
	// get environment variable: EXCLUDED_NAMESPACES
	val, ok := GetEnvironmentVariable("EXCLUDED_NAMESPACES")
	if !ok {
		return []string{"kube-system", "kube-public", "kube-node-lease"}
	}
	// split the string by comma
	return strings.Split(val, ",")
}
