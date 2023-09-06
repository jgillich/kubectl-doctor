package triage

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const pvcLostPhase = "Lost"

// TriagePVC gets a coreclient and checks if there are any pvcs that are in lost state
func TriagePVC(ctx context.Context, coreClient coreclient.CoreV1Interface) (*Triage, error) {
	listOfTriages := make([]string, 0)
	pvcs, err := coreClient.PersistentVolumeClaims("").List(ctx, v1.ListOptions{})
	if client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	for _, i := range pvcs.Items {
		if i.Status.Phase == pvcLostPhase {
			listOfTriages = append(listOfTriages, i.GetName())
		}
	}
	return NewTriage("PVC", "Found PVC in Lost State!", listOfTriages), nil
}
