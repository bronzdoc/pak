package api

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/bronzdoc/pak/util"
	"github.com/jhoonb/archivex"
	"github.com/pkg/errors"
)

// Promote will promote the artifact to a given promote label
func Promote(artifactName string, options map[string]interface{}) (string, error) {
	inspectOptions := map[string]interface{}{}

	// get metadata
	jsonMetadata, err := Inspect(artifactName, inspectOptions)
	if err != nil {
		return "", errors.Wrap(err, "failed to inspect artifact metadata")
	}

	var metadata map[string]interface{}
	json.Unmarshal([]byte(jsonMetadata), &metadata)

	// Get metadata of the label we are promoting
	var labelMetadata map[string]interface{}
	var label string
	if value, ok := options["label"]; ok && value != "" {
		labelMetadata = metadata[value.(string)].(map[string]interface{})
		label = value.(string)
	}

	// Resolve env variables of the metadata
	for key, value := range labelMetadata {
		var newValue string
		var err error

		switch value.(type) {
		case string:
			newValue, err = util.ResolveEnvVar(value.(string))
			labelMetadata[key] = newValue
		case map[string]interface{}:
			resolvedMetadata := make(map[string]interface{})
			for k, v := range value.(map[string]interface{}) {
				newValue, err = util.ResolveEnvVar(v.(string))
				resolvedMetadata[k] = newValue
			}

			labelMetadata[key] = resolvedMetadata
		}

		if err != nil {
			return "", errors.Wrap(err, "failed to resolve env var")
		}
	}

	metadata[label] = labelMetadata

	// Create promote artifact
	var newArtifactname string
	if value, ok := labelMetadata["name"]; ok {
		newArtifactname = value.(string)
	} // generate random name if name empty?

	fullArtifactName := fmt.Sprintf("%s.tar", newArtifactname)
	newArtifact := new(archivex.TarFile)
	newArtifact.Create(fullArtifactName)

	// create new pak.metadata with the old data and the new one
	jsonString, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		panic(err)
	}

	// Store metadata in the new artifact
	newArtifact.Add("pak.metadata", jsonString)

	// Add packaged files from the old artifact to the new one
	pkg, err := os.Open(artifactName)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open artifact %s", artifactName)
	}

	tr := tar.NewReader(pkg)

	// Search for the metadata file inside the artifact
	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return "", errors.Wrap(err, "could not get tar header")
		}

		fileContent, err := ioutil.ReadAll(tr)
		if err != nil {
			return "", errors.Wrapf(err, "could not read artifact %s", artifactName)
		}

		// Don't add metadata file since we will be adding a new one with new metadata
		metadataFileName := "pak.metadata"
		if header.Name != metadataFileName {
			newArtifact.Add(header.Name, fileContent)
		}
	}

	newArtifact.Close()

	return fullArtifactName, nil
}
