package api

import (
	"fmt"

	"github.com/jhoonb/archivex"
)

func Build(pakfile *PakFile) {
	artifact := new(archivex.TarFile)
	artifact.Create(pakfile.ArtifactName)

	// Store metadata in the package
	for key, value := range pakfile.Metadata {
		keyValPair := fmt.Sprintf("%s=%s", key, value)
		artifact.Add("pak.metadata", []byte(keyValPair))
	}

	artifact.AddAll(pakfile.Path, true)

	artifact.Close()
}
