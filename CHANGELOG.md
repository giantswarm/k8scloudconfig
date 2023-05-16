# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project's packages adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [17.0.0] - 2023-05-16

### Removed

- Remove feature gate `TTLAfterFinished` (removed in k8s 1.25).
- Remove `PodSecurityPolicy` admission plugin (removed in k8s 1.25).

### Changed

- Require k8s 1.25 and calico 1.25.
- Update calico manifests for version 1.25.

## [16.2.0] - 2023-05-05

### Changed

- Disable ETCD compaction request from apiserver.

## [16.1.0] - 2023-04-04

### Changed

- Allow customizing etcd's `--quota-backend-bytes` flag.
- Remove `--enable-v2` flag from etcd systemd unit.
- Add `--auto-compaction-mode=revision` and `--auto-compaction-retention=1` to etcd unit.
- Run etcd defragmentation every hour.

## [16.0.0] - 2023-03-31

## Changed

- **BREAKING** bump release-operator dependency to `v4`.
- Change go version to `v1.18`.

## [15.7.0] - 2023-03-06

### Changed

- Remove `--api-endpoint` flag from `k8s-api-healthz` static pod manifest as the default value 127.0.0.1 is safe to be used now.
- Enable `CronJobTimeZone` feature gate through feature flag.

## [15.6.0] - 2023-02-07

### Changed

- Enable `CronJobTimeZone` feature gate.

## [15.5.0] - 2023-01-27

### Changed

- Improve reliability of calico CRD installer job.

## [15.4.4] - 2023-01-24

### Fixed

- Add Docker token to avoid rate limits for pulling images.

## [15.4.3] - 2023-01-19

### Added

- Remove pod limit for AWS CNI if subnet prefix is enabled. 

## [15.4.2] - 2023-01-17

- Allowed the use of all seccomp profiles for components under the restricted podsecurity policy.
- Set the default seccomp profile to runtime/default under the restricted podsecurity policy.

## [15.4.1] - 2023-01-17

### Added

- Allow customers to configure all `net.*` kernel parameters in pods.


## [15.4.0] - 2023-01-11

### Changed

- Lower apiserver's cpu request to be 1/2 of the available CPUs in the VM.

## [15.3.0] - 2022-11-29

### Changed

- Label master nodes with node-role.kubernetes.io/control-plane to comply with kubeadm/CAPI.

## [15.2.0] - 2022-11-24

### Added

- Add `component` label to scheduler and controller-manager's manifests.

### Fixed

- Add missing registry mirror to containerd config.

## [15.1.1] - 2022-11-03

### Fixed

- Remove leftover from api-server manifest.

## [15.1.0] - 2022-10-24

### Changed

- Set CPU and Memory requests for Api server.

### Fixed

- Use 'MemTotal' rather than 'MemFree' to get memory limit for api server.

## [15.0.1] - 2022-10-11

### Changed

- calico-crd-installer: Tolerate `node.cloudprovider.kubernetes.io/uninitialized`.

## [15.0.0] - 2022-09-07

### Added

- Automatically set `--max-requests-inflight`, `--max-mutating-requests-inflight` and resource limits to API Server's manifest based on node size.

### Changed

- Remove dockershim parameters from kubelet systemd unit.
- Remove --address flag from scheduler's manifest.
- Remove unused ImagePullProgressDeadline setting.

## [14.6.0] - 2022-10-24

### Changed

- Allow specifying `--service-account-signing-key-file` and `--service-account-key-file` API Server flags.

### Removed

- IRSA-specific Params.

## [14.5.2] - 2022-10-11

### Changed

- calico-crd-installer: Tolerate `node.cloudprovider.kubernetes.io/uninitialized`.

## [14.5.1] - 2022-08-31

### Fixed

- Set `SystemdCgroup` to false in `containerd` config for node pools using cgroups v1.

## [14.5.0] - 2022-08-30

### Fixed

- Check hostname is up to date on azure before starting the kubelet.

## [14.4.0] - 2022-08-29

### Changed

- Enable authn/authz for scheduler/ccm to allow prometheus scraping metrics.

## [14.3.0] - 2022-08-25

### Changed

- Allow disabling in-tree cloud controller manager.

## [14.2.1] - 2022-08-24

### Removed

- Remove `aws-cni.yaml` file creation from ignition config. File is long gone from this repo and is provided by aws-operator.

## [14.2.0] - 2022-08-09

### Removed

- Remove `giantswarm` user definition from ignition template. Same user is injected by operators when generating the ignition file.

### Changed

- Use docker image instead of binary for `envsubst`.

### Fixed

- Don't fail while parsing release component versions if calico is not present in the Release.

## [14.1.2] - 2022-08-02

### Fixed

- Avoid setting `clusterCIDR` in kube proxy's config if calico is disabled.

## [14.1.1] - 2022-07-27

### Changed

- Set `api-server`'s `--kubelet-preferred-address-types` to `Hostname` for AWS to fix prometheus scraping of host network pods.

## [14.1.0] - 2022-07-25

### Changed

- Revert applying `aws-cloud-controller-manager`.

## [14.0.1] - 2022-07-15

### Fixed

- Copy `crictl.yaml` on worker nodes.

## [14.0.0] - 2022-07-14

### Changed

- **Breaking** Mandatory Kubernetes Version >=v1.23 - Change liveness probe ports and metric ports of controller-manager and scheduler.
- **Breaking** Mandatory Kubernetes Version >=v1.23 - Change apiVersion of kubescheduler to `v1beta2`.
- **Breaking** Use `containerd` socket instead of `dockershim` in the kubelet config.
- **Breaking** AWS: Use external cloud provider.
- Update `pause` container to the latest image version.

## [13.9.1] - 2022-07-07

### Fixed

- Fix syntax error in k8s-addons.

## [13.9.0] - 2022-07-04

### Added

- Override default auditd configuration to capture `execve` syscalls.

