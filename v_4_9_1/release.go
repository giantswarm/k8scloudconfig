package v_4_9_1

import (
	"fmt"

	"github.com/giantswarm/apiextensions/pkg/apis/release/v1alpha1"
)

func BuildImages(registryDomain string, versions Versions) Images {
	return Images{
		CalicoCNI:             buildImage("giantswarm/cni", registryDomain, versions.Calico),
		CalicoKubeControllers: buildImage("giantswarm/kube-controllers", registryDomain, versions.Calico),
		CalicoNode:            buildImage("giantswarm/node", registryDomain, versions.Calico),
		Etcd:                  buildImage("giantswarm/etcd", registryDomain, versions.Etcd),
		Hyperkube:             buildImage("giantswarm/hyperkube", registryDomain, versions.Kubernetes),
		Kubectl:               buildImage("giantswarm/docker-kubectl", registryDomain, versions.Kubectl),
		KubernetesAPIHealthz:  buildImage("giantswarm/k8s-api-healthz", registryDomain, versions.KubernetesAPIHealthz),
	}
}

func ExtractComponentVersions(releaseComponents []v1alpha1.ReleaseSpecComponent) (Versions, error) {
	var versions Versions

	{
		component, err := findComponent(releaseComponents, "kubernetes")
		if err != nil {
			return Versions{}, err
		}
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
