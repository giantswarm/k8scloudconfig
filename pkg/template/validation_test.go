package template

import (
	"testing"
)

const (
	kubernetesVersionGood = "1.17.6"
	kubernetesVersionBad  = "1.20.0"

	calicoVersionGood = "3.14.1"
	calicoVersionBad  = "3.15.0"

	etcdVersionGood = "3.4.9"
	etcdVersionBad  = "3.3.17"
)

func Test_Params_Validation(t *testing.T) {
	nilErrorMatcher := func(t *testing.T, err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	invalidConfigErrorMatcher := func(t *testing.T, err error) {
		if !IsInvalidConfig(err) {
			t.Fatal(err)
		}
	}
	testCases := []struct {
		errorMatcher func(t *testing.T, err error)
		name         string
		versions     Versions
	}{
		{
			errorMatcher: nilErrorMatcher,
			name:         "case 0: normal versions are valid",
			versions: Versions{
				Calico:                       calicoVersionGood,
				CRITools:                     "",
				Etcd:                         etcdVersionGood,
				Kubernetes:                   kubernetesVersionGood,
				KubernetesAPIHealthz:         "",
				KubernetesNetworkSetupDocker: "",
			},
		},
		{
			errorMatcher: invalidConfigErrorMatcher,
			name:         "case 1: unsupported versions are invalid",
			versions: Versions{
				Calico:                       calicoVersionBad,
				CRITools:                     "",
				Etcd:                         etcdVersionBad,
				Kubernetes:                   kubernetesVersionBad,
				KubernetesAPIHealthz:         "",
				KubernetesNetworkSetupDocker: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			params := Params{
				Images:   BuildImages("quay.io", tc.versions),
				Versions: tc.versions,
			}
			tc.errorMatcher(t, params.Validate())
		})
	}
}
