package v_4_0_0

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"path"
	"runtime"
	"text/template"

	"github.com/giantswarm/microerror"
)

// Files is map[string]string for files that we fetched from disk and then filled with data.
type Files map[string]string

func RenderFiles(filespath string, ctx interface{}) (*Files, error) {
	files := Files{}
	dirList, err := ioutil.ReadDir(filespath)
	if err != nil {
		return nil, microerror.Maskf(err, "Failed to read files dir: %s, error: %#v", filespath, err)
	}

	for _, dir := range dirList {
		fileList, err := ioutil.ReadDir(path.Join(filespath, dir.Name()))
		if err != nil {
			return nil, microerror.Maskf(err, "Failed to read dir: %s, error: %#v", path.Join(filespath, dir.Name()), err)
		}

		for _, file := range fileList {
			tmpl, err := template.ParseFiles(path.Join(filespath, dir.Name(), file.Name()))
			if err != nil {
				return nil, microerror.Maskf(err, "Failed to file: %s, error: %#v", path.Join(filespath, dir.Name(), file.Name()), err)
			}
			var data bytes.Buffer
			tmpl.Execute(&data, ctx)

			files[dir.Name()+"/"+file.Name()] = base64.StdEncoding.EncodeToString(data.Bytes())
		}
	}

	return &files, nil
}

// GetFilesPath retrieves runtime path for the ignition templates
func GetFilesPath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", microerror.New("failed to retrieve runtime information")
	}
	filesPath := path.Join(path.Dir(filename), FilesDir)

	return filesPath, nil
}
