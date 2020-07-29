package template

import (
	"encoding/base64"
	"path"
	"testing"

	"github.com/giantswarm/k8scloudconfig/v7/pkg/ignition"
)

func TestCloudConfig(t *testing.T) {
	tests := []struct {
		name     string
		template string
		params   Params
	}{
		{
			name:     "master",
			template: MasterTemplate,
			params: Params{
				Etcd: Etcd{
					ClientPort: 443,
				},
			},
		},
		{
			name:     "worker",
			template: WorkerTemplate,
			params: Params{
				Etcd: Etcd{
					ClientPort: 443,
				},
			},
		},
		{
			name:     "worker",
			template: WorkerTemplate,
			params: Params{
				Etcd: Etcd{
					ClientPort: 2379,
				},
			},
		},
	}

	for _, tc := range tests {
		c := CloudConfigConfig{
			Params:   tc.params,
			Template: "",
		}

		packagePath, err := GetPackagePath()
		if err != nil {
			t.Errorf("failed to retrieve package path, %v:", err)
		}
		filesPath := path.Join(packagePath, filesDir)
		files, err := RenderFiles(filesPath, tc.params)
		if err != nil {
			t.Errorf("failed to render ignition files, %v:", err)
		}
		tc.params.Extension = nopExtension{}
		tc.params.Files = files
		tc.params.Versions = Versions{
			Calico:                       "3.14.1",
			CRITools:                     "1.17.0",
			Etcd:                         "3.4.9",
			Kubernetes:                   "1.17.7",
			KubernetesAPIHealthz:         "0.1.1",
			KubernetesNetworkSetupDocker: "0.2.0",
		}
		tc.params.Images = BuildImages("docker.io", tc.params.Versions)

		c.Params = tc.params
		c.Template = tc.template

		cloudConfig, err := NewCloudConfig(c)
		if err != nil {
			t.Fatal(err)
		}

		if err := cloudConfig.ExecuteTemplate(); err != nil {
			t.Fatal(err)
		}

		cloudConfigBase64 := cloudConfig.Base64()
		if _, err := base64.StdEncoding.DecodeString(cloudConfigBase64); err != nil {
			t.Errorf("The string isn't Base64 encoded: %v", err)
		}

		_, err = ignition.ConvertTemplatetoJSON([]byte(cloudConfig.String()))
		if err != nil {
			t.Fatalf("failed to validate ignition %#q config, %v:", tc.name, err)
		}

	}
}
