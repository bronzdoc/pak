package api

import (
	"fmt"

	pak "github.com/bronzdoc/pak/pakfile"
	"github.com/jhoonb/archivex"
)

func Build(pakfile *pak.PakFile) {
	artifact := new(archivex.TarFile)
	artifact.Create(fmt.Sprintf("%s.tar", pakfile.ArtifactName))

	var metadataContent []byte
	for key, value := range pakfile.Metadata {
		keyValPair := fmt.Sprintf("%s=%s\n", key, value)
		metadataContent = append(metadataContent, keyValPair...)
	}

	// Store metadata in the package
	artifact.Add("pak.metadata", metadataContent)

	artifact.AddAll(pakfile.Path, true)

	artifact.Close()
}
