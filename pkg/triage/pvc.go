package triage

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PersistentVolumeClaimLost struct{}

func (*PersistentVolumeClaimLost) Severity() Severity {
	return Error
}

func (*PersistentVolumeClaimLost) Description() string {
	return "PersistentVolumeClaim volume does not exist any longer and data on it was lost."
}

func (*PersistentVolumeClaimLost) Triage(ctx context.Context, cl client.Client) (anomalies []Anomaly, err error) {
	var list corev1.PersistentVolumeList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	for _, pvc := range list.Items {
		if pvc.Status.Phase == "Lost" {
			anomalies = append(anomalies, Anomaly{NamespacedName: nn(&pvc)})
		}
	}

	return
}
