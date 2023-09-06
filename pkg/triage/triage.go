package triage

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Severity int64

const (
	InfoSeverity Severity = iota
	WarningSeverity
	ErrorSeverity
	FatalSeverity
)

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
	Id() string
	Severity() Severity
	Triage(context.Context, client.Client) ([]Anomaly, error)
}

type Anomaly struct {
	Name   types.NamespacedName
	Reason string
}

func nn(obj client.Object) types.NamespacedName {
	return types.NamespacedName{Name: obj.GetName(), Namespace: obj.GetNamespace()}
}
