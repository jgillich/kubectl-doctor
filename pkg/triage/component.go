package triage

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ComponentUnhealthy struct{}

func (*ComponentUnhealthy) Id() string {
	return "ComponentUnhealthy"
}

func (*ComponentUnhealthy) Severity() Severity {
	return Fatal
}

func (*ComponentUnhealthy) Triage(ctx context.Context, cl client.Client) ([]Anomaly, error) {
	var list corev1.ComponentStatusList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	var anomalies []Anomaly
	for _, componentStatus := range list.Items {
		for _, cond := range componentStatus.Conditions {
			if cond.Status != "True" {
				anomalies = append(anomalies, Anomaly{NamespacedName: nn(&componentStatus)})
			}
		}
	}

	return anomalies, nil
}