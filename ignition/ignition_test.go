package ignition

import (
	"bytes"
	"testing"
)

func TestConvertTemplatetoJSON(t *testing.T) {
	tests := []struct {
		yamlContent         []byte
		expectedJSONContent []byte
	}{
		{
			yamlContent: []byte(`---
ignition:
  version: "2.2.0"
  config: {}
  security: {}
  tls: {}
  timeout: {}
storage:
  disks: []
  filesystems: []
  files: []
systemd:
  units: []
users: []
networkd:
  units: []
`),
			expectedJSONContent: []byte(`{
  "ignition": {
    "config": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "2.2.0"
  },
  "networkd": {},
  "passwd": {},
  "storage": {},
  "systemd": {}
}
`),
		},
	}

	for _, tc := range tests {
		converted, err := ConvertTemplatetoJSON(tc.yamlContent)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(converted, tc.expectedJSONContent) {
			t.Fatalf("expected %#v, got %#v", string(tc.expectedJSONContent[:]), string(converted[:]))
		}
	}
}
