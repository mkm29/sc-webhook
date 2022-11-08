package mutation

import (
	"encoding/json"

	"github.com/mkm29/sc-webhook/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/wI2L/jsondiff"
	corev1 "k8s.io/api/core/v1"
)

// Mutator is a container for mutation
type Mutator struct {
	Logger *logrus.Entry
}

// NewMutator returns an initialised instance of Mutator
func NewMutator(logger *logrus.Entry) *Mutator {
	return &Mutator{Logger: logger}
}

// podMutators is an interface used to group functions mutating pods
type podMutator interface {
	Mutate(*corev1.Pod) (*corev1.Pod, error)
	Name() string
}

// MutatePodPatch returns a json patch containing all the mutations needed for
// a given pod
func (m *Mutator) MutatePodPatch(pod *corev1.Pod) ([]byte, error) {
	// check Pod namespace and if in excluded namespaces, return
	// get list of namespaces to ignore
	xns := utils.GetExcludedNamespaces()
	// check if the pod is in the excluded namespaces
	for _, ns := range xns {
		if pod.ObjectMeta.Namespace == ns {
			// skip mutation
			return nil, nil
		}
	}

	var podName string
	if pod.Name != "" {
		podName = pod.Name
	} else {
		if pod.ObjectMeta.GenerateName != "" {
			podName = pod.ObjectMeta.GenerateName
		}
	}
	log := logrus.WithField("pod_name", podName)

	// list of all mutations to be applied to the pod
	mutations := []podMutator{
		injectSecurityContext{Logger: log},
	}

	mpod := pod.DeepCopy()

	// apply all mutations
	for _, m := range mutations {
		var err error
		mpod, err = m.Mutate(mpod)
		if err != nil {
			return nil, err
		}
	}

	// generate json patch
	patch, err := jsondiff.Compare(pod, mpod)
	if err != nil {
		return nil, err
	}

	patchb, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	return patchb, nil
}
