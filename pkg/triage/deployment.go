package triage

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DeploymentNotAvailable struct{}

func (*DeploymentNotAvailable) Id() string {
	return "DeploymentNotAvailable"
}

func (*DeploymentNotAvailable) Severity() Severity {
	return Error
}

func (*DeploymentNotAvailable) Triage(ctx context.Context, cl client.Client) (anomalies []Anomaly, err error) {
	var list appsv1.DeploymentList
	if err = cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return
	}

	for _, deployment := range list.Items {
		for _, cond := range deployment.Status.Conditions {
			if cond.Type == appsv1.DeploymentAvailable && cond.Status != "True" {
				anomalies = append(anomalies, Anomaly{NamespacedName: nn(&deployment)})
			}
		}
	}
	return
}

type DeploymentIdle struct{}

func (*DeploymentIdle) Id() string {
	return "DeploymentIdle"
}

func (*DeploymentIdle) Severity() Severity {
	return Warning
}

func (*DeploymentIdle) Triage(ctx context.Context, cl client.Client) (anomalies []Anomaly, err error) {
	var list appsv1.DeploymentList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	for _, deployment := range list.Items {
		if deployment.Status.Replicas == 0 && deployment.Status.AvailableReplicas == 0 {
			anomalies = append(anomalies, Anomaly{NamespacedName: nn(&deployment)})
		}
	}
	return
}
