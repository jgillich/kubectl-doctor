## Current list of anomaly checks

* core component health (etcd cluster members, scheduler, controller-manager)
* orphan endpoints (endpoints with no ipv4 attached)
* persistent-volume available & unclaimed
* persistent-volume-claim in lost state
* k8s nodes that are not in ready state
* orphan replicasets (desired number of replicas are bigger than 0 but the available replicas are 0)
* leftover replicasets (desired number of replicas and the available # of replicas are 0)
* orphan deployments (desired number of replicas are bigger than 0 but the available replicas are 0)
* lefover deployments (desired number of replicas and the available # of replicas are 0)