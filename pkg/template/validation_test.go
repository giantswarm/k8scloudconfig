package template

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/Masterminds/semver/v3"
)

var releaseVersionsAWS_11_5_0 = Versions{
	Calico:                       "3.10.4",
	CRITools:                     "1.16.0",
	Etcd:                         "3.4.9",
	Kubernetes:                   "1.16.13",
	KubernetesAPIHealthz:         "0.1.1",
	KubernetesNetworkSetupDocker: "0.2.0",
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
	invalidConfigErrorMatcher := func(t *testing.T, err error) {
		if !IsInvalidConfig(err) {
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
			name:         "case 0: aws release 11.5.0 versions are valid",
			versions:     releaseVersionsAWS_11_5_0,
		},
		{
			errorMatcher: invalidVersionMatcher,
			name:         "case 1: empty kubernetes version is invalid",
			versions:     editVersions(releaseVersionsAWS_11_5_0, "Kubernetes", ""),
		},
		{
			errorMatcher: invalidConfigErrorMatcher,
			name:         "case 2: old kubernetes version is invalid",
			versions:     editVersions(releaseVersionsAWS_11_5_0, "Kubernetes", "1.15.0"),
		},
	}

	for i, tc := range testCases {
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
