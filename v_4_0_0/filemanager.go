package v_4_0_0

import (
	"bytes"
	"encoding/base64"
	"os"
	"path/filepath"
	"text/template"

	"github.com/giantswarm/microerror"
)

// Files is map[string]string (k: filename, v: contents) for files that are fetched from disk
// and then filled with data.
type Files map[string]string

// RenderFiles walks over filesdir and parses all regular files with
// text/template. Parsed templates are then rendered with ctx, base64 encoded
// and added to returned Files.
//
// filesdir must not contain any other files than templates that can be parsed
// with text/template.
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
