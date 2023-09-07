package triage

import (
	"context"
	"reflect"

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
	&ComponentUnhealthy{},
	&DeploymentIdle{},
	&DeploymentNotAvailable{},
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
