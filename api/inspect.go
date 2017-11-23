package api

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

func Inspect(packageName string, options map[string]interface{}) (string, error) {
	const metadataFileName string = "pak.metadata"

	pkg, err := os.Open(packageName)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open package %s", packageName)
	}

	tr := tar.NewReader(pkg)

	// Search for the metadata file inside the package
	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return "", errors.Wrap(err, "could not get tar header")
		}

		if header.Name == metadataFileName {
			break
		}
	}

	content, err := ioutil.ReadAll(tr)
	if err != nil {
		return "", errors.Wrapf(err, "could not read package %s", packageName)
	}

	if len(content) <= 0 {
		return "", fmt.Errorf("no metadata found")
	}

	// check if we need to inspect only subset of the metadata
	if value, ok := options["label"]; ok && value != "" {
		var mapContent map[string]map[string]string

		label := value.(string)

		json.Unmarshal(content, &mapContent)

		metadata := mapContent[label]

		if isKeyVal, ok := options["is_key_val"]; ok && isKeyVal.(bool) {
			var content []byte

			for key, value := range metadata {
				str := fmt.Sprintf("%s=%s\n", key, value)
				content = append(content, str...)
			}

			return string(content), nil
		} else {
			content, err := json.MarshalIndent(metadata, "", "  ")
			if err != nil {
				panic(err)
			}
			return string(content), nil
		}
	}

	return string(content), nil
}
