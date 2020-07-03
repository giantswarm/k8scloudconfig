package template

import (
	"encoding/base64"
	"path"
	"testing"

	"github.com/giantswarm/k8scloudconfig/v6/pkg/ignition"
)

func TestCloudConfig(t *testing.T) {
	tests := []struct {
		name             string
		template         string
		params           Params
		customEtcdPort   int
		expectedEtcdPort int
	}{
		{
			name:             "master",
			template:         MasterTemplate,
			params:           DefaultParams(),
			expectedEtcdPort: 443,
		},
		{
			name:             "worker",
			template:         WorkerTemplate,
			params:           DefaultParams(),
			expectedEtcdPort: 443,
		},
		{
			name:             "worker",
			template:         WorkerTemplate,
			params:           DefaultParams(),
			customEtcdPort:   2379,
			expectedEtcdPort: 2379,
		},
	}

	for _, tc := range tests {
		c := DefaultCloudConfigConfig()

		tc.params.Extension = nopExtension{}

		if tc.customEtcdPort != 0 {
			tc.params.Etcd.ClientPort = tc.customEtcdPort
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
		tc.params.Files = files
		tc.params.Images = BuildImages(tc.params.RegistryDomain, Versions{})

		c.Params = tc.params
		c.Template = tc.template

		cloudConfig, err := NewCloudConfig(c)
		if err != nil {
			t.Fatal(err)
		}

		if cloudConfig.params.Etcd.ClientPort != tc.expectedEtcdPort {
			t.Errorf("expected etcd port %q, got %q", tc.expectedEtcdPort, cloudConfig.params.Etcd.ClientPort)
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
