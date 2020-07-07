package template

import (
	"strconv"
	"testing"
)

func Test_validateImagesRegsitry(t *testing.T) {
	testCases := []struct {
		name         string
		inputImages  Images
		inputMirrors []string
		errorMatcher func(err error) bool
	}{
		{
			name:         "case 0: ok",
			inputImages:  BuildImages("docker.io", Versions{}),
			inputMirrors: []string{"giantswarm.azurecr.io"},
			errorMatcher: nil,
		},
		{
			name:         "case 1: ok - no mirrors",
			inputImages:  BuildImages("quay.io", Versions{}),
			inputMirrors: nil,
			errorMatcher: nil,
		},
		{
			name:         "case 2: non-docker registry when mirrors are set",
			inputImages:  BuildImages("quay.io", Versions{}),
			inputMirrors: []string{"giantswarm.azurecr.io"},
			errorMatcher: IsInvalidConfig,
		},
		{
			name:         "case 3: empty Images",
			inputImages:  Images{},
			errorMatcher: IsInvalidConfig,
		},
		{
			name: "case 4: manually set images",
			inputImages: Images{
				CalicoCNI:                    "docker.io/giantswarm/image:1.2.3",
				CalicoKubeControllers:        "docker.io/giantswarm/image:1.2.3",
				CalicoNode:                   "docker.io/giantswarm/image:1.2.3",
				Etcd:                         "docker.io/giantswarm/image:1.2.3",
				Hyperkube:                    "docker.io/giantswarm/image:1.2.3",
				KubeApiserver:                "docker.io/giantswarm/image:1.2.3",
				KubeControllerManager:        "docker.io/giantswarm/image:1.2.3",
				KubeScheduler:                "docker.io/giantswarm/image:1.2.3",
				KubeProxy:                    "docker.io/giantswarm/image:1.2.3",
				KubernetesAPIHealthz:         "docker.io/giantswarm/image:1.2.3",
				KubernetesNetworkSetupDocker: "docker.io/giantswarm/image:1.2.3",
				Pause:                        "docker.io/giantswarm/image:1.2.3",
			},
			errorMatcher: nil,
		},
		{
			name: "case 5: manually set images - one is wrong",
			inputImages: Images{
				CalicoCNI:                    "docker.io/giantswarm/image:1.2.3",
				CalicoKubeControllers:        "docker.io/giantswarm/image:1.2.3",
				CalicoNode:                   "docker.io/giantswarm/image:1.2.3",
				Etcd:                         "docker.io/giantswarm/image:1.2.3",
				Hyperkube:                    "docker.io/giantswarm/image:1.2.3",
				KubeApiserver:                "docker.io/giantswarm/image:1.2.3",
				KubeControllerManager:        "quay.io/giantswarm/image:1.2.3",
				KubeScheduler:                "docker.io/giantswarm/image:1.2.3",
				KubeProxy:                    "docker.io/giantswarm/image:1.2.3",
				KubernetesAPIHealthz:         "docker.io/giantswarm/image:1.2.3",
				KubernetesNetworkSetupDocker: "docker.io/giantswarm/image:1.2.3",
				Pause:                        "docker.io/giantswarm/image:1.2.3",
			},
			errorMatcher: IsInvalidConfig,
		},
		{
			name: "case 6: manually set images - one is empty",
			inputImages: Images{
				CalicoCNI:                    "docker.io/giantswarm/image:1.2.3",
				CalicoKubeControllers:        "docker.io/giantswarm/image:1.2.3",
				CalicoNode:                   "",
				Etcd:                         "docker.io/giantswarm/image:1.2.3",
				Hyperkube:                    "docker.io/giantswarm/image:1.2.3",
				KubeApiserver:                "docker.io/giantswarm/image:1.2.3",
				KubeControllerManager:        "docker.io/giantswarm/image:1.2.3",
				KubeScheduler:                "docker.io/giantswarm/image:1.2.3",
				KubeProxy:                    "docker.io/giantswarm/image:1.2.3",
				KubernetesAPIHealthz:         "docker.io/giantswarm/image:1.2.3",
				KubernetesNetworkSetupDocker: "docker.io/giantswarm/image:1.2.3",
				Pause:                        "docker.io/giantswarm/image:1.2.3",
			},
			errorMatcher: IsInvalidConfig,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			err := validateImagesRegsitry(tc.inputImages, tc.inputMirrors)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

			if tc.errorMatcher != nil {
				return
			}
		})
	}
}
