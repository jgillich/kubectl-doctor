# kubectl-doctor

Find anomalies in your Kubernetes cluster.

```
$ kubectl doctor triage
 TYPE                       SEVERITY  NAMESPACE    NAME                                       REASON                                                                                                                                                     
 PersistentVolumeAvailable      Info               pv-available                                                                                                                                                                                          
 DeploymentNotAvailable        Error  default      deployment-exit                            MinimumReplicasUnavailable                                                                                                                                 
 DeploymentNotAvailable        Error  default      deployment-invalid-image                   MinimumReplicasUnavailable                                                                                                                                 
 PodNotReady                   Error  default      deployment-exit-9bb67b49b-mvbb5            CrashLoopBackOff(nginx): back-off 5m0s restarting failed container=nginx pod=deployment-exit-9bb67b49b-mvbb5_default(89a7b5d9-d798-4624-a5a2-b5876a32b564) 
 PodNotReady                   Error  default      deployment-invalid-image-76dd98b76d-d2g82  ImagePullBackOff(nginx): Back-off pulling image "nginx:00000"                                                                                              
 PodWithoutOwner             Warning  kube-system  storage-provisioner
 ```

# License

Copyright 2023 Jakob Gillich

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
