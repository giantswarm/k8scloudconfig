package key

// These consts are used by template.Params.Validate() to check that provided versions
// fall within the known-good ranges.
const (
	CalicoVersionConstraint                       = ">= 3.10.0 < 3.16.0"
	EtcdVersionConstraint                         = ">= 3.4.0 < 3.5.0"
	KubernetesVersionConstraint                   = ">= 1.16.0"
	KubernetesApiHealthzVersionConstraint         = ">= 0.1.1"
	KubernetesNetworkSetupDockerVersionConstraint = ">= 0.2.0"
)
