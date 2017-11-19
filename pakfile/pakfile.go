package pakfile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"time"

	"github.com/pkg/errors"
)

type PakFile struct {
	ArtifactName string            `json:"artifact_name"`
	Path         string            `json:"path"`
	Metadata     map[string]string `json:"metadata"`
}

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

func New(jsonString []byte) *PakFile {
	pakfile := PakFile{}

	if err := json.Unmarshal(jsonString, &pakfile); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal Pakfile.json content"))
	}

	return &pakfile
}

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

func (p *PakFile) GetMetadata() error {
	if len(p.Metadata) <= 0 {
		return nil
	}

	for key, value := range p.Metadata {
		valueIsEnvVar, err := match(value, `^\${.+}`)
		if err != nil {
			return errors.Wrap(err, "could not match value")
		}

		if valueIsEnvVar {
			envVar, err := search(value, `\w+`)
			if err != nil {
				return errors.Wrap(err, "could not search value")
			}

			p.Metadata[key] = os.Getenv(envVar)
		}
	}

	return nil
}

func (p *PakFile) GetArtifactName() error {
	nameIsEnvVar, err := match(p.ArtifactName, `^\${.+}`)
	if err != nil {
		return errors.Wrap(err, "could not match value")
	}

	if nameIsEnvVar {
		envVar, err := search(p.ArtifactName, `\w+`)
		if err != nil {
			return errors.Wrap(err, "could not search value")
		}

		p.ArtifactName = os.Getenv(envVar)
	}

	// Check if no artifact name was given
	if p.ArtifactName == "" {
		// Build random artifact name if not given
		rand.Seed(int64(time.Now().Nanosecond()))
		p.ArtifactName = fmt.Sprintf("pak-artifact.%d", rand.Int())
	}

	return nil
}

func (p *PakFile) GetPath() error {
	// Check if path was given
	if p.Path == "" {
		// pak can't build an artifact if no build path is given
		return fmt.Errorf("Path can not be empty")
	}

	nameIsEnvVar, err := match(p.Path, `^\${.+}`)
	if err != nil {
		return errors.Wrap(err, "could not match value")
	}

	if nameIsEnvVar {
		envVar, err := search(p.Path, `\w+`)
		if err != nil {
			return errors.Wrap(err, "could not search value")
		}

		p.Path = os.Getenv(envVar)
	}

	return nil
}

func match(str, pattern string) (bool, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false, errors.Wrap(err, "failed to compile regex pattern")
	}

	return regex.MatchString(str), nil
}

func search(str, pattern string) (string, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return "", errors.Wrap(err, "failed to compile regex pattern")
	}

	return regex.FindString(str), nil
}
