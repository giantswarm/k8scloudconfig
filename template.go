package cloudconfig

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/k8scloudconfig/v_0_1_0"
)

type version string

const (
	V_0_1_0 version = "v_0_1_0"
)

type Template struct {
	Master string
	Worker string
}

func NewTemplate(v version) (Template, error) {
	var template Template

	switch v {
	case V_0_1_0:
		template.Master = v_0_1_0.MasterTemplate
		template.Worker = v_0_1_0.WorkerTemplate
	default:
		return Template{}, microerror.Maskf(notFoundError, "template version '%s'", v)
	}

	return template, nil
}
