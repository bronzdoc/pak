package pak

import (
	"github.com/jhoonb/archivex"
)

func Build(pakfile *PakFile) {
	artifact := new(archivex.TarFile)
	artifact.Create(pakfile.ArtifactName)
	artifact.Add("pak.metadata", []byte(pakfile.Metadata["name"]))
	//artifact.AddFile("test.sample")
	artifact.AddAll(pakfile.Path, true)

	artifact.Close()
}
