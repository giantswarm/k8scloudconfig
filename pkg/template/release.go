package template

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/giantswarm/apiextensions/pkg/apis/release/v1alpha1"
)

func BuildImages(registryDomain string, versions Versions) Images {
	return Images{
		CalicoCNI:                    buildImage(registryDomain, "giantswarm/cni", versions.Calico),
		CalicoKubeControllers:        buildImage(registryDomain, "giantswarm/kube-controllers", versions.Calico),
		CalicoNode:                   buildImage(registryDomain, "giantswarm/node", versions.Calico),
		Etcd:                         buildImage(registryDomain, "giantswarm/etcd", versions.Etcd),
		KubeApiserver:                buildImage(registryDomain, "giantswarm/kube-apiserver", versions.Kubernetes),
		KubeControllerManager:        buildImage(registryDomain, "giantswarm/kube-controller-manager", versions.Kubernetes),
		KubeScheduler:                buildImage(registryDomain, "giantswarm/kube-scheduler", versions.Kubernetes),
		Kubectl:                      buildImage(registryDomain, "giantswarm/docker-kubectl", versions.Kubectl),
		Kubelet:                      buildImage(registryDomain, "giantswarm/kubelet", versions.Kubernetes),
		KubernetesAPIHealthz:         buildImage(registryDomain, "giantswarm/k8s-api-healthz", versions.KubernetesAPIHealthz),
		KubernetesNetworkSetupDocker: buildImage(registryDomain, "giantswarm/k8s-setup-network-environment", versions.KubernetesNetworkSetupDocker),
	}
}

func ExtractComponentVersions(releaseComponents []v1alpha1.ReleaseSpecComponent) (Versions, error) {
	var versions Versions

	{
		component, err := findComponent(releaseComponents, "kubernetes")
		if err != nil {
			return Versions{}, err
		}
		// cri-tools is released for each k8s minor version
		parsedVersion, err := semver.NewVersion(component.Version)
		if err != nil {
			return Versions{}, err
		}
		versions.CRITools = fmt.Sprintf("v%d.%d.0", parsedVersion.Major(), parsedVersion.Minor())
		versions.Kubernetes = fmt.Sprintf("v%s", component.Version)
	}

	{
		component, err := findComponent(releaseComponents, "etcd")
		if err != nil {
			return Versions{}, err
		}
		versions.Etcd = fmt.Sprintf("v%s", component.Version)
	}

	{
		component, err := findComponent(releaseComponents, "calico")
		if err != nil {
			return Versions{}, err
		}
		versions.Calico = fmt.Sprintf("v%s", component.Version)
	}

	return versions, nil
}

func buildImage(registryDomain string, repo string, tag string) string {
	return fmt.Sprintf("%s/%s:%s", registryDomain, repo, tag)
}

func findComponent(releaseComponents []v1alpha1.ReleaseSpecComponent, name string) (*v1alpha1.ReleaseSpecComponent, error) {
	for _, component := range releaseComponents {
		if component.Name == name {
			return &component, nil
		}
	}
	return nil, componentNotFoundError
}
