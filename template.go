package cloudconfig

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/k8scloudconfig/v_0_1_0_3500632"
)

type Template struct {
	Master string
	Worker string
}

func NewTemplate(v string) (Template, error) {
	var template Template

	switch v {
	case "v_0_1_0_3500632":
		template.Master = v_0_1_0_3500632.MasterTemplate
		template.Worker = v_0_1_0_3500632.WorkerTemplate
	default:
		return Template{}, microerror.Maskf(notFoundError, "template version '%s'", v)
	}

	return template, nil
}
