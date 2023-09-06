package triage

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StandalonePods gets a kubernetes.Clientset and a specific namespace string
// then proceeds to search if there are leftover deployments
// the criteria is that a pod has no ownership (Deployment/Statefulset)
func TerminatingNamespaces(ctx context.Context, kubeCli *kubernetes.Clientset) (*Triage, error) {
	listOfTriages := make([]string, 0)
	namespaces, err := kubeCli.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	for _, i := range namespaces.Items {
		if i.Status.Phase == "Terminating" {
			listOfTriages = append(listOfTriages, i.GetName())
		}
	}
	return NewTriage("Namespaces", "Found Terminating Namespaces: ", listOfTriages), nil
}
