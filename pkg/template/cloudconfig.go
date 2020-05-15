package template

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/Masterminds/semver/v3"
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
		ImagePullProgressDeadline: defaultImagePullProgressDeadline,
		RegistryDomain:            "quay.io",
		Etcd: Etcd{
			ClientPort:       etcdPort,
			HighAvailability: false,
		},
		Versions: Versions{
			Calico:     "1.0.0",
			CRITools:   "1.0.0",
			Kubernetes: "v1.17.5",
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
	if config.Params.Etcd.NodeName == "" {
		if config.Params.Etcd.HighAvailability {
			return nil, microerror.Maskf(invalidConfigError,
				"config.%T must be specified for HA etcd",
				config.Params.Etcd.NodeName)
		}
		config.Params.Etcd.NodeName = nodeName(1)
	}

	if config.Params.Etcd.InitialCluster == "" {
		etcdClusterSize := 1
		if config.Params.Etcd.HighAvailability {
			etcdClusterSize = 3
		}
		config.Params.Etcd.InitialCluster = defaultEtcdMultiCluster(config.Params.BaseDomain, etcdClusterSize)
	}

	if !strings.Contains(config.Params.Etcd.InitialCluster, fmt.Sprintf("%s=", config.Params.Etcd.NodeName)) {
		return nil, microerror.Maskf(invalidConfigError,
			"initial cluster, %s, must contain node ID, %s",
			config.Params.Etcd.InitialCluster,
			config.Params.Etcd.NodeName)
	}

	{
		kubernetesVersion, err := semver.NewVersion(config.Params.Versions.Kubernetes)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		// Before Kubernetes 1.17, the hyperkube image only contained the hyperkube
		// binary which itself provided all of the constituent Kubernetes components
		// such as kubelet through subcommands like `hyperkube kubelet`. In 17, hyperkube
		// began to become unsupported becoming a wrapper which simply called binaries
		// bundled into the docker image. For 1.16, we need to copy hyperkube out and
		// create wrappers, for 1.17+, we need to just copy binaries out.
		config.Params.Kubernetes.HyperkubeWrappers = kubernetesVersion.Minor() < 17
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

func nodeName(index int) string {
	return fmt.Sprintf("etcd%d", index)
}

func defaultEtcdMultiCluster(baseDomain string, count int) string {
	var cluster string
	for i := 1; i < count+1; i++ {
		id := nodeName(i)
		cluster = fmt.Sprintf("%s,%s=%s.%s:2380", cluster, id, id, baseDomain)
	}
	return strings.TrimPrefix(cluster, ",")
}
