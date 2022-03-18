package template

import (
	"github.com/Masterminds/semver/v3"
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/k8scloudconfig/v12/pkg/key"
)

func validateComponentVersion(name, versionString, constraintString string) error {
	if constraintString == "" {
		return nil
	}

	version, err := semver.NewVersion(versionString)
	if err != nil {
		return microerror.Mask(err)
	}

	constraint, err := semver.NewConstraint(constraintString)
	if err != nil {
		return microerror.Mask(err)
	}

	if !constraint.Check(version) {
		return microerror.Maskf(validationError, "component %#q requires version following constraint %#q, got %#q", name, constraintString, versionString)
	}

	return nil
}

func (p *Params) Validate() error {
	if err := validateImagesRegistry(p.Images, p.RegistryMirrors); err != nil {
		return microerror.Mask(err)
	}

	if err := validateComponentVersion("Kubernetes", p.Versions.Kubernetes, key.KubernetesVersionConstraint); err != nil {
		return microerror.Mask(err)
	}

	calicoVersionConstraint := key.CalicoVersionConstraint
	if p.CalicoPolicyOnly {
		calicoVersionConstraint = key.CalicoPolicyOnlyVersionConstraint
	}

	if err := validateComponentVersion("Calico", p.Versions.Calico, calicoVersionConstraint); err != nil {
		return microerror.Mask(err)
	}

	if err := validateComponentVersion("Etcd", p.Versions.Etcd, key.EtcdVersionConstraint); err != nil {
		return microerror.Mask(err)
	}

	// CRI tools follow kubernetes versioning so we'll reuse the version constraint
	if err := validateComponentVersion("CRITools", p.Versions.CRITools, key.KubernetesVersionConstraint); err != nil {
		return microerror.Mask(err)
	}

	if err := validateComponentVersion("KubernetesAPIHealthz", p.Versions.KubernetesAPIHealthz, key.KubernetesApiHealthzVersionConstraint); err != nil {
		return microerror.Mask(err)
	}

	if err := validateComponentVersion("KubernetesNetworkSetupDocker", p.Versions.KubernetesNetworkSetupDocker, key.KubernetesNetworkSetupDockerVersionConstraint); err != nil {
		return microerror.Mask(err)
	}

	return nil
}
