package triage

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PersistentVolumeAvailable struct{}

func (*PersistentVolumeAvailable) Severity() Severity {
	return Info
}

func (*PersistentVolumeAvailable) Description() string {
	return "PersistentVolume is not bound to a PersistentVolumeClaim."
}

func (*PersistentVolumeAvailable) Triage(ctx context.Context, cl client.Client) (anomalies []Anomaly, err error) {
	var list corev1.PersistentVolumeList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	for _, pv := range list.Items {
		if pv.Status.Phase == "Available" {
			anomalies = append(anomalies, Anomaly{NamespacedName: nn(&pv)})
		}
	}
	return
}
