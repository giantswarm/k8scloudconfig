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

	"github.com/giantswarm/k8scloudconfig/v9/pkg/ignition"
)

const (
	InitialClusterStateNew      = "new"
	InitialClusterStateExisting = "existing"
)

type CloudConfigConfig struct {
	Params   Params
	Template string
}

type CloudConfig struct {
	config   string
	params   Params
	template string
}

func NewCloudConfig(config CloudConfigConfig) (*CloudConfig, error) {
	if err := config.Params.Validate(); err != nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.%s", config, err)
	}
	if config.Template == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Template must not be empty")
	}

	if config.Params.Etcd.NodeName == "" {
		if config.Params.Etcd.HighAvailability {
			// We can't guess the node name in this case so must return an error
			return nil, microerror.Maskf(invalidConfigError,
				"config.%T must be specified for HA etcd",
				config.Params.Etcd.NodeName)
		}
		config.Params.Etcd.NodeName = etcdNodeName(1, 1)
	}

	if config.Params.Etcd.InitialCluster == "" {
		config.Params.Etcd.InitialCluster = etcdInitialCluster(config.Params.BaseDomain, config.Params.Etcd.HighAvailability)
	}

	if !strings.Contains(config.Params.Etcd.InitialCluster, fmt.Sprintf("%s=", config.Params.Etcd.NodeName)) {
		return nil, microerror.Maskf(invalidConfigError,
			"initial cluster, %s, must contain node ID, %s",
			config.Params.Etcd.InitialCluster,
			config.Params.Etcd.NodeName)
	}

	if config.Params.DockerhubToken == "" {
		return nil, microerror.Maskf(
			invalidConfigError,
			"config.Params.DockerhubToken must be specified",
		)
	}

	{
		kubernetesVersion, err := semver.NewVersion(config.Params.Versions.Kubernetes)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		// Before Kubernetes 1.17, the hyperkube image only contained the hyperkube
		// binary which itself provided all of the constituent Kubernetes components
		// such as kubelet through subcommands like `hyperkube kubelet`. In 1.17, hyperkube
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

func etcdClusterSize(highAvailability bool) int {
	if highAvailability {
		return 3
	}
	return 1
}

func etcdInitialCluster(baseDomain string, highAvailability bool) string {
	var cluster string
	clusterSize := etcdClusterSize(highAvailability)
	for i := 1; i < clusterSize+1; i++ {
		id := etcdNodeName(i, clusterSize)
		cluster = fmt.Sprintf("%s,%s=https://%s.%s:2380", cluster, id, id, baseDomain)
	}
	return strings.TrimPrefix(cluster, ",")
}

func etcdNodeName(index int, clusterSize int) string {
	if clusterSize == 1 {
		return "etcd" // skip suffix for non-HA clusters for backwards compatibility
	}
	return fmt.Sprintf("etcd%d", index)
}
