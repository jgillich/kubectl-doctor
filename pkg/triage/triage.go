package triage

import (
	"context"
	"reflect"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Severity int64

const (
	Info Severity = iota
	Warning
	Error
	Fatal
)

func (s Severity) String() string {
	switch s {
	case Info:
		return "Info"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	case Fatal:
		return "Fatal"
	default:
		return ""
	}
}

var List = []Triage{
	&CloudnativePGClusterNotReady{},
	&CloudnativePGClusterContinuousArchivingFailing{},
	&ComponentUnhealthy{},
	&DeploymentIdle{},
	&DeploymentNotAvailable{},
	&FluxHelmReleaseNotReady{},
	&FluxKustomizationNotReady{},
	&NamespaceTerminating{},
	&PodNotReady{},
	&PodWithoutOwner{},
	&PersistentVolumeAvailable{},
	&PersistentVolumeClaimLost{},
}

type Triage interface {
	Severity() Severity
	Description() string
	Triage(context.Context, client.Client) ([]Anomaly, error)
}

type Anomaly struct {
	NamespacedName types.NamespacedName
	Reason         string
}

func nn(obj client.Object) types.NamespacedName {
	return types.NamespacedName{Name: obj.GetName(), Namespace: obj.GetNamespace()}
}

func Id(t Triage) string {
	return reflect.TypeOf(t).Elem().Name()
}

func resourceExists(ctx context.Context, c client.Client, gvk schema.GroupVersionResource) (bool, error) {
	var list apiextensionsv1.CustomResourceDefinitionList
	if err := c.List(ctx, &list); err != nil {
		return false, err
	}

	for _, crd := range list.Items {
		if crd.Spec.Group != gvk.Group {
			continue
		}
		if crd.Spec.Names.Kind != gvk.Resource {
			continue
		}
		for _, version := range crd.Spec.Versions {
			if version.Name == gvk.Version {
				return true, nil
			}
		}
	}

	return false, nil
}
