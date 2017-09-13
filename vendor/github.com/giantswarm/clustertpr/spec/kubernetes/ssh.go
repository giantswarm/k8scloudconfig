package kubernetes

import "github.com/giantswarm/clustertpr/spec/kubernetes/ssh"

type SSH struct {
	// PublicKeys is a list of SSH public keys being added to each Kubernetes
	// node. It can contain admin specific public keys as well as customer
	// specific ones.
	UserList []ssh.User
}