## [13.8.0] - 2022-06-07

### Fixed

- Fixed syntax error in kube-apiserver manifest introduced in 13.7.0.

### Added

- Improve k8s-kubelet unit definition to prevent nodes from joining as 'localhost'.

## [13.7.0] - 2022-05-31

### Changed

- For HA clusters, use the Load Balancer endpoint for etcd rather than localhost.

## [13.6.0] - 2022-05-23

### Changed

- Switch kubelet's `cgroupDriver` to `systemd` unless `ForceCgroupsV1` is set.

## [13.5.0] - 2022-04-19

### Added

- Add extra IRSA key args.

## [13.4.0] - 2022-04-04

### Added

- Added systemd unit to create VPA for kube-proxy.

## [13.3.0] - 2022-04-01

### Removed

- Remove VPA CR for kube-proxy.

## [13.2.0] - 2022-04-01

### Added

- Add VPA CR for kube-proxy.

## [13.1.0] - 2022-03-24

### Changed

- Leverage flatcar `3033.2.4` feature to make use of cgroups v1. 

## [13.0.0] - 2022-03-23

### Changed

- Bump apiextensions to `v6.0.0`.
- Bump go module version in `go.mod`.

## [12.1.0] - 2022-03-18

### Changed

- Increase max storage size for etcd to 8GB.

### Fixed

- Bump go module version in `go.mod`.

## [12.0.0] - 2022-03-15

### Changed

- Bump apiextensions to `v5.0.1`, use `Release` CRD from `release-operator` repository, remove `cluster-api` dependency.

## [11.1.2] - 2022-02-28

### Changed

- Allow `projected` volumes for the `restricted` psp policy.

## [11.1.1] - 2022-02-18

### Fixed

- Fix `calico-kube-controllers` permissions for `networkpolicies`.

## [11.1.0] - 2022-02-15

### Removed

- Remove `rpc-statd.service` dependency on `kubelet.service` on `AWS` and `Azure`.

## [11.0.1] - 2022-02-01

### Changed

- Bump module version to v11.

## [11.0.0] - 2022-02-01

### Added

- Add feature to force cgroup v1 on Flatcar 3033.2.0 and above. This will not work with older Flatcar version.

## [10.16.0] - 2022-01-17

### Added

- New EC2 instance types.

## [10.15.0] - 2022-01-14

### Changed

- Updated calico-policy-only manifests for calico 3.21.

## [10.14.0] - 2021-10-12

- Add net dir mount to calico-node container

## [10.13.0] - 2021-10-05
- Bump calico end etcd constraints.

## Added

- Add Calico Typha to `calico-all.yaml` to reduce load on Kubernetes API.

## [10.12.2] - 2021-09-28

## Fixed

- Apply `calico-policy-only` manifest when `aws-cni` is used.

## [10.12.1] - 2021-09-13

### Fixed

- Avoid running Calico datastore migration pods when Calico is in policy only mode.

## [10.12.0] - 2021-09-09

### Added

- Add access to networkpolicies to calico-kube-controllers service account.


## [10.11.0] - 2021-09-01

### Changed

- Update manifests for Calico 3.19 compatibility.
- Separate Calico version constraint for policy-only deployment.
- Migrate Calico datastore from etcd to Kubernetes (KVM-only).

## [10.10.0] - 2021-08-25

### Changed

- Temporarily re-enable `ssh-rsa` `CASignatureAlgorithm` in sshd until it is fully removed

## [10.9.1] - 2021-08-20

### Changed

- Add check to only include `InTreePluginAWSUnregister` flag if set to true.

## [10.9.0] - 2021-08-16

### Added

- Set `service-account-jwks-uri` api server value to align with `service-account-issuer` value

### Change

- Replace `EnableCSIMigrationAWSComplete` feature gate flag with `InTreePluginAWSUnregister`

## [10.8.1] - 2021-07-01

### Added

- Set `kernelMemcgNotification` kubelet configuration to true

## [10.8.0] - 2021-05-25

### Changed

- Enable AWS CSI migration.

## [10.7.1] - 2021-05-24

### Fixed

- Fix tags in Worker Ignition

## [10.7.0] - 2021-05-20

### Added

- Added `KVMWorkerHostVolumes` in worker template.

## [10.6.0] - 2021-05-20

### Added

- Added `enableServer` config in kubelet config file
- Propagate proxy environments to the node templates.

### Changed

- Change deprecated `--dry-run` kubectl flag to `--dry-run=client` in k8s-addons script

### Removed

- Removed deprecated kubernetes api-server flag `--kubelet-https`
- Removed deprecated kubernetes api-server flag `--insecure-port`
- Removed `--enable-server` kubelet flag since it is now deprecated

## [10.5.0] - 2021-05-12

### Added

- Added `serviceAccountIssuer`, `serviceAccountKeyFile` and `serviceAccountSigningKeyFile` flags for k8s-api-server. Required in k8s v1.20

## [10.4.0] - 2021-05-03

### Added

- Add `--bind-address-hard-fail` flag to kubeproxy to hard fail on bind failure

### Fixed

- Wait for hostname to be set before running k8s-addons.

## [10.3.0] - 2021-04-29

### Changed

- Install Calico CRDs using a separate app, `calico-crd-installer`.

## [10.2.1] - 2021-04-19

### Changed

- Retrieve Calico CRDs using HTTPS rather than base64 embedded in the ignition to fix an issue with slow startup times.

## [10.2.0] - 2021-03-19

### Changed

- Enable `anonymous-auth` in API server to comply with CAPI (needed by `kubeadm`).

## [10.1.0] - 2021-02-23

### Changed

- Move Calico (full and policy-only) CRDs into a separate file (`/srv/calico-crds.yaml`) and upgrade to CRD v1 API.
- Set `streamingConnectionIdleTimeout` to 1hr (was previously unset, default is 4h).
- Set `api-server` request timeout to 1 minute (previously unset, default is 1 minute).

