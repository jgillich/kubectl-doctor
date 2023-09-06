package triage

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PersistentVolumeClaimLost struct{}

func (*PersistentVolumeClaimLost) Id() string {
	return "PersistentVolumeClaimLost"
}

func (*PersistentVolumeClaimLost) Severity() Severity {
	return Error
}

// TriagePVC gets a coreclient and checks if there are any pvcs that are in lost state
func (*PersistentVolumeClaimLost) Triage(ctx context.Context, cl client.Client) ([]Anomaly, error) {
	var list corev1.PersistentVolumeList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	var anomalies []Anomaly
	for _, pvc := range list.Items {
		if pvc.Status.Phase == "Lost" {
			anomalies = append(anomalies, Anomaly{NamespacedName: nn(&pvc)})
		}
	}

	return anomalies, nil
}
