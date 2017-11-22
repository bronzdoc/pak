package api

import (
	"encoding/json"
	"fmt"

	pak "github.com/bronzdoc/pak/pakfile"
	"github.com/jhoonb/archivex"
)

func Build(pakfile *pak.PakFile) {
	artifact := new(archivex.TarFile)
	artifact.Create(fmt.Sprintf("%s.tar", pakfile.ArtifactName))

	var metadataContent []byte
	buildMetadata := make(map[string]string)

	for key, value := range pakfile.Metadata {
		buildMetadata[key] = value
	}

	labels := map[string]map[string]string{
		"build": buildMetadata,
	}

	// Get labels and metadata from the Promote map
	for key, metadata := range pakfile.Promote {
		for _, value := range metadata {
			labels[key] = value
		}
	}

	jsonString, err := json.MarshalIndent(labels, "", "  ")
	if err != nil {
		panic(err)
	}

	metadataContent = append(metadataContent, jsonString...)

	// Store metadata in the package
	artifact.Add("pak.metadata", metadataContent)

	artifact.AddAll(pakfile.Path, true)

	artifact.Close()
}
