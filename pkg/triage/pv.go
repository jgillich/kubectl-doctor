package triage

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PersistentVolumeAvailable struct{}

func (*PersistentVolumeAvailable) Id() string {
	return "PersistentVolumeAvailable"
}

func (*PersistentVolumeAvailable) Description() string {
	return "PersistentVolume is not bound to a PersistentVolumeClaim"
}

func (*PersistentVolumeAvailable) Severity() Severity {
	return InfoSeverity
}

func (*PersistentVolumeAvailable) Triage(ctx context.Context, cl client.Client) ([]Anomaly, error) {
	var list corev1.PersistentVolumeList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	var anomalies []Anomaly
	for _, pv := range list.Items {
		if pv.Status.Phase == "Available" {
			anomalies = append(anomalies, Anomaly{Name: nn(&pv)})
		}
	}
	return anomalies, nil
}
