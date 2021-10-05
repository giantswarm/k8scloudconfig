package template

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/Masterminds/semver/v3"
)

var releaseVersionsValid = Versions{
	Calico:                       "3.20.0",
	CRITools:                     "1.19.0",
	Etcd:                         "3.4.14",
	Kubernetes:                   "1.19.4",
	KubernetesAPIHealthz:         "0.1.1",
	KubernetesNetworkSetupDocker: "0.2.0",
}

var releaseVersionsInvalid = Versions{
	Calico:                       "3.18.0",
	CRITools:                     "1.15.0",
	Etcd:                         "3.3.15",
	Kubernetes:                   "1.15.5",
	KubernetesAPIHealthz:         "0.1.0",
	KubernetesNetworkSetupDocker: "0.1.0",
}

func editVersions(versions Versions, fieldName string, version string) Versions {
	entityType := reflect.TypeOf(versions)
	for i := 0; i < entityType.NumField(); i++ {
		structField := entityType.Field(i)
		if fieldName != structField.Name {
			continue
		}
		field := reflect.ValueOf(&versions).Elem().Field(i)
		if field.CanSet() {
			field.SetString(version)
		}
	}
	return versions
}

func Test_Params_Validation(t *testing.T) {
	nilErrorMatcher := func(t *testing.T, err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	validationErrorMatcher := func(t *testing.T, err error) {
		if !IsValidationError(err) {
			t.Fatal(err)
		}
	}
	invalidVersionMatcher := func(t *testing.T, err error) {
		if !errors.Is(err, semver.ErrInvalidSemVer) {
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
			name:         "case 0: valid release versions are valid",
			versions:     releaseVersionsValid,
		},
		{
			errorMatcher: invalidVersionMatcher,
			name:         "case 1: empty kubernetes version is invalid",
			versions:     editVersions(releaseVersionsValid, "Kubernetes", ""),
		},
		{
			errorMatcher: validationErrorMatcher,
			name:         "case 2: old kubernetes version is invalid",
			versions:     editVersions(releaseVersionsValid, "Kubernetes", releaseVersionsInvalid.Kubernetes),
		},
		{
			errorMatcher: validationErrorMatcher,
			name:         "case 3: old calico version is invalid",
			versions:     editVersions(releaseVersionsValid, "Calico", releaseVersionsInvalid.Calico),
		},
		{
			errorMatcher: validationErrorMatcher,
			name:         "case 4: old etcd version is invalid",
			versions:     editVersions(releaseVersionsValid, "Etcd", releaseVersionsInvalid.Etcd),
		},
		{
			errorMatcher: validationErrorMatcher,
			name:         "case 5: old critools version is invalid",
			versions:     editVersions(releaseVersionsValid, "CRITools", releaseVersionsInvalid.CRITools),
		},
		{
			errorMatcher: validationErrorMatcher,
			name:         "case 6: old api healthz version is invalid",
			versions:     editVersions(releaseVersionsValid, "KubernetesAPIHealthz", releaseVersionsInvalid.KubernetesAPIHealthz),
		},
		{
			errorMatcher: validationErrorMatcher,
			name:         "case 7: old network setup version is invalid",
			versions:     editVersions(releaseVersionsValid, "KubernetesNetworkSetupDocker", releaseVersionsInvalid.KubernetesNetworkSetupDocker),
		},
	}

	for i, tc := range testCases {
		tc := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)
			params := Params{
				Images:   BuildImages("quay.io", tc.versions),
				Versions: tc.versions,
			}
			tc.errorMatcher(t, params.Validate())
		})
	}
}