### Removed

- Drop bgppeer KeepOriginalNextHop default field.

## [10.0.0] - 2020-12-10

### Removed

- Drop support for Kubernetes 1.16-1.18.
- Replace `k8s-extract-hyperkube-wrappers` and `k8s-extract-binaries` scripts with `k8s-extract`.

### Changed

- Move scheduler config from `v1alpha1` to `v1beta1`.
- Rename module from `github.com/giantswarm/k8scloudconfig/v9` to `github.com/giantswarm/k8scloudconfig/v10`.
- Trim "v" prefix from hyperkube image to match new image tag format.

## [9.3.0] - 2020-12-07

## [9.2.0] - 2020-12-01

### Changed

- Remove explicit registry pull limits defaulting to less restrictive upstream settings.

### Removed

- `runtime-cgroup` kubelet flag

## [9.1.3] - 2020-11-24

### Changed

- Shorten `calico-node` wait timeout in `k8s-addons` and add retry for faster cluster initialization.
- Synchronize `calico-node` pod template labels between `calico-all.yaml` and `calico-policy-only.yaml`.
- Remove non-functional `aws-node` wait in `k8s-addons`.
- Remove unused Kubernetes scheduler configuration fields preventing strict YAML unmarshalling.

## [9.1.2] - 2020-11-23

### Changed

- AWS - decrease `hostnetworkpod` count for calculating pod limit due `cert-exporter`  running without host network since `v1.3.0`.

## [9.1.1] - 2020-11-09

### Added

- Set configurable labels with `k8s-label-node` unit as well to update old labels when node identity doesn't change on upgrade.

## [9.1.0] - 2020-10-29

### Added

- Add dockerhub authentication for `kubelet`.
- Use `root` user explicitly for `kubelet` systemd unit.

## [9.0.0] - 2020-10-27

### Changed

- Enable `kubelet` flag to protect kernel defaults
- Set `scheduler` address to local address `127.0.0.1`
- Update apiextensions to v3 and replace CAPI with Giant Swarm fork.
- Prepare module v9.
- (KVM only) Update Calico etcd certs and hostPath mounts corresponding to changes in v8.0.4.

### Removed

- Disable `kubelet` read-only port

## [8.0.4] - 2020-10-21

### Changed

- Updated certificates used by kubelet for Calico etcd datastore to match new location in certs@v3 library.

## [8.0.3] - 2020-10-05

### Added

- Add timer to run `k8s-label-node.service` every hour to ensure core labels are present.

## [8.0.2] - 2020-09-30

### Fixed

- Removed extra line break when setting k8s api server arguments.

### Changed

- Allow parallel download of images in kubelet templates

## [8.0.1] - 2020-09-17

### Removed

- When calico is used only for Network Policies it will not install the CNI binaries. The CNI in each provider will take care of installing the required binaries.

### Added

- Add monitoring annotations `prometheus.io/*` and `giantswarm.io/monitoring*` to kube-proxy, k8s-scheduler, k8s-controller-manager and calico.

### Changed

- Changed the path of the ETCD certificate files used in the etcdctl alias.
- Exposed some of the etcd3.service systemd unit settings via environment variables to make customizations in the configuration easier.

## [8.0.0] - 2020-08-11

### Added

- Added validation of versions in the cloud config Params struct. Versions outside of supported ranges will cause an
  error to be returned from cloud config-related functions.

### Changed

- Updated backward incompatible Kubernetes dependencies to v1.18.5.

### Removed

- Removed `DefaultParams` and `DefaultCloudConfigConfig` functions from the `template` package. Defaults should be
  established by the consumer of the library instead.

## [7.0.5] - 2020-07-30

### Changed

- Adjusted number of host network pods on worker node for aws-cni.

## [7.0.4] - 2020-07-29

### Changed

- Adjusted `MAX_PODS` for master and worker nodes to max IP's per ENI when using AWS CNI.

## [7.0.3] - 2020-07-23

### Fixed

- Set etcd data dir permission to `0700` to comply with etcd 3.10.4 requirements.

## [7.0.2] - 2020-07-22

### Removed
- Removed PV node limits for AWS as the feature gate is no longer supported in 1.17+.

## [7.0.1] - 2020-07-20

### Changed

- Changed wrong v6 reference to use the latest v7 module.

## [7.0.0] - 2020-07-07

### Added

- Add `Params.RegistryMirrors` allowing to configure docker registry mirrors.

### Changed

- Fail if all images do not have the same registry.

## [6.4.0] - 2020-07-06

### Added

- Add back registry domain configuration.

### Changed

- Change `kube-apiserver` image to include certs.
- Delete kube-proxy and Calico DaemonSets/Deployments with `--cascade=false`
  when upgrading from clusters using k8scloudconfig v6.0.3/v6.1.0 so that
  upgrades can continue without manual intervention.

### Removed

- Remove `RegistryDomain` template parameter.

## [6.3.0] - 2020-06-22

## Changed

- Upgrade calico to 3.14.1
- Slightly changed the configuration interface for Calico

## [6.2.6] - 2020-06-19

### Removed

- Remove quay.io as registry mirror over security concerns

## [6.2.5] - 2020-06-18

### Added

- Add quay.io as docker.io mirror in dockerd config

### Changed

- Use templated registry domain value for docker registry mirror
- Use `giantswarm/pause:3.1` container as pod infra container instead of default container, hosted on gcr.
- Bind kubelet health check endpoint to all IPv4 addresses.

### Removed

- Remove registry domain availability as we have no failover

## [6.2.4] 2020-06-16

### Changed

- Drop `clusterCIDR` from `kube-proxy` config on Azure.

### Fixed

- Fix worker's `$IMAGE` k8s-setup-network-env systemd unit to pick up value
  from `.Images` instead of `.Cluster`.

## [6.2.3] 2020-06-09

### Changed

- Enable felix metrics for calico policy-only manifest.

### Removed

