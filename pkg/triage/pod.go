package triage

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PodWithoutOwner struct{}

func (*PodWithoutOwner) Severity() Severity {
	return Warning
}

func (*PodWithoutOwner) Description() string {
	return "Pod does not have any owner references and may be lost when evicted."
}

func (*PodWithoutOwner) Triage(ctx context.Context, cl client.Client) (anomalies []Anomaly, err error) {
	var list corev1.PodList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	for _, pod := range list.Items {
		if len(pod.GetOwnerReferences()) == 0 {
			anomalies = append(anomalies, Anomaly{NamespacedName: nn(&pod)})
		}
	}
	return
}

type PodNotReady struct{}

func (*PodNotReady) Severity() Severity {
	return Error
}

func (*PodNotReady) Description() string {
	return "Pod is not ready."
}

func (*PodNotReady) Triage(ctx context.Context, cl client.Client) (anomalies []Anomaly, err error) {
	var list corev1.PodList
	if err := cl.List(ctx, &list); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	for _, pod := range list.Items {
		for _, cond := range pod.Status.Conditions {
			if cond.Type == corev1.PodReady && cond.Status != "True" && cond.Reason != "PodCompleted" {
				reason := cond.Reason
				if cond.Reason == "ContainersNotReady" {
					reason = ""
					for _, s := range pod.Status.ContainerStatuses {
						if s.State.Waiting != nil {
							reason += fmt.Sprintf("%s(%s)", s.State.Waiting.Reason, s.Name)
							if len(s.State.Waiting.Message) > 0 {
								reason += fmt.Sprintf(": %s", s.State.Waiting.Message)
							}
							// if s.State.Waiting.Reason == "CrashLoopBackOff" {
							// 	reason +=
							// }
							break
						} else if s.State.Terminated != nil {
							reason += fmt.Sprintf("%s(%s)", s.State.Terminated.Reason, s.Name)
							if len(s.State.Terminated.Message) > 0 {
								reason += fmt.Sprintf(": %s", s.State.Terminated.Message)
							}
							break
						}
					}
				}
				anomalies = append(anomalies, Anomaly{NamespacedName: nn(&pod), Reason: reason})
			}
		}
	}
	return
}
