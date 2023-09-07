package triage

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func init() {
	List = append(List, &NodeNotReady{})
}

type NodeNotReady struct{}

func (*NodeNotReady) Severity() Severity {
	return Error
}

func (*NodeNotReady) Triage(ctx context.Context, cl client.Client) (anomalies []Anomaly, err error) {
	var list corev1.NodeList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	for _, node := range list.Items {
		for _, cond := range node.Status.Conditions {
			if cond.Type == corev1.NodeReady && cond.Status != "True" {
				anomalies = append(anomalies, Anomaly{NamespacedName: nn(&node), Reason: cond.Reason})
			}
		}
	}
	return
}
