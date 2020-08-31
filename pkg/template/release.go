package template

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/giantswarm/apiextensions/pkg/apis/release/v1alpha1"
	"github.com/giantswarm/microerror"
)

func BuildImages(registryDomain string, versions Versions) Images {
	return Images{
		CalicoCNI:                    buildImage(registryDomain, "giantswarm/cni", versions.Calico, ""),
		CalicoKubeControllers:        buildImage(registryDomain, "giantswarm/kube-controllers", versions.Calico, ""),
		CalicoNode:                   buildImage(registryDomain, "giantswarm/node", versions.Calico, ""),
		Etcd:                         buildImage(registryDomain, "giantswarm/etcd", versions.Etcd, ""),
		Hyperkube:                    buildImage(registryDomain, "giantswarm/hyperkube", versions.Kubernetes, ""),
		KubeApiserver:                buildImage(registryDomain, "giantswarm/kube-apiserver", versions.Kubernetes, "-giantswarm"),
		KubeControllerManager:        buildImage(registryDomain, "giantswarm/kube-controller-manager", versions.Kubernetes, ""),
		KubeProxy:                    buildImage(registryDomain, "giantswarm/kube-proxy", versions.Kubernetes, ""),
		KubeScheduler:                buildImage(registryDomain, "giantswarm/kube-scheduler", versions.Kubernetes, ""),
		KubernetesAPIHealthz:         buildImage(registryDomain, "giantswarm/k8s-api-healthz", versions.KubernetesAPIHealthz, ""),
		KubernetesNetworkSetupDocker: buildImage(registryDomain, "giantswarm/k8s-setup-network-environment", versions.KubernetesNetworkSetupDocker, ""),
		Pause:                        buildImage(registryDomain, "giantswarm/pause", "3.2", ""),
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

func buildImage(registryDomain, repo, tag, suffix string) string {
	return fmt.Sprintf("%s/%s:%s%s", registryDomain, repo, tag, suffix)
}

func findComponent(releaseComponents []v1alpha1.ReleaseSpecComponent, name string) (*v1alpha1.ReleaseSpecComponent, error) {
	for _, component := range releaseComponents {
		if component.Name == name {
			return &component, nil
		}
	}
	return nil, componentNotFoundError
}

func validateImagesRegsitry(images Images, mirrors []string) error {
	data, err := json.Marshal(images)
	if err != nil {
		return microerror.Mask(err)
	}

	var m map[string]string
	err = json.Unmarshal(data, &m)
	if err != nil {
		return microerror.Mask(err)
	}

	var firstImage string
	var firstKey string
	var firstRegistry string

	for k, image := range m {
		split := strings.Split(image, "/")
		r := split[0]
		if firstImage == "" {
			firstImage = image
			firstKey = k
			firstRegistry = r
		}

		if r == "" {
			return microerror.Maskf(invalidConfigError, "%T.%s image %#q registry domain must not be empty", images, k, image)
		}

		if len(mirrors) > 0 && r != "docker.io" {
			return microerror.Maskf(invalidConfigError, "%T.%s image %#q registry domain must be %#q when mirrors are set", images, k, image, "docker.io")
		}

		if r != firstRegistry {
			return microerror.Maskf(invalidConfigError, "%T.%s image %#q and %T.%s image %#q have different registries domains", images, firstKey, firstImage, images, k, image)
		}
	}

	return nil
}
