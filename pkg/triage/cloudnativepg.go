package triage

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
)

type CloudnativePGClusterNotReady struct{}

func (*CloudnativePGClusterNotReady) Severity() Severity {
	return Error
}

func (*CloudnativePGClusterNotReady) Description() string {
	return "CloudnativePG cluster is not ready."
}

func (*CloudnativePGClusterNotReady) Triage(ctx context.Context, cl client.Client) (anomalies []Anomaly, err error) {
	if ok, err := resourceExists(ctx, cl, cnpgv1.ClusterGVK); !ok {
		return nil, err
	}

	var list cnpgv1.ClusterList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	for _, cluster := range list.Items {
		for _, cond := range cluster.Status.Conditions {
			if cond.Type == "Healthy" && cond.Status != "True" {
				anomalies = append(anomalies, Anomaly{NamespacedName: nn(&cluster), Reason: cond.Reason})
			}
		}
	}
	return
}

type CloudnativePGClusterContinuousArchivingFailing struct{}

func (*CloudnativePGClusterContinuousArchivingFailing) Severity() Severity {
	return Warning
}

func (*CloudnativePGClusterContinuousArchivingFailing) Description() string {
	return "CloudnativePG cluster continuous archiving is failing."
}

func (*CloudnativePGClusterContinuousArchivingFailing) Triage(ctx context.Context, cl client.Client) (anomalies []Anomaly, err error) {
	if ok, err := resourceExists(ctx, cl, cnpgv1.ClusterGVK); !ok {
		return nil, err
	}

	var list cnpgv1.ClusterList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	for _, cluster := range list.Items {
		for _, cond := range cluster.Status.Conditions {
			if cond.Type == "ContinuousArchiving" && cond.Status != "True" {
				anomalies = append(anomalies, Anomaly{NamespacedName: nn(&cluster), Reason: cond.Reason})
			}
		}
	}
	return
}
