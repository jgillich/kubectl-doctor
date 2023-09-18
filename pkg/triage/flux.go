package triage

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	helmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
)

type FluxHelmReleaseNotReady struct{}

func (*FluxHelmReleaseNotReady) Severity() Severity {
	return Error
}

func (*FluxHelmReleaseNotReady) Description() string {
	return "Flux HelmRelease not ready."
}

func (*FluxHelmReleaseNotReady) Triage(ctx context.Context, cl client.Client) (anomalies []Anomaly, err error) {
	if ok, err := resourceExists(ctx, cl, schema.GroupVersionResource{Group: helmv2beta1.GroupVersion.Group, Version: helmv2beta1.GroupVersion.Version, Resource: "HelmRelease"}); !ok {
		return nil, err
	}

	var list helmv2beta1.HelmReleaseList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	for _, resource := range list.Items {
		for _, cond := range resource.Status.Conditions {
			if cond.Type == "Ready" && cond.Status != "True" {
				anomalies = append(anomalies, Anomaly{NamespacedName: nn(&resource), Reason: cond.Reason})
			}
		}
	}
	return
}

type FluxKustomizationNotReady struct{}

func (*FluxKustomizationNotReady) Severity() Severity {
	return Error
}

func (*FluxKustomizationNotReady) Description() string {
	return "Flux Kustomization not ready."
}

func (*FluxKustomizationNotReady) Triage(ctx context.Context, cl client.Client) (anomalies []Anomaly, err error) {
	if ok, err := resourceExists(ctx, cl, schema.GroupVersionResource{Group: kustomizev1.GroupVersion.Group, Version: kustomizev1.GroupVersion.Version, Resource: "Kustomization"}); !ok {
		return nil, err
	}

	var list kustomizev1.KustomizationList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	for _, resource := range list.Items {
		for _, cond := range resource.Status.Conditions {
			if cond.Type == "Ready" && cond.Status != "True" {
				anomalies = append(anomalies, Anomaly{NamespacedName: nn(&resource), Reason: cond.Reason})
			}
		}
	}
	return
}

// status:
// conditions:
// - lastTransitionTime: "2023-09-18T10:16:07Z"
// 	message: Release reconciliation succeeded
// 	reason: ReconciliationSucceeded
// 	status: "True"
// 	type: Ready
