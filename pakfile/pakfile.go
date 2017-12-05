package pakfile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/bronzdoc/pak/util"
	"github.com/pkg/errors"
)

// PakFile represents a Pakfile.json
type PakFile struct {
	ArtifactName string                            `json:"name"`
	Path         string                            `json:"path"`
	Metadata     map[string]string                 `json:"metadata"`
	Promote      map[string]map[string]interface{} `json:"promote"`
}

// Factory will create a new Pakfile will its necesary data
func Factory() (*PakFile, error) {
	var pakfile *PakFile
	jsonPakfile := "Pakfile.json"

	// check for a valid pakfile in the current directory
	if _, err := os.Stat(jsonPakfile); os.IsNotExist(err) {
		return pakfile, errors.Wrap(err, "Pakfile.json not found in the current directory")
	}

	pakfileContent, err := ioutil.ReadFile(jsonPakfile)
	if err != nil {
		errors.Wrap(err, "failed to read Pakfile.json")
	}

	pakfile = New(pakfileContent)

	if err := pakfile.GetData(); err != nil {
		return pakfile, errors.Wrap(err, "pakfile failed to get data")
	}

	return pakfile, nil
}

// New will create a new Pakfile
func New(jsonString []byte) *PakFile {
	pakfile := PakFile{}

	if err := json.Unmarshal(jsonString, &pakfile); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal Pakfile.json content"))
	}

	return &pakfile
}

// GetData will get the name, path and metadata of a Pakfile
func (p *PakFile) GetData() error {
	if err := p.GetArtifactName(); err != nil {
		return errors.Wrap(err, "failed to get artifact name")
	}

	if err := p.GetPath(); err != nil {
		return errors.Wrap(err, "failed to get path to build from")
	}

	if err := p.GetMetadata(); err != nil {
		return errors.Wrap(err, "failed to get metadata")
	}

	return nil
}

// GetMetadata will get a Pakfile metadata
func (p *PakFile) GetMetadata() error {
	if len(p.Metadata) <= 0 {
		return nil
	}

	for key, value := range p.Metadata {
		newValue, err := util.ResolveEnvVar(value)
		if err != nil {
			return errors.Wrapf(err, "failed to resolve env var \"%s\"", value)
		}

		p.Metadata[key] = newValue
	}

	return nil
}

// GetArtifactName will get a Pakfile artiface name
func (p *PakFile) GetArtifactName() error {
	// If artifact name is an env var, resolve it
	newArtifactName, err := util.ResolveEnvVar(p.ArtifactName)
	if err != nil {
		return errors.Wrapf(err, "failed to resolve env var \"%s\"", p.ArtifactName)
	}

	p.ArtifactName = newArtifactName

	// Check if no artifact name was given
	if p.ArtifactName == "" {
		// Build random artifact name if not given
		rand.Seed(int64(time.Now().Nanosecond()))
		p.ArtifactName = fmt.Sprintf("pak-artifact.%d", rand.Int())
	}

	return nil
}

// GetPath will get a Pakfile path
func (p *PakFile) GetPath() error {
	// If path is an env var, resolve it
	newPath, err := util.ResolveEnvVar(p.Path)
	if err != nil {
		return errors.Wrapf(err, "failed to resolve env var \"%s\"", p.ArtifactName)
	}

	p.Path = newPath

	// Check if path is empty
	if p.Path == "" {
		// pak can't build an artifact if no build path is given
		return fmt.Errorf("path can not be empty")
	}

	return nil
}
