package pakfile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"time"
)

type PakFile struct {
	ArtifactName string            `json:"artifact_name"`
	Path         string            `json:"path"`
	Metadata     map[string]string `json:"metadata"`
}

func New(jsonPath string) *PakFile {
	if _, err := os.Stat(jsonPath); err != nil {
		fmt.Println(err)
	}

	jsonContent, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		fmt.Println(err)
	}

	pakfile := PakFile{}
	json.Unmarshal(jsonContent, &pakfile)

	pakfile.parseArtifactName()
	pakfile.parseMetadata()
	pakfile.parsePath()

	if pakfile.ArtifactName == "" {
		// Build random artifact name if not given
		rand.Seed(int64(time.Now().Nanosecond()))
		pakfile.ArtifactName = fmt.Sprintf("pak-artifact.%d", rand.Int())
	}

	return &pakfile
}

func (p *PakFile) parseMetadata() {
	if len(p.Metadata) <= 0 {
		return
	}

	for key, value := range p.Metadata {
		valueIsEnvVar, err := match(value, `^\${.+}`)
		if err != nil {
			fmt.Println(err)
			return
		}

		if valueIsEnvVar {
			envVar, err := search(value, `\w+`)
			if err != nil {
				fmt.Println(err)
				return
			}

			p.Metadata[key] = os.Getenv(envVar)
		}
	}
}

func (p *PakFile) parseArtifactName() {
	nameIsEnvVar, err := match(p.ArtifactName, `^\${.+}`)
	if err != nil {
		fmt.Println(err)
		return
	}

	if nameIsEnvVar {
		envVar, err := search(p.ArtifactName, `\w+`)
		if err != nil {
			fmt.Println(err)
			return
		}

		p.ArtifactName = os.Getenv(envVar)
	}
}

func (p *PakFile) parsePath() {
	nameIsEnvVar, err := match(p.Path, `^\${.+}`)
	if err != nil {
		fmt.Println(err)
		return
	}

	if nameIsEnvVar {
		envVar, err := search(p.Path, `\w+`)
		if err != nil {
			fmt.Println(err)
			return
		}

		p.Path = os.Getenv(envVar)
	}
}

func match(str, pattern string) (bool, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}

	return regex.MatchString(str), nil
}

func search(str, pattern string) (string, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	return regex.FindString(str), nil
}