- Remove typha deployment for calico policy-only manifest.

## [6.2.2] 2020-06-04

### Added

- Add option to set etcd initial state to 'new' or 'existing'.

### Changed

- Explicitly set TLS cipher suites.
- Specify zap logger for etcd as capnslog is deprecated in v3.4.

### Fixed

- Fix `rpc-statd.service` not running before kubelet.
- Fix regression in kubelet installation systemd unit for 1.16 clusters.
- Fix runtime cgroups configuration for kubelet.

### Removed

- Remove `resourceContainer` from `kube-proxy` configuration file.

## [6.2.1] 2020-05-26

### Added

- Add `bird-live` flag to calico node liveness probe.
- RBAC permissions allowing calico node to get configmaps.
- Parameter to disable deletion of master nodes for HA masters.

## [6.2.0] 2020-05-20

### Added

- Support for highly available etcd clusters.
- Kubernetes 1.17 compatibility.

### Removed

- Remove limits from calico-kube-controllers.

## [6.1.1] 2020-05-07

### Fixed

- Revert changes to deployment label selectors causing k8s-addons to fail.

## [6.1.0] 2020-05-06

### Changed

- Fix conntrack configuration structure for `kube-proxy`.
- Flatten directory structure. Only the most recent version lives in this repo now.
  Go module version becomes synonymous with cloudconfig version.

### Removed

- All cloudconfig versions prior to v6.0.0.
- Remove IC performance improvements from OS provisioning.

## [6.0.3] 2020-04-16

### Added

- A new template variable `EnableAWSCNI` which should be set to `true` to get AWS CNI specific files/config.
- Disable profiling for Controller Manager and Scheduler.


## [6.0.2]

### Changed

- Remove init limits from calico-node
- Limit PV per node on AWS
- Hardcode registry domain AWS


## [6.0.1]

### Fixed

- Fix go module.



## [6.0.0]

### Changed

- Extract images and versions out from k8scloudconfig and make them templatable by importer.
- Switch from dep to go modules.
- Use architect orb.
- Add persistent volume node limit for AWS.

### Added

- Add `conntrackMaxPerCore` parameter in `kube-proxy` manifest.

## [5.2.0]

### Changed

