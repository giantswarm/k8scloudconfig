module github.com/giantswarm/k8scloudconfig/v16

go 1.18

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/giantswarm/apiextensions/v6 v6.0.0
	github.com/giantswarm/microerror v0.4.0
	github.com/stretchr/testify v1.8.1
	sigs.k8s.io/yaml v1.3.0
)

require github.com/giantswarm/release-operator/v4 v4.0.0

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/giantswarm/k8smetadata v0.19.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apimachinery v0.24.3 // indirect
	k8s.io/klog/v2 v2.70.1 // indirect
	k8s.io/utils v0.0.0-20220713171938-56c0de1e6f5e // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
)

replace (
	github.com/coreos/etcd v3.3.10+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/gogo/protobuf v1.3.1 => github.com/gogo/protobuf v1.3.2
	github.com/gorilla/websocket v1.4.0 => github.com/gorilla/websocket v1.4.2
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v1.11.1
	golang.org/x/text => golang.org/x/text v0.3.8
	sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v1.0.4
)
