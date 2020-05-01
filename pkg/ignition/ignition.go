package ignition

import (
	"encoding/json"

	"sigs.k8s.io/yaml"

	"github.com/giantswarm/microerror"
)

func ConvertTemplatetoJSON(dataIn []byte) ([]byte, error) {
	cfg := Config{}

	if err := yaml.Unmarshal(dataIn, &cfg); err != nil {
		return nil, microerror.Mask(err)
	}

	dataOut, err := json.MarshalIndent(&cfg, "", "  ")
	if err != nil {
		return nil, microerror.Mask(err)
	}
	dataOut = append(dataOut, '\n')

	return dataOut, nil
}
