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

	labels := map[string]map[string]interface{}{
		"build": map[string]interface{}{
			"name":     pakfile.ArtifactName,
			"metadata": pakfile.Metadata,
		},
	}

	// Get labels and metadata from the Promote map
	for key, value := range pakfile.Promote {
		labels[key] = value
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
