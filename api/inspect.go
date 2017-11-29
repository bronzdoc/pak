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

	var mapContent map[string]interface{}
	if err := json.Unmarshal(content, &mapContent); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal content")
	}

	// check if we need to inspect only subset of the metadata
	if label, ok := options["label"]; ok && label != "" {
		filteredMetadata := filterContentByLabel(mapContent, label.(string))

		if isKeyVal, ok := options["is_key_value"]; ok && isKeyVal.(bool) {
			var metadata []byte

			for key, value := range filteredMetadata {
				switch value.(type) {
				default:
					str := fmt.Sprintf("%s=\"%s\"\n", key, value)
					metadata = append(metadata, str...)
				case map[string]interface{}:
					for k, v := range value.(map[string]interface{}) {
						str := fmt.Sprintf("%s=\"%s\"\n", k, v)
						metadata = append(metadata, str...)
					}
				}
			}

			return string(metadata), nil
		} else {
			metadata, err := json.MarshalIndent(filteredMetadata, "", "  ")
			if err != nil {
				return "", errors.Wrap(err, "failded to marshal metadata")
			}

			return string(metadata), nil
		}
	}

	// Inspect all metadata
	if isKeyVal, ok := options["is_key_value"]; ok && isKeyVal.(bool) {
		var metadata []byte

		for key, value := range mapContent {
			str := fmt.Sprintf("#%s\n", key)
			metadata = append(metadata, str...)
			switch value.(type) {
			default:
				return "", errors.Wrap(fmt.Errorf("Pakfile.json error"), "invalid Pakfile.json")
			case map[string]interface{}:
				for key, v := range value.(map[string]interface{}) {
					switch v.(type) {
					default:
						str := fmt.Sprintf("  %s=\"%s\"\n", key, v)
						metadata = append(metadata, str...)
					case map[string]interface{}:
						for k, v1 := range v.(map[string]interface{}) {
							str := fmt.Sprintf("  %s=\"%s\"\n", k, v1)
							metadata = append(metadata, str...)
						}
					}
				}
			}
		}

		return string(metadata), nil
	} else {
		return string(content), nil
	}
}

func filterContentByLabel(metadata map[string]interface{}, label string) map[string]interface{} {
	if _, ok := metadata[label]; ok {
		return metadata[label].(map[string]interface{})
	}

	return make(map[string]interface{})
}