- Reserve ports `30000-32767` from ephemeral port range for `kube-apiserver` use.
- Make provisioning idempotent by generating `/boot/coreos/first_boot` file on every boot.
- Use [AWS VPC CNI](https://github.com/aws/amazon-vpc-cni-k8s) for pod networking and Calico for ensuring network policies.
- Enable ':9393/metrics' prometheus endpoint in docker daemon.

## [5.1.1] - Unreleased

### Changed

- Update Kubernetes to `1.16.7`.

## [5.1.0] - 2020-01-21

### Changed

- Lowercase $(hostname) to match k8s node name e.g. when using with kubectl.
- Extend ignition with debug options.

## [5.0.0] - 2020-01-02

### Changed

- Moved kubelet from container to host process (`--containerized` flag is removed in Kubernetes 1.16).
- Changed `restricted` PodSecurityPolicy to restrict the allowed range of user IDs for PODs.
- Update Kubernetes to `1.16.3`.
- Update Calico to `3.10.1` along with corresponding RBAC rules.
- Update etcd to `3.3.17`.
- Update `calicoctl` (debug tool) to `3.10.1`.
- Update `crictl` (debug tool) to `1.16.1`.
- Clean up k8s-addons (use system `kubectl`, avoid `kubectl get cs`).
- Apply kubelet restricted role labels using new systemd service.
- Increase `fs.inotify.max_user_instances` to 8192.
- Change Priority Class for `calico-node` to `system-node-critical`.
- Use registry domain for k8s-api-healthz and wait for domains script for AWS China.

### Added

- Add eviction hard setting for image file system in kubelet.
- Add Deny All as default Network Policy in `kube-system` and `giantswarm` namespaces.

## [4.9.2] - 2020-04-15

### Changed

- Remove debug profiling from Controller Manager and Scheduler

## [4.9.1] - 2020-03-10

### Added

- Add `conntrackMaxPerCore` parameter in `kube-proxy` manifest.

### Changed

-  Remove limit of calico node init container.

## [4.9.0] - 2019-10-17

### Changed

- Bind kube-proxy metrics address to 0.0.0.0 instead of default 127.0.0.1 in config file.
- Remove Calico Node limits.
- Update Kubernetes to `1.15.5` (including CVE-2019-11251).
- Update Calico to `3.9.1`.
- Update etcd to `3.3.15`.
- Update `calicoctl` (debug tool) to `3.9.1`.
- Update `crictl` (debug tool) to `1.15.0`.
- Change `calico-node` `DaemonSet` from `v1beta1` to `v1`.
- Change `calico-kube-controllers` `Deployment` from `v1beta1` to `v1`.
- Use `/bin/calico-node -felix-live` for `calico-node` liveness probe instead of `httpGet`.
- Generally minimize differences between [Calico v3.9 yaml](https://docs.projectcalico.org/v3.9/manifests/calico.yaml) and `calico-all.yaml`.

## [4.8.1] - 2019-12-31

### Changed

- Update Kubernetes to 1.14.10, includes fixes for CVE-2019-11253 and some Azure fixes.
- Increase `fs.inotify.max_user_instances` to 8192.

## [4.8.0]

### Added

- Add k8s-api-healthz service to master node to enable proper LB health checks to api and etcd.
- Set api-server listen address to 0.0.0.0.

## [4.7.0]

### Added

- Enable TTLAfterFinished feature gate. This allows a TTL controller to clean up jobs after they finish execution.

### Changed

- Update kubernetes to 1.14.6, includes fixes for CVE-2019-9512, CVE-2019-9514
- Update calico to 3.8.2

## [4.6.0]

### Added

- Systemd unit, which sets certificates group owner to `giantswarm`, so that cert-exporter running as user `giantswarm` is able to read certificates.

### Changed

- Mount relevant directories so that the command `docker` can run in `Kubelet`. This is needed for `rbd` to mount `Ceph` volumes on the nodes.
- Add explicit cgroups for finer grained resources management over operating system components and container runtime.
- Make --image-pull-progress-deadline configurable for kubelets so a longer
duration can be used in AWS China regions to mitigate slow image pulls.
- Harden `restricted` podsecuritypolicy.

### Fixed

- Update `giantswarm-critical` priority class manifest to use `v1` stable.

## [4.5.1]

### Changed

- Update kubernetes to 1.14.5 CVE-2019-1002101, CVE-2019-11246

## [4.5.0]

### Changed

- Add configuration necessery for generic support of rbd storage.
- Add `name` label for `kube-system` and `default` namespaces.

## [4.4.0]

### Changed

- Change Felix configuration to add metric server and expose data to be scraped for prometheus.
- Add `k8s-app` label for `api-server`, `controller-manager` and `scheduler`.
- Harden SSH config and tuned networking kernel settings
- Update kubernetes to 1.14.3
- Update calico to 3.7.1
- Update etcd to 3.3.13.


## [4.3.0]

### Changed
- Update kubernetes to 1.14.1
- Update calico to 3.6.1
- Change node labels for master and workers
- Update kube-proxy and calico to tolerate every taint effects and CriticalAddonsOnly
- Add managed giantswarm label to calico daemonset

## [4.2.0]

### Changed
- Fix race condition issue with systemd units.

### Removed

- Remove `UsePrivilegeSeparation` option from sshd configuration.

## [4.1.2]
### Changed
- Pin calico-kube-controllers to master.
- Fix calico-node felix severity log level.
- Enable `serviceaccount` controller in calico-kube-controller.
- Remove 'staticPodPath' from worker kubelet configuration.

## [4.1.1]
### Changed
- Update kubernetes to 1.13.4 CVE-2019-1002100

## [4.1.0]
### Changed
- Intall calicoctl, crictl and configure etcctl tooling in masters.
- Update kubernetes to 1.13.3.
- Update etcd to 3.3.12.
- Update calico to 3.5.1.
- Add fine-grained Audit Policy

## [4.0.1]
### Changed
- Update kubernetes to 1.12.6 CVE-2019-1002100

## [3.8.0] WIP
- Update kubernetes to 1.12.6 CVE-2019-1002100

## [4.0.0]

### Changed
- Switched from cloudinit to ignition.
- Double the inotify watches.
- Switch kube-proxy from `iptables` to `ipvs`.
- Add worker node labels.
- Increase timeouts for etcd defragmentaion.

### Removed

- Ingress Controller and CoreDNS manifests. Now migrated to chart-operator.
- Removed nodename_file_optional from calico configmap.

## [3.7.5]
- Update kubernetes to 1.12.6 CVE-2019-1002100

## [3.7.4]

### Changed
- Double the inotify watches.

### Removed
- Removed nodename_file_optional from calico configmap.

## [3.7.3]

### Changed
- update kubernetes to 1.12.3 (CVE-2018-1002105)

## [3.6.4]

### Changed
- Update `libreadline` version

## [3.6.3]
- update kubernetes to 1.11.5 (CVE-2018-1002105)

### Changed
- update kubernetes to 1.10.11 (CVE-2018-1002105)

## [3.5.3]

### Changed
- Update `libreadline` version

## [3.5.2]

### Changed

## [3.7.2]

### Changed
- Remove the old master from the k8s api before upgrading calico (k8s-addons)
- Wait until etcd DNS is resolvable before upgrading calico. Networking pods crashlooping isn't fun!

## [3.7.1]

### Changed
- The pod priority class for calico got lost. We found it again!
- kube-proxy is now installed before calico during cluster creation and upgrades.

## [3.7.0]

### Changed
- Updated Kubernetes to 1.12.2
- Updated etcd to 3.3.9
- Kubernetes and etcd images are now held in one place
- Updated audit policy version
- Moved audit policy out of static pod path
- Updated rbac resources to v1
- Remove static pod path from worker nodes
- Remove readonly port from kubelet
- Add DBUS socket and ClusterCIDR to kube-proxy

## [3.6.2]

### Changed
- Updated Calico to 3.2.3
- Updated Calico manifest with resource limits to get QoS policy guaranteed.
- Enabled admission plugins: DefaultTolerationSeconds, MutatingAdmissionWebhook, ValidatingAdmissionWebhook.

## [3.6.1]

### Changed
- Use patched GiantSwarm build of Kubernetes (`hyperkube:v1.11.1-cec4fb8023db783fbf26fb056bf6c76abfcd96cf-giantswarm`).

## [3.6.0]

### Added
- Added template flag for removing CoreDNS resources (will be managed by
chart-operator).

### Changed
- Updated Kubernetes (hyperkube) to version 1.11.1.
- Updated Calico to version 3.2.0.

### Removed


## [3.5.1]


## [3.5.0]

### Changed
- Disabled HSTS headers in nginx-ingress-controller.
- Added separate parameter for disabling the Ingress Controller service manifest.

### Removed


## [3.4.0]

### Added
- Added SSO public key into ssh trusted CA.
- Added RBAC rules for node-operator.
- Added RBAC rules for prometheus.
- Enabled monitoring for ingress controller metrics.
- Parameterize image registry domain.
- Set "worker-processes" to 4 for nginx-ingress-controller.
- Added `--feature-gates=CustomResourceSubresources=true` to enable status subresources for CRDs.

### Changed
- Add memory eviction thresholds for kubelet in order to preserve node in case of heavy load.
- Updated etcd version to 3.3.8

### Removed


## [3.3.4]

### Changed
- Added parameter for disabling Ingress Controller related components.
- Increased start timeout for k8s-kubelet.service.

### Removed


## [3.3.3]

### Changed

- Remove Nginx version from `Server` header in Ingress Controller
- Updated hyperkube to version 1.10.4.

### Removed


## [3.3.2]

### Changed
- Updated hyperkube to version 1.10.2 due to regression in 1.10.3 with configmaps.

### Removed
- Removed node-exporter related components (will be managed by chart-operator).

## [3.3.1]

### Changed
- Changed some remaining images to be pulled from Giant Swarm's registry.
- Updated Alpine sidecar for Ingress Controller to version 3.7.
- Fixed mkfs.xfs for containerized kubelet.
- Updated Kubernetes (hyperkube) to version 1.10.3.

### Removed


## [3.3.0]

### Changed
- Updated hyperkube to version 1.10.2.

### Removed
- Removed kube-state-metrics related components (will be managed by
chart-operator).


## [3.2.6]

### Changed
- Changed node-exporter to have named ports.
- Added RBAC rules for configmaps, secrets and hpa for kube-state-metrics.
- Synced privileged PSP with upstream (adding all added capabilities and seccomp profiles)
- Downgraded hyperkube to version 1.9.5.

### Removed


## [3.2.5]

### Changed
- Updated kube-state-metrics to version 1.3.1.
- Updated hyperkube to version 1.10.1.
- Changed kubelet bind mount mode from "shared" to "rshared".
- Disabled etcd3-defragmentation service in favor systemd timer.
- Added /lib/modules mount for kubelet.
- Updated CoreDNS to 1.1.1.
- Fixed node-exporter running in container by adjusting host mounts.
- Updated Calico to 3.0.5.
- Updated Etcd to 3.3.3.
- Added trusted certificate CNs to aggregation API allowed names.
- Disabled SSL passthrough in nginx-ingress-controller.
- Removed explicit hostname labeling for kubelet.

### Removed
- Removed docker flag "--disable-legacy-registry".
- Removed calico-ipip-pinger.


## [3.2.4]

### Changed
- Masked systemd-networkd-wait-online unit.
- Makes injection of Kubernetes encryption key configurable.
- Enabled volume resizing feature.



## [3.2.3]

### Changed
- Updated Kubernetes with version 1.9.5.
- Updated nginx-ingress-controller to version 0.12.0.

### Removed
- Removed hard limits from core kubernetes components.



## [3.2.2]

### Removed
- Removed set-ownership-etcd-data-dir.service.



## [3.2.1]

### Added
- Added priority classes core-components, critical-pods and important pods.
- Added Guaranteed QoS for api/scheduler/controller-manager pods by assigning resources limits.

### Changed
- Enabled aggregation layer in Kubernetes API server.
- Ordered Kubernetes cluster components scheduling process by assigning PriorityClass to pods.

## [3.1.1]

### Added
- Added calico-ipip-pinger.

### Changed
- Change etcd data path to /var/lib/etcd.
- Fix `StartLimitIntervalSec` parameter location in `etcd3` systemd unit.
- Add `feature-gates` flag in api server enabling `ExpandPersistentVolumes` feature.
- Updated calico to 3.0.2.
- Updated etcd to 3.3.1.
- Tune kubelet flags for protecting key units (kubelet and container runtime) from workload overloads.
- Updated nginx-ingress-controller to 0.11.0.
- Updated coredns to 1.0.6.

## [3.1.0]

### Changed
- Systemd units for Kubernetes components (api-server, scheduler and controller-manager)
  replaced with self-hosted pods.

## [3.0.0]

### Added
- Add encryption config template for API etcd data encryption experimental
  feature.
- Add tests for template evaluations.
- Allow adding extra manifests.
- Allow extract Hyperkube options.
- Allow setting custom K8s API address for master nodes.
- Allow setting etcd port.
- Add node-exporter.
- Add kube-state-metrics.

### Changed
- Unify CloudConfig struct construction.
- Update calico to 3.0.1.
- Update hyperkube to v1.9.2.
- Update nginx-ingress-controller to 0.10.2.
- Use vanilla (previously coreos) hyperkube image.
- kube-dns replaced with CoreDNS 1.0.5.
- Fix Kubernetes API audit log.

### Removed
- Remove calico-ipip-pinger.
- Remove calico-node-controller.

## [2.0.2]

### Added
- Add fix for scaled workers to ensure they have a kube-proxy.

## [2.0.1]

### Changed
- Fix audit logging.

## [2.0.0]

### Added
- Disable API etcd data encryption experimental feature.

### Changed
- Updated calico to 2.6.5.

### Removed
- Removed calico-ipip-pinger.
- Removed calico-node-controller.

## [1.1.0]

### Added
- Use Cluster type from https://github.com/giantswarm/apiextensions.

## [1.0.0]

### Removed
- Disable API etcd data encryption experimental feature.

## [0.1.0]

[Unreleased]: https://github.com/giantswarm/k8scloudconfig/compare/v17.0.0...HEAD
[17.0.0]: https://github.com/giantswarm/k8scloudconfig/compare/v16.2.0...v17.0.0
[16.2.0]: https://github.com/giantswarm/k8scloudconfig/compare/v16.1.0...v16.2.0
[16.1.0]: https://github.com/giantswarm/k8scloudconfig/compare/v16.0.0...v16.1.0
[16.0.0]: https://github.com/giantswarm/k8scloudconfig/compare/v15.7.0...v16.0.0
[15.7.0]: https://github.com/giantswarm/k8scloudconfig/compare/v15.6.0...v15.7.0
[15.6.0]: https://github.com/giantswarm/k8scloudconfig/compare/v15.5.0...v15.6.0
[15.5.0]: https://github.com/giantswarm/k8scloudconfig/compare/v15.4.4...v15.5.0
[15.4.4]: https://github.com/giantswarm/k8scloudconfig/compare/v15.4.3...v15.4.4
[15.4.3]: https://github.com/giantswarm/k8scloudconfig/compare/v15.4.2...v15.4.3
[15.4.2]: https://github.com/giantswarm/k8scloudconfig/compare/v15.4.1...v15.4.2
[15.4.1]: https://github.com/giantswarm/k8scloudconfig/compare/v15.4.1...v15.4.1
[15.4.1]: https://github.com/giantswarm/k8scloudconfig/compare/v15.4.0...v15.4.1
[15.4.0]: https://github.com/giantswarm/k8scloudconfig/compare/v15.3.0...v15.4.0
[15.3.0]: https://github.com/giantswarm/k8scloudconfig/compare/v15.2.0...v15.3.0
[15.2.0]: https://github.com/giantswarm/k8scloudconfig/compare/v15.1.1...v15.2.0
[15.1.1]: https://github.com/giantswarm/k8scloudconfig/compare/v15.1.0...v15.1.1
[15.1.0]: https://github.com/giantswarm/k8scloudconfig/compare/v15.0.1...v15.1.0
[15.0.1]: https://github.com/giantswarm/k8scloudconfig/compare/v15.0.0...v15.0.1
[15.0.0]: https://github.com/giantswarm/k8scloudconfig/compare/v14.6.0...v15.0.0
[14.6.0]: https://github.com/giantswarm/k8scloudconfig/compare/v14.5.2...v14.6.0
[14.5.2]: https://github.com/giantswarm/k8scloudconfig/compare/v14.5.1...v14.5.2
[14.5.1]: https://github.com/giantswarm/k8scloudconfig/compare/v14.5.0...v14.5.1
[14.5.0]: https://github.com/giantswarm/k8scloudconfig/compare/v14.4.0...v14.5.0
[14.4.0]: https://github.com/giantswarm/k8scloudconfig/compare/v14.3.0...v14.4.0
[14.3.0]: https://github.com/giantswarm/k8scloudconfig/compare/v14.2.1...v14.3.0
[14.2.1]: https://github.com/giantswarm/k8scloudconfig/compare/v14.2.0...v14.2.1
[14.2.0]: https://github.com/giantswarm/k8scloudconfig/compare/v14.1.2...v14.2.0
[14.1.2]: https://github.com/giantswarm/k8scloudconfig/compare/v14.1.1...v14.1.2
[14.1.1]: https://github.com/giantswarm/k8scloudconfig/compare/v14.1.0...v14.1.1
[14.1.0]: https://github.com/giantswarm/k8scloudconfig/compare/v14.0.1...v14.1.0
[14.0.1]: https://github.com/giantswarm/k8scloudconfig/compare/v14.0.0...v14.0.1
[14.0.0]: https://github.com/giantswarm/k8scloudconfig/compare/v13.9.1...v14.0.0
[13.9.1]: https://github.com/giantswarm/k8scloudconfig/compare/v13.9.0...v13.9.1
[13.9.0]: https://github.com/giantswarm/k8scloudconfig/compare/v13.8.0...v13.9.0
[13.8.0]: https://github.com/giantswarm/k8scloudconfig/compare/v13.7.0...v13.8.0
[13.7.0]: https://github.com/giantswarm/k8scloudconfig/compare/v13.6.0...v13.7.0
[13.6.0]: https://github.com/giantswarm/k8scloudconfig/compare/v13.5.0...v13.6.0
[13.5.0]: https://github.com/giantswarm/k8scloudconfig/compare/v13.4.0...v13.5.0
[13.4.0]: https://github.com/giantswarm/k8scloudconfig/compare/v13.3.0...v13.4.0
[13.3.0]: https://github.com/giantswarm/k8scloudconfig/compare/v13.2.0...v13.3.0
[13.2.0]: https://github.com/giantswarm/k8scloudconfig/compare/v13.1.0...v13.2.0
[13.1.0]: https://github.com/giantswarm/k8scloudconfig/compare/v13.0.0...v13.1.0
[13.0.0]: https://github.com/giantswarm/k8scloudconfig/compare/v12.1.0...v13.0.0
[12.1.0]: https://github.com/giantswarm/k8scloudconfig/compare/v12.0.0...v12.1.0
[12.0.0]: https://github.com/giantswarm/k8scloudconfig/compare/v11.1.2...v12.0.0
[11.1.2]: https://github.com/giantswarm/k8scloudconfig/compare/v11.1.1...v11.1.2
[11.1.1]: https://github.com/giantswarm/k8scloudconfig/compare/v11.1.0...v11.1.1
[11.1.0]: https://github.com/giantswarm/k8scloudconfig/compare/v11.0.1...v11.1.0
[11.0.1]: https://github.com/giantswarm/k8scloudconfig/compare/v11.0.0...v11.0.1
[11.0.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.16.0...v11.0.0
[10.16.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.15.0...v10.16.0
[10.15.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.14.0...v10.15.0
[10.14.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.13.0...v10.14.0
[10.13.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.12.2...v10.13.0
[10.12.2]: https://github.com/giantswarm/k8scloudconfig/compare/v10.12.1...v10.12.2
[10.12.1]: https://github.com/giantswarm/k8scloudconfig/compare/v10.12.0...v10.12.1
[10.12.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.11.0...v10.12.0
[10.11.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.10.0...v10.11.0
[10.10.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.9.1...v10.10.0
[10.9.1]: https://github.com/giantswarm/k8scloudconfig/compare/v10.9.0...v10.9.1
[10.9.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.8.1...v10.9.0
[10.8.1]: https://github.com/giantswarm/k8scloudconfig/compare/v10.8.0...v10.8.1
[10.8.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.7.1...v10.8.0
[10.7.1]: https://github.com/giantswarm/k8scloudconfig/compare/v10.7.0...v10.7.1
[10.7.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.6.0...v10.7.0
[10.6.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.5.0...v10.6.0
[10.5.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.4.0...v10.5.0
[10.4.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.3.0...v10.4.0
[10.3.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.2.1...v10.3.0
[10.2.1]: https://github.com/giantswarm/k8scloudconfig/compare/v10.2.0...v10.2.1
[10.2.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.1.0...v10.2.0
[10.1.0]: https://github.com/giantswarm/k8scloudconfig/compare/v10.0.0...v10.1.0
[10.0.0]: https://github.com/giantswarm/k8scloudconfig/compare/v9.3.0...v10.0.0
[9.3.0]: https://github.com/giantswarm/k8scloudconfig/compare/v9.2.0...v9.3.0
[9.2.0]: https://github.com/giantswarm/k8scloudconfig/compare/v9.1.3...v9.2.0
[9.1.3]: https://github.com/giantswarm/k8scloudconfig/compare/v9.1.2...v9.1.3
[9.1.2]: https://github.com/giantswarm/k8scloudconfig/compare/v9.1.1...v9.1.2
[9.1.1]: https://github.com/giantswarm/k8scloudconfig/compare/v9.1.0...v9.1.1
[9.1.0]: https://github.com/giantswarm/k8scloudconfig/compare/v9.0.0...v9.1.0
[9.0.0]: https://github.com/giantswarm/k8scloudconfig/compare/v8.0.4...v9.0.0
[8.0.4]: https://github.com/giantswarm/k8scloudconfig/compare/v8.0.3...v8.0.4
[8.0.3]: https://github.com/giantswarm/k8scloudconfig/compare/v8.0.2...v8.0.3
[8.0.2]: https://github.com/giantswarm/k8scloudconfig/compare/v8.0.1...v8.0.2
[8.0.1]: https://github.com/giantswarm/k8scloudconfig/compare/v8.0.0...v8.0.1
[8.0.0]: https://github.com/giantswarm/k8scloudconfig/compare/v7.0.5...v8.0.0
[7.0.5]: https://github.com/giantswarm/k8scloudconfig/compare/v7.0.4...v7.0.5
[7.0.4]: https://github.com/giantswarm/k8scloudconfig/compare/v7.0.3...v7.0.4
[7.0.3]: https://github.com/giantswarm/k8scloudconfig/compare/v7.0.2...v7.0.3
[7.0.2]: https://github.com/giantswarm/k8scloudconfig/compare/v7.0.1...v7.0.2
[7.0.1]: https://github.com/giantswarm/k8scloudconfig/compare/v7.0.0...v7.0.1
[7.0.0]: https://github.com/giantswarm/k8scloudconfig/compare/v6.4.0...v7.0.0
[6.4.0]: https://github.com/giantswarm/k8scloudconfig/compare/v6.3.0...v6.4.0
[6.3.0]: https://github.com/giantswarm/k8scloudconfig/compare/v6.2.6...v6.3.0
[6.2.6]: https://github.com/giantswarm/k8scloudconfig/compare/v6.2.5...v6.2.6
[6.2.5]: https://github.com/giantswarm/k8scloudconfig/compare/v6.2.4...v6.2.5
[6.2.4]: https://github.com/giantswarm/k8scloudconfig/compare/v6.2.3...v6.2.4
[6.2.3]: https://github.com/giantswarm/k8scloudconfig/compare/v6.2.2...v6.2.3
[6.2.2]: https://github.com/giantswarm/k8scloudconfig/compare/v6.2.1...v6.2.2
[6.2.1]: https://github.com/giantswarm/k8scloudconfig/compare/v6.2.0...v6.2.1
[6.2.0]: https://github.com/giantswarm/k8scloudconfig/compare/v6.1.1...v6.2.0
[6.1.1]: https://github.com/giantswarm/k8scloudconfig/compare/v6.1.0...v6.1.1
[6.1.0]: https://github.com/giantswarm/k8scloudconfig/compare/v6.0.3...v6.1.0
[6.0.3]: https://github.com/giantswarm/k8scloudconfig/compare/v6.0.2...v6.0.3
[6.0.2]: https://github.com/giantswarm/k8scloudconfig/compare/v6.0.1...v6.0.2
[6.0.1]: https://github.com/giantswarm/k8scloudconfig/compare/v6.0.0...v6.0.1
[6.0.0]: https://github.com/giantswarm/k8scloudconfig/releases/tag/v6.0.0
[5.2.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_5_2_0
[5.1.1]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_5_1_1
[5.1.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_5_1_0
[5.0.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_5_0_0
[4.9.2]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_9_2
[4.9.1]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_9_1
[4.9.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_9_0
[4.8.1]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_8_1
[4.8.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_8_0
[4.7.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_7_0
[4.6.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_6_0
[4.5.1]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_5_1
[4.5.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_5_0
[4.4.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_4_0
[4.3.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_3_0
[4.2.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_2_0
[4.1.2]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_1_2
[4.1.1]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_1_1
[4.1.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_1_0
[4.0.1]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_0_1
[4.0.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_4_0_0
[3.8.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_8_0
[3.7.5]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_7_5
[3.7.4]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_7_4
[3.7.3]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_7_3
[3.6.4]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_6_4
[3.6.3]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_6_3
[3.5.2]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_5_2
[3.7.2]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_7_2
[3.7.1]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_7_1
[3.7.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_7_0
[3.6.2]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_6_2
[3.6.1]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_6_1
[3.6.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_6_0
[3.5.3]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_5_3
[3.5.1]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_5_1
[3.5.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_5_0
[3.4.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_4_0
[3.3.4]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_3_4
[3.3.3]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_3_3
[3.3.2]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_3_2
[3.3.1]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_3_1
[3.3.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_3_0
[3.2.6]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_2_6
[3.2.5]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_2_5
[3.2.4]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_2_4
[3.2.3]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_2_3
[3.2.2]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_2_2
[3.2.1]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_2_1
[3.1.1]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_1_1
[3.1.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_1_0
[3.0.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_3_0_0
[2.0.2]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_2_0_2
[2.0.1]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_2_0_1
[2.0.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v2
[1.1.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v1_1
[1.0.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v1
[0.1.0]: https://github.com/giantswarm/k8scloudconfig/commits/master/v_0_1_0
