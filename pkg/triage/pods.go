package triage

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PodWithoutOwner struct{}

func (*PodWithoutOwner) Id() string {
	return "PodWithoutOwner"
}

func (*PodWithoutOwner) Severity() Severity {
	return WarningSeverity
}

func (*PodWithoutOwner) Triage(ctx context.Context, cl client.Client) ([]Anomaly, error) {
	var list corev1.PodList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	var anomalies []Anomaly
	for _, pod := range list.Items {
		if len(pod.GetOwnerReferences()) == 0 {
			anomalies = append(anomalies, Anomaly{Name: nn(&pod)})
		}
	}
	return anomalies, nil
}

type PodNotReady struct{}

func (*PodNotReady) Id() string {
	return "PodNotReady"
}

func (*PodNotReady) Severity() Severity {
	return ErrorSeverity
}

func (*PodNotReady) Triage(ctx context.Context, cl client.Client) ([]Anomaly, error) {
	var list corev1.PodList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	var anomalies []Anomaly
	for _, pod := range list.Items {
		for _, cond := range pod.Status.Conditions {
			if cond.Type == corev1.PodReady && cond.Status != "True" && cond.Reason != "PodCompleted" {
				anomalies = append(anomalies, Anomaly{Name: nn(&pod)})
			}
		}
	}
	return anomalies, nil
}
