package key

// These consts are used by template.Params.Validate() to check that provided versions
// fall within the known-good ranges.
const (
	CalicoPolicyOnlyVersionConstraint             = ">= 3.21.0 < 3.22.0"
	CalicoVersionConstraint                       = ">= 3.20.0 < 3.21.0"
	EtcdVersionConstraint                         = ">= 3.4.0 < 3.6.0"
	KubernetesVersionConstraint                   = ">= 1.19.0"
	KubernetesApiHealthzVersionConstraint         = ">= 0.1.1"
	KubernetesNetworkSetupDockerVersionConstraint = ">= 0.2.0"
)
