package v_4_0_0

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/giantswarm/microerror"
)

// Files is map[string]string for files that we fetched from disk and then filled with data.
type Files map[string]string

func RenderFiles(filesdir string, ctx interface{}) (Files, error) {
	files := Files{}

	err := filepath.Walk(filesdir, func(path string, f os.FileInfo, err error) error {
		if f.Mode().IsRegular() {
			tmpl, err := template.ParseFiles(path)
			if err != nil {
				return microerror.Maskf(err, "failed to parse file #%q", path)
			}
			var data bytes.Buffer
			tmpl.Execute(&data, ctx)

			relativePath, err := filepath.Rel(filesdir, path)
			fmt.Printf("path: %s", relativePath)
			if err != nil {
				return microerror.Mask(err)
			}
			files[relativePath] = base64.StdEncoding.EncodeToString(data.Bytes())
		}
		return nil
	})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return files, nil
}
