package v_3_0_0

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"reflect"
	"text/template"

	"github.com/giantswarm/microerror"
)

type Config struct {
	Params   Params
	Template string
}

func DefaultConfig() Config {
	return Config{
		Params:   Params{},
		Template: "",
	}
}

type CloudConfig struct {
	config   string
	params   Params
	template string
}

func New(config Config) (*CloudConfig, error) {
	if reflect.DeepEqual(config.Params, Params{}) {
		return nil, microerror.Maskf(invalidConfigError, "config.Params must not be empty")
	}
	if config.Template == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Template must not be empty")
	}

	c := &CloudConfig{
		config:   "",
		params:   config.Params,
		template: config.Template,
	}

	return c, nil
}

func (c *CloudConfig) ExecuteTemplate() error {
	tmpl, err := template.New("cloudconfig").Parse(c.template)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, c.params)
	if err != nil {
		return err
	}
	c.config = buf.String()

	return nil
}

func (c *CloudConfig) Base64() string {
	cloudConfigBytes := []byte(c.config)

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(cloudConfigBytes)
	w.Close()

	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func (c *CloudConfig) String() string {
	return c.config
}
