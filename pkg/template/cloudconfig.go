package template

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"log"
	"text/template"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/k8scloudconfig/v6/pkg/ignition"
)

const (
	defaultImagePullProgressDeadline = "1m"
	etcdPort                         = 443
)

type CloudConfigConfig struct {
	Params   Params
	Template string
}

func DefaultCloudConfigConfig() CloudConfigConfig {
	return CloudConfigConfig{
		Params:   Params{},
		Template: "",
	}
}

func DefaultParams() Params {
	return Params{
		EtcdPort:                  etcdPort,
		ImagePullProgressDeadline: defaultImagePullProgressDeadline,
		RegistryDomain:            "quay.io",
		MultiMasters: MultiMasters{
			Enabled:            false,
			EtcdInitialCluster: "",
			MasterID:           1,
		},
		Versions: Versions{
			Calico:   "1.0.0",
			CRITools: "1.0.0",
		},
	}
}

type CloudConfig struct {
	config   string
	params   Params
	template string
}

func NewCloudConfig(config CloudConfigConfig) (*CloudConfig, error) {
	if err := config.Params.Validate(); err != nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Params.%s", err)
	}
	if config.Template == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Template must not be empty")
	}
	if config.Params.MultiMasters.Enabled && config.Params.MultiMasters.EtcdInitialCluster == "" {
		config.Params.MultiMasters.EtcdInitialCluster = defaultEtcdMultiCluster(config.Params.BaseDomain)
	}
	c := &CloudConfig{
		config:   "",
		params:   config.Params,
		template: config.Template,
	}

	return c, nil
}

func (c *CloudConfig) ExecuteTemplate() error {
	tmpl, err := template.New("cloudconfig").Parse(c.template)
	if err != nil {
		return microerror.Mask(err)
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, c.params)
	if err != nil {
		return microerror.Mask(err)
	}

	ignitionJSON, err := ignition.ConvertTemplatetoJSON(buf.Bytes())
	if err != nil {
		return microerror.Mask(err)
	}

	c.config = string(ignitionJSON)

	return nil
}

func (c *CloudConfig) Base64() string {
	cloudConfigBytes := []byte(c.config)

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(cloudConfigBytes)
	if err != nil {
		log.Printf("failed to write gzip, reason: %#q", err.Error())
		return ""
	}
	err = w.Close()
	if err != nil {
		log.Printf("failed to close gzip, reason: %#q", err.Error())
		return ""
	}

	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func (c *CloudConfig) String() string {
	return c.config
}

func defaultEtcdMultiCluster(baseDomain string) string {
	return fmt.Sprintf("etcd1=etcd1.%s:2380,etcd2=etcd2.%s:2380,etcd3=etcd3.%s:2380", baseDomain, baseDomain, baseDomain)
}
