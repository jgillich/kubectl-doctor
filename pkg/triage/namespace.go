package triage

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type NamespaceTerminating struct {
}

func (*NamespaceTerminating) Severity() Severity {
	return Error
}

func (*NamespaceTerminating) Triage(ctx context.Context, cl client.Client) (anomalies []Anomaly, err error) {
	var list corev1.NamespaceList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	for _, namespace := range list.Items {
		if namespace.DeletionTimestamp != nil {
			anomalies = append(anomalies, Anomaly{NamespacedName: nn(&namespace)})
		}
	}
	return
}
