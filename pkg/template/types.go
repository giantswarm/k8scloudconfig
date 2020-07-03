package template

import (
	"github.com/giantswarm/apiextensions/pkg/apis/provider/v1alpha1"
)

type Params struct {
	// APIServerEncryptionKey is AES-CBC with PKCS#7 padding key to encrypt API
	// etcd data.
	APIServerEncryptionKey string
	BaseDomain             string
	Cluster                v1alpha1.Cluster
	// Debug options
	Debug Debug
	// EnableAWSCNI flag. When set to true will use AWS CNI for pod networking
	// and Calico only for network policies.
	EnableAWSCNI bool
	// CalicoPolicyOnly flag. When set to true will deploy calico for network policies only.
	CalicoPolicyOnly bool
	// DisableEncryptionAtREST flag. When set removes all manifests from the cloud
	// config related to Kubernetes encryption at REST.
	DisableEncryptionAtREST bool
	// DisableIngressControllerService flag. When set removes the manifest for
	// the Ingress Controller service. This allows us to migrate providers to
	// chart-operator independently.
	DisableIngressControllerService bool
	Etcd                            Etcd
	Extension                       Extension
	// ExtraManifests allows to specify extra Kubernetes manifests in
	// /opt/k8s-addons script. The manifests are applied after calico is
	// ready.
	//
	// The general use-case is to create a manifest file with Extension and
	// then apply the manifest by adding it to ExtraManifests.
	ExtraManifests []string
	Files          Files
	// ImagePullProgressDeadline is the duration after which image pulling is
	// cancelled if no progress has been made.
	ImagePullProgressDeadline string
	// Container images used in the cloud-config templates
	Images Images
	// Kubernetes components allow the passing of extra `docker run` and
	// `command` arguments to image commands. This allows, for example,
	// the addition of cloud provider extensions.
	Kubernetes      Kubernetes
	Node            v1alpha1.ClusterNode
	RegistryDomain  string
	RegistryMirrors []string
	SSOPublicKey    string
	Versions        Versions
}

func (p *Params) Validate() error {
	return nil
}

type Versions struct {
	Calico                       string
	CRITools                     string
	Etcd                         string
	Kubernetes                   string
	KubernetesAPIHealthz         string
	KubernetesNetworkSetupDocker string
}

type Debug struct {
	Enabled    bool
	LogsPrefix string
	LogsToken  string
}

type Images struct {
	CalicoCNI                    string
	CalicoKubeControllers        string
	CalicoNode                   string
	Etcd                         string
	Hyperkube                    string
	KubeApiserver                string
	KubeControllerManager        string
	KubeScheduler                string
	KubeProxy                    string
	KubernetesAPIHealthz         string
	KubernetesNetworkSetupDocker string
	Pause                        string
}

type Kubernetes struct {
	Apiserver         KubernetesPodOptions
	ControllerManager KubernetesPodOptions
	HyperkubeWrappers bool
	Kubelet           KubernetesDockerOptions
}

type KubernetesDockerOptions struct {
	RunExtraArgs     []string
	CommandExtraArgs []string
}

type KubernetesPodOptions struct {
	HostExtraMounts  []KubernetesPodOptionsHostMount
	CommandExtraArgs []string
}

type KubernetesPodOptionsHostMount struct {
	Name     string
	Path     string
	ReadOnly bool
}

type Etcd struct {
	// ClientPort allows the port for clients to be specified.
	// aws-operator sets this to the Etcd listening port so Calico on the
	// worker nodes can access via a CNAME record to the master.
	ClientPort int
	// Enabled when set to true will cause rendering master template for cluster of 3 masters. Single master otherwise.
	// Defaults to false.
	HighAvailability bool
	// InitialCluster is config which define which etcd are members of the cluster.
	// The format should look like this: `etcd1=https://etcd1.example.com:2380,etcd2=https://etcd2.example.com:2380,etcd3=https://etcd3.example.com:2380`
	// Where etcd1.example.com, etcd2.example.com, and etcd3.example.com can be either the IP or DNS of the master machine
	// where is etcd listening.
	InitialCluster string
	// Initial cluster state for the etcd cluster. Should have values either `new` or `existing`.
	InitialClusterState string
	// NodeName is the name of the current etcd cluster node.
	NodeName string
}

type FileMetadata struct {
	AssetContent string
	Path         string
	Owner        Owner
	Compression  bool
	Permissions  int
}

type Owner struct {
	Group Group
	User  User
}

// Group object reflects spec for ignition Group object.
// If both ID and name are specified, ID is preferred.
type Group struct {
	ID   int
	Name string
}

// User object reflects spec for ignition User object.
// If both ID and name are specified, ID is preferred.
type User struct {
	ID   int
	Name string
}

type FileAsset struct {
	Metadata FileMetadata
	Content  string
}

type UnitMetadata struct {
	AssetContent string
	Name         string
	Enabled      bool
}

type UnitAsset struct {
	Metadata UnitMetadata
	Content  []string
}

// VerbatimSection is a blob of YAML we want to add to the
// CloudConfig, with no variable interpolation.
type VerbatimSection struct {
	Name    string
	Content string
}

type Extension interface {
	Files() ([]FileAsset, error)
	Units() ([]UnitAsset, error)
	VerbatimSections() []VerbatimSection
}
