package triage

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const componentHealthy = "True"

type Component struct{}

// TriageComponents gets a coreclient and checks if core components are in healthy state
// such as etcd cluster members, scheduler, controller-manager
func TriageComponents(ctx context.Context, coreClient coreclient.CoreV1Interface) (*Triage, error) {
	components, err := coreClient.ComponentStatuses().List(ctx, v1.ListOptions{})
	if client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	listOfTriages := make([]string, 0)
	for _, i := range components.Items {
		for _, y := range i.Conditions {
			if y.Status != componentHealthy {
				listOfTriages = append(listOfTriages, i.GetName())
			}
		}
	}
	return NewTriage("ComponentStatuses", "Found unhealthy components!", listOfTriages), nil
}
