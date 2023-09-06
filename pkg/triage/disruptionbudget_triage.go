package triage

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// OrphanedDisruptionBudget gets a kubernetes.Clientset and a specific namespace string
// then proceeds to search if there are orphan PodDisruptionBudgets
// the criteria is that the desired number of replicas are bigger than 0 but the healthy replicas are 0
func OrphanedDisruptionBudget(ctx context.Context, kubeCli *kubernetes.Clientset, namespace string) (*Triage, error) {
	listOfTriages := make([]string, 0)
	rs, err := kubeCli.PolicyV1().PodDisruptionBudgets(namespace).List(ctx, v1.ListOptions{})
	if client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	for _, i := range rs.Items {
		if i.Status.DesiredHealthy > 0 && i.Status.CurrentHealthy == 0 {
			listOfTriages = append(listOfTriages, i.GetName())
		}
	}
	return NewTriage("PodDisruptionBudgets", "Found orphan PodDisruptionBudget in namespace: "+namespace, listOfTriages), nil
}
