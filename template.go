package cloudconfig

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/k8scloudconfig/v2"
	"github.com/giantswarm/k8scloudconfig/v_0_1_0"
)

type version string

const (
	V_0_1_0 version = "v_0_1_0"
	V2      version = "v2"
)

type Template struct {
	Master string
	Worker string
}

// NewTemplate returns a template structure containing cloud config templates
// versioned as provided.
//
// NOTE that version is a private type to this package to prevent users from
// accidentially providing wrong versions.
func NewTemplate(v version) (Template, error) {
	var template Template

	switch v {
	case V_0_1_0:
		template.Master = v_0_1_0.MasterTemplate
		template.Worker = v_0_1_0.WorkerTemplate
	case V2:
		template.Master = v2.MasterTemplate
		template.Worker = v2.WorkerTemplate
	default:
		return Template{}, microerror.Maskf(notFoundError, "template version '%s'", v)
	}

	return template, nil
}
