package v_4_0_0

import (
	"encoding/base64"
	"testing"
)

const (
	FilesDir = "files"
)

func TestCloudConfig(t *testing.T) {
	tests := []struct {
		template         string
		params           Params
		expectedEtcdPort int
	}{
		{
			template: MasterTemplate,
			params: Params{
				Extension: nopExtension{},
			},
			expectedEtcdPort: 443,
		},
		{
			template: WorkerTemplate,
			params: Params{
				Extension: nopExtension{},
			},
			expectedEtcdPort: 443,
		},
		{
			template: WorkerTemplate,
			params: Params{
				EtcdPort:  2379,
				Extension: nopExtension{},
			},
			expectedEtcdPort: 2379,
		},
	}

	for _, tc := range tests {
		c := DefaultCloudConfigConfig()

		filesPath, err := GetFilesPath()
		if err != nil {
			t.Error(err)
		}

		files, err := RenderFiles(filesPath, tc.params)
		if err != nil {
			t.Errorf("failed to render ignition files, %v:", err)
		}
		tc.params.Files = files

		c.Params = tc.params
		c.Template = tc.template

		cloudConfig, err := NewCloudConfig(c)
		if err != nil {
			t.Fatal(err)
		}

		if cloudConfig.params.EtcdPort != tc.expectedEtcdPort {
			t.Errorf("expected etcd port %q, got %q", tc.expectedEtcdPort, cloudConfig.params.EtcdPort)
		}

		if err := cloudConfig.ExecuteTemplate(); err != nil {
			t.Fatal(err)
		}

		cloudConfigBase64 := cloudConfig.Base64()
		if _, err := base64.StdEncoding.DecodeString(cloudConfigBase64); err != nil {
			t.Errorf("The string isn't Base64 encoded: %v", err)
		}
	}
}
