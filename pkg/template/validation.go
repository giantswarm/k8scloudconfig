package template

import (
	"github.com/Masterminds/semver/v3"
	"github.com/giantswarm/microerror"
)

type versionConstraints struct {
	calico     string
	criTools   string
	etcd       string
	kubernetes string
}

var knownVersionConstraints = versionConstraints{
	calico:     ">= 3.10.0 < 3.15.0",
	criTools:   "",
	etcd:       ">= 3.4.0 < 3.5.0",
	kubernetes: ">= 1.16.0 < 1.18.0",
}

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
		return microerror.Maskf(invalidConfigError, "component %#q requires version following constraint %#q, got %#q", name, constraintString, versionString)
	}
	return nil
}

func (p *Params) Validate() error {
	if err := validateImagesRegsitry(p.Images, p.RegistryMirrors); err != nil {
		return microerror.Mask(err)
	}
	if err := validateComponentVersion("Kubernetes", p.Versions.Kubernetes, knownVersionConstraints.kubernetes); err != nil {
		return microerror.Mask(err)
	}
	if err := validateComponentVersion("Calico", p.Versions.Calico, knownVersionConstraints.calico); err != nil {
		return microerror.Mask(err)
	}
	if err := validateComponentVersion("Etcd", p.Versions.Etcd, knownVersionConstraints.etcd); err != nil {
		return microerror.Mask(err)
	}
	if err := validateComponentVersion("CRITools", p.Versions.CRITools, knownVersionConstraints.criTools); err != nil {
		return microerror.Mask(err)
	}
	return nil
}
