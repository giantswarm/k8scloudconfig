package v_3_0_0

import "github.com/giantswarm/apiextensions/pkg/apis/provider/v1alpha1"

type Params struct {
	Cluster   v1alpha1.Cluster
	Extension Extension
	Node      v1alpha1.ClusterNode
}

type FileMetadata struct {
	AssetContent string
	Path         string
	Owner        string
	Encoding     string
	Permissions  int
}

type FileAsset struct {
	Metadata FileMetadata
	Content  []string
}

type UnitMetadata struct {
	AssetContent string
	Name         string
	Enable       bool
	Command      string
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
