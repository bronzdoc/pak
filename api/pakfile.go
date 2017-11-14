package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type PakFile struct {
	ArtifactName string            `json:"artifact_name"`
	Path         string            `json:"path"`
	Metadata     map[string]string `json:"metadata"`
}

func NewPakFile(jsonPath string) *PakFile {
	if _, err := os.Stat(jsonPath); err != nil {
		fmt.Println(err)
	}

	jsonContent, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		fmt.Println(err)
	}

	pakfile := PakFile{}
	json.Unmarshal(jsonContent, &pakfile)

	return pakfile
}
