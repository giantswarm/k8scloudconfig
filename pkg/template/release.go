package template

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/giantswarm/apiextensions/pkg/apis/release/v1alpha1"
)

func BuildImages(registryDomain string, versions Versions) Images {
	return Images{
		CalicoCNI:                    buildImage("giantswarm/cni", versions.Calico, ""),
		CalicoKubeControllers:        buildImage("giantswarm/kube-controllers", versions.Calico, ""),
		CalicoNode:                   buildImage("giantswarm/node", versions.Calico, ""),
		Etcd:                         buildImage("giantswarm/etcd", versions.Etcd, ""),
		Hyperkube:                    buildImage("giantswarm/hyperkube", versions.Kubernetes, ""),
		KubeApiserver:                buildImage("giantswarm/kube-apiserver", versions.Kubernetes, "-giantswarm"),
		KubeControllerManager:        buildImage("giantswarm/kube-controller-manager", versions.Kubernetes, ""),
		KubeProxy:                    buildImage("giantswarm/kube-proxy", versions.Kubernetes, ""),
		KubeScheduler:                buildImage("giantswarm/kube-scheduler", versions.Kubernetes, ""),
		KubernetesAPIHealthz:         buildImage("giantswarm/k8s-api-healthz", versions.KubernetesAPIHealthz, ""),
		KubernetesNetworkSetupDocker: buildImage("giantswarm/k8s-setup-network-environment", versions.KubernetesNetworkSetupDocker, ""),
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

func buildImage(repo string, tag string, suffix string) string {
	return fmt.Sprintf("%s:%s%s", repo, tag, suffix)
}

func findComponent(releaseComponents []v1alpha1.ReleaseSpecComponent, name string) (*v1alpha1.ReleaseSpecComponent, error) {
	for _, component := range releaseComponents {
		if component.Name == name {
			return &component, nil
		}
	}
	return nil, componentNotFoundError
}
