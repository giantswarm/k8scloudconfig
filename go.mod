module github.com/giantswarm/k8scloudconfig/v10

go 1.14

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/giantswarm/apiextensions/v3 v3.13.0
	github.com/giantswarm/microerror v0.3.0
	github.com/stretchr/testify v1.7.0
	sigs.k8s.io/yaml v1.2.0
)

replace sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.10-gs
